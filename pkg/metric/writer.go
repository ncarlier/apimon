package metric

import (
	"errors"
	"fmt"
	"net/url"
)

// Writer that write array byte to a custom output
type Writer interface {
	Write(metric Metric) error
	Close() error
}

func getMetricWriter(target, format string) (Writer, error) {
	formatter, err := getMetricFormatter(format)
	if err != nil {
		return nil, fmt.Errorf("unable to get metric formatter: %s", err)
	}
	var writer Writer
	switch {
	case target == "", target == "stdout":
		writer = newStdoutWriter(formatter)
	case isValidURLWithScheme(target, "http") || isValidURLWithScheme(target, "https"):
		writer = newHTTPWriter(target, formatter)
	case isValidURLWithScheme(target, "file"):
		writer, err = newFileWriter(target, formatter)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("non supported metric writer: " + target)
	}
	return writer, nil
}

func isValidURLWithScheme(toTest, scheme string) bool {
	u, err := url.ParseRequestURI(toTest)
	if err != nil || u.Scheme != scheme {
		return false
	}
	return true
}
