package format

import (
	"github.com/ncarlier/apimon/pkg/model"
)

func newStringMetricFormatter() (*StringMetricFormatter, error) {
	return &StringMetricFormatter{}, nil
}

// StringMetricFormatter String metric formatter
type StringMetricFormatter struct{}

// Format a metric to raw string
func (f *StringMetricFormatter) Format(metric model.Metric) string {
	return metric.String()
}

// ContentType gets formatter content-type
func (f *StringMetricFormatter) ContentType() string {
	return "text/plain"
}
