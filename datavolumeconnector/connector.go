package datavolumeconnector

import (
	"context"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
	"time"
)

type connectorImp struct {
	config          Config
	metricsConsumer consumer.Metrics
	logger          *zap.Logger
	component.StartFunc
	component.ShutdownFunc
}

const (
	dataTypeAttributeKey          = "data_type"
	dataTypeLogsAttributeValue    = "logs"
	dataTypeTracesAttributeValue  = "traces"
	dataTypeMetricsAttributeValue = "metrics"
)

var (
	plogSizer    = plog.ProtoMarshaler{}
	ptraceSizer  = ptrace.ProtoMarshaler{}
	pmetricSizer = pmetric.ProtoMarshaler{}
)

func newConnector(logger *zap.Logger, config component.Config) (*connectorImp, error) {
	cfg := config.(*Config)

	return &connectorImp{
		config: *cfg,
		logger: logger,
	}, nil
}

func (c *connectorImp) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

func (c *connectorImp) ConsumeLogs(ctx context.Context, logs plog.Logs) error {
	outputMetrics := pmetric.NewMetrics()

	for i := 0; i < logs.ResourceLogs().Len(); i++ {
		timestamp := pcommon.NewTimestampFromTime(time.Now())
		resourceLogs := logs.ResourceLogs().At(i)

		resourceAttributes := resourceLogs.Resource().Attributes()
		rawAttributes := resourceAttributes.AsRaw()

		metricAttrMap := map[string]any{}

		metricAttrMap[dataTypeAttributeKey] = dataTypeLogsAttributeValue
		for _, key := range c.config.LabelResourceAttributes {
			if rawAttributes[key] != nil {
				metricAttrMap[key] = rawAttributes[key]
			} else {
				metricAttrMap[key] = "unknown"
			}
		}

		resourceMetrics := outputMetrics.ResourceMetrics().AppendEmpty()
		err := resourceMetrics.Resource().Attributes().FromRaw(metricAttrMap)
		if err != nil {
			c.logger.Error("error adding attributes to datavolume metric for logs measurement", zap.Error(err), zap.Any("attributes_map", metricAttrMap))
		}
		scopeMetric := resourceMetrics.ScopeMetrics().AppendEmpty()

		if c.config.CountMetricName != "" {
			countValue := int64(resourceLogs.ScopeLogs().Len())
			addOutputMetricToScopeMetrics(scopeMetric, c.config.CountMetricName, "", timestamp, countValue)
		}

		if c.config.BytesMetricName != "" {
			bytes := int64(plogSizer.LogsSize(logs))
			addOutputMetricToScopeMetrics(scopeMetric, c.config.BytesMetricName, "bytes", timestamp, bytes)
		}
	}

	return c.metricsConsumer.ConsumeMetrics(ctx, outputMetrics)
}

func (c *connectorImp) ConsumeTraces(ctx context.Context, traces ptrace.Traces) error {
	outputMetrics := pmetric.NewMetrics()

	for i := 0; i < traces.ResourceSpans().Len(); i++ {
		timestamp := pcommon.NewTimestampFromTime(time.Now())
		resourceSpans := traces.ResourceSpans().At(i)

		resourceAttributes := resourceSpans.Resource().Attributes()
		rawAttributes := resourceAttributes.AsRaw()

		metricAttrMap := map[string]any{}

		metricAttrMap[dataTypeAttributeKey] = dataTypeTracesAttributeValue
		for _, key := range c.config.LabelResourceAttributes {
			if rawAttributes[key] != nil {
				metricAttrMap[key] = rawAttributes[key]
			}
		}

		resourceMetrics := outputMetrics.ResourceMetrics().AppendEmpty()
		err := resourceMetrics.Resource().Attributes().FromRaw(metricAttrMap)
		if err != nil {
			c.logger.Error("error adding attributes to datavolume metric for traces measurement", zap.Error(err), zap.Any("attributes_map", metricAttrMap))
		}
		scopeMetric := resourceMetrics.ScopeMetrics().AppendEmpty()

		if c.config.CountMetricName != "" {
			countValue := int64(resourceSpans.ScopeSpans().Len())
			addOutputMetricToScopeMetrics(scopeMetric, c.config.CountMetricName, "", timestamp, countValue)
		}

		if c.config.BytesMetricName != "" {
			bytes := int64(ptraceSizer.TracesSize(traces))
			addOutputMetricToScopeMetrics(scopeMetric, c.config.BytesMetricName, "bytes", timestamp, bytes)
		}
	}

	return c.metricsConsumer.ConsumeMetrics(ctx, outputMetrics)
}

func (c *connectorImp) ConsumeMetrics(ctx context.Context, metrics pmetric.Metrics) error {
	outputMetrics := pmetric.NewMetrics()

	for i := 0; i < metrics.ResourceMetrics().Len(); i++ {
		timestamp := pcommon.NewTimestampFromTime(time.Now())
		resourceLogs := metrics.ResourceMetrics().At(i)

		resourceAttributes := resourceLogs.Resource().Attributes()
		rawAttributes := resourceAttributes.AsRaw()

		metricAttrMap := map[string]any{}

		metricAttrMap[dataTypeAttributeKey] = dataTypeMetricsAttributeValue
		for _, key := range c.config.LabelResourceAttributes {
			if rawAttributes[key] != nil {
				metricAttrMap[key] = rawAttributes[key]
			} else {
				metricAttrMap[key] = "unknown"
			}
		}

		resourceMetrics := outputMetrics.ResourceMetrics().AppendEmpty()
		err := resourceMetrics.Resource().Attributes().FromRaw(metricAttrMap)
		if err != nil {
			c.logger.Error("error adding attributes to datavolume metric for metrics measurement", zap.Error(err), zap.Any("attributes_map", metricAttrMap))
		}
		scopeMetric := resourceMetrics.ScopeMetrics().AppendEmpty()

		if c.config.CountMetricName != "" {
			countValue := int64(resourceLogs.ScopeMetrics().Len())
			addOutputMetricToScopeMetrics(scopeMetric, c.config.CountMetricName, "", timestamp, countValue)
		}

		if c.config.BytesMetricName != "" {
			bytes := int64(pmetricSizer.MetricsSize(metrics))
			addOutputMetricToScopeMetrics(scopeMetric, c.config.BytesMetricName, "bytes", timestamp, bytes)
		}
	}

	return c.metricsConsumer.ConsumeMetrics(ctx, metrics)
}

func addOutputMetricToScopeMetrics(scopeMetric pmetric.ScopeMetrics, metricName string, unit string, timestamp pcommon.Timestamp, bytes int64) {
	metric := scopeMetric.Metrics().AppendEmpty()
	metric.SetName(metricName)
	if unit != "" {
		metric.SetUnit(unit)
	}
	sum := metric.SetEmptySum()
	sum.SetIsMonotonic(true)
	sum.SetAggregationTemporality(pmetric.AggregationTemporalityDelta)
	dataPoints := sum.DataPoints()
	dataPoint := dataPoints.AppendEmpty()
	dataPoint.SetTimestamp(timestamp)
	dataPoint.SetStartTimestamp(timestamp)
	dataPoint.SetIntValue(bytes)
}
