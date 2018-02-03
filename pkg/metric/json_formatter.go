package metric

import (
	"encoding/json"
	"fmt"
)

// JSONMetricFormatter JSON metric formatter
type JSONMetricFormatter struct{}

// Format a metric to a JSON string
func (f *JSONMetricFormatter) Format(metric Metric) string {
	b, err := json.Marshal(metric)
	if err != nil {
		return fmt.Sprintf("{\"error\": \"%s\"}", err)
	}
	return string(b)
}

// ContentType gets formatter content-type
func (f *JSONMetricFormatter) ContentType() string {
	return "application/json"
}
