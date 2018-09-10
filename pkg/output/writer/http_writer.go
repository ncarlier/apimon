package writer

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/ncarlier/apimon/pkg/model"
	"github.com/ncarlier/apimon/pkg/output/format"
)

// HTTPWriter HTTP writer
type HTTPWriter struct {
	Formatter format.Formatter
	URL       string
}

func newHTTPWriter(url string, formatter format.Formatter) *HTTPWriter {
	return &HTTPWriter{
		Formatter: formatter,
		URL:       url,
	}
}

// Write post metric to HTTP endpoint
func (w *HTTPWriter) Write(metric model.Metric) error {
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
