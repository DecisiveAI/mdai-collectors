package datavolumeconnector

import (
	"fmt"
)

type Config struct {
	LabelResourceAttributes []string `mapstructure:"label_resource_attributes"`
	BytesMetricName         string   `mapstructure:"bytes_metric_name"`
	CountMetricName         string   `mapstructure:"count_metric_name"`
}

func (c *Config) Validate() error {
	if c.BytesMetricName == "" && c.CountMetricName == "" {
		return fmt.Errorf("one of bytes_metric_name and/or count_metric_name must be specified")
	}
	return nil
}
