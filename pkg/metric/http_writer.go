package metric

import (
	"bytes"
	"fmt"
	"net/http"
)

// HTTPWriter HTTP writer
type HTTPWriter struct {
	Formatter Formatter
	URL       string
}

func newHTTPWriter(url string, formatter Formatter) *HTTPWriter {
	return &HTTPWriter{
		Formatter: formatter,
		URL:       url,
	}
}

// Write post metric to HTTP endpoint
func (w *HTTPWriter) Write(metric Metric) error {
	contentType := w.Formatter.ContentType()
	body := w.Formatter.Format(metric)
	resp, err := http.Post(w.URL, contentType, bytes.NewBufferString(body))
	if err != nil {
		return err
	} else if resp.StatusCode >= 300 {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}
	return nil
}

// Close close the metric writer
func (w *HTTPWriter) Close() error {
	return nil
}
