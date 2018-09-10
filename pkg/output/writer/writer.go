package writer

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/ncarlier/apimon/pkg/model"
	"github.com/ncarlier/apimon/pkg/output/format"
)

// Writer that write array byte to a custom output
type Writer interface {
	Write(metric model.Metric) error
	Close() error
}

func NewOutputWriter(target, _format string) (Writer, error) {
	formatter, err := format.NewMetricFormatter(_format)
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
		return nil, errors.New("unsupported output writer: " + target)
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
