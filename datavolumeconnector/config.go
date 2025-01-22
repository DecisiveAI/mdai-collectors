package datavolumeconnector

import (
	"fmt"
)

type Config struct {
	// Resource attributes that will be extracted from resources and appended to output metrics
	LabelResourceAttributes []string `mapstructure:"label_resource_attributes"`
	// The name of the bytes measurement metric name. Required if count_metric_name is not present. Byte measurement will not occur if this is not present.
	BytesMetricName string `mapstructure:"bytes_metric_name"`
	// The name of the scope item measurement metric name. Required if bytes_metric_name is not present. Scope item measurement will not occur if this is not present.
	CountMetricName string `mapstructure:"count_metric_name"`
}

func (c *Config) Validate() error {
	if c.BytesMetricName == "" && c.CountMetricName == "" {
		return fmt.Errorf("one of bytes_metric_name and/or count_metric_name must be specified")
	}
	return nil
}
