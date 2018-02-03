package metric

import "errors"

// Formatter that format a metric to a string
type Formatter interface {
	Format(metric Metric) string
	ContentType() string
}

func getMetricFormatter(format string) (Formatter, error) {
	switch format {
	case "", "influxdb":
		return &InfluxDBMetricFormatter{}, nil
	case "json":
		return &JSONMetricFormatter{}, nil
	default:
		return nil, errors.New("non supported format: " + format)
	}
}
