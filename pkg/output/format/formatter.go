package format

import (
	"errors"

	"github.com/ncarlier/apimon/pkg/model"
)

// Formatter that format a metric to a string
type Formatter interface {
	Format(metric model.Metric) string
	ContentType() string
}

// NewMetricFormatter creats new metric formatter
func NewMetricFormatter(format string) (Formatter, error) {
	switch format {
	case "", "influxdb":
		return newInfluxDBMetricFormatter()
	case "json":
		return newJSONMetricFormatter()
	default:
		return nil, errors.New("non supported format: " + format)
	}
}
