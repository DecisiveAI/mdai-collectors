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

var (
	plogMarshaler = plog.ProtoMarshaler{}
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
	metrics := pmetric.NewMetrics()

	for i := 0; i < logs.ResourceLogs().Len(); i++ {
		timestamp := pcommon.NewTimestampFromTime(time.Now())
		resourceLogs := logs.ResourceLogs().At(i)

		resourceAttributes := resourceLogs.Resource().Attributes()
		rawAttributes := resourceAttributes.AsRaw()

		metricAttrMap := map[string]any{}

		metricAttrMap["data_type"] = "logs"
		for _, key := range c.config.LabelResourceAttributes {
			if rawAttributes[key] != nil {
				metricAttrMap[key] = rawAttributes[key]
			} else {
				metricAttrMap[key] = "unknown"
			}
		}

		resourceMetrics := metrics.ResourceMetrics().AppendEmpty()
		scopeMetric := resourceMetrics.ScopeMetrics().AppendEmpty()
		scopeMetric.Scope().Attributes().FromRaw(metricAttrMap) // TODO: Handle error

		if c.config.CountMetricName != "" {
			countMetric := scopeMetric.Metrics().AppendEmpty()
			countMetric.SetName(c.config.CountMetricName)
			countSum := countMetric.SetEmptySum()
			countSum.SetIsMonotonic(true)
			countSum.SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)
			countDataPoints := countSum.DataPoints()
			countDataPoint := countDataPoints.AppendEmpty()
			countDataPoint.SetTimestamp(timestamp)
			countDataPoint.SetIntValue(int64(resourceLogs.ScopeLogs().Len()))
		}

		if c.config.BytesMetricName != "" {
			bytes := plogMarshaler.LogsSize(logs)
			volumeMetric := scopeMetric.Metrics().AppendEmpty()
			volumeMetric.SetName(c.config.BytesMetricName)
			volumeMetric.SetUnit("bytes")
			volumeSum := volumeMetric.SetEmptySum()
			volumeSum.SetIsMonotonic(true)
			volumeSum.SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)
			volumeDataPoints := volumeSum.DataPoints()
			volumeDataPoint := volumeDataPoints.AppendEmpty()
			volumeDataPoint.SetTimestamp(timestamp)
			volumeDataPoint.SetIntValue(int64(bytes))
		}
	}

	return c.metricsConsumer.ConsumeMetrics(ctx, metrics)
}

// ConsumeTraces method is called for each instance of a trace sent to the connector
func (c *connectorImp) ConsumeTraces(ctx context.Context, td ptrace.Traces) error {
	//panic("ConsumeTraces is not implemented")
	return nil
}

func (c *connectorImp) ConsumeMetrics(ctx context.Context, td pmetric.Metrics) error {
	//panic("ConsumeMetrics is not implemented")
	return nil
}
