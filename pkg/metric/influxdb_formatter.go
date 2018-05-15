package metric

import (
	"fmt"
	"strings"
	"time"
)

// InfluxDBMetricFormatter InfluxDB metric formatter
type InfluxDBMetricFormatter struct{}

// Format a metric to a InfluxDB string
func (f *InfluxDBMetricFormatter) Format(metric Metric) string {
	var status int8
	if metric.Status == "UP" {
		status = 1
	}
	duration := int64(metric.Duration / time.Millisecond)
	ts := metric.Timestamp.UnixNano()
	if metric.Error != "" {
		reason := strings.SplitN(metric.Error, ":", 2)[0]
		reason = strings.ToLower(reason)
		return fmt.Sprintf(
			"http_health_check,name=%s value=%d,reason=\"%s\",duration=%d %d",
			metric.Name,
			status,
			reason,
			duration,
			ts)
	}
	return fmt.Sprintf(
		"http_health_check,name=%s value=%d,duration=%d %d",
		metric.Name,
		status,
		duration,
		ts)
}

// ContentType gets formatter content-type
func (f *InfluxDBMetricFormatter) ContentType() string {
	return "text/plain"
}
