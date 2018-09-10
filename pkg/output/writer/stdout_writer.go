package writer

import (
	"fmt"

	"github.com/ncarlier/apimon/pkg/model"
	"github.com/ncarlier/apimon/pkg/output/format"
)

// StdoutWriter writes data to STDOUT
type StdoutWriter struct {
	Formatter format.Formatter
}

func newStdoutWriter(formatter format.Formatter) *StdoutWriter {
	return &StdoutWriter{
		Formatter: formatter,
	}
}

// Write writes metric to STDOUT
func (w *StdoutWriter) Write(metric model.Metric) error {
	fmt.Println(w.Formatter.Format(metric))
	return nil
}

// Close close the metric writer
func (w *StdoutWriter) Close() error {
	return nil
}
