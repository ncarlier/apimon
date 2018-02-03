package metric

import (
	"fmt"
)

// StdoutWriter writes data to STDOUT
type StdoutWriter struct {
	Formatter Formatter
}

func newStdoutWriter(formatter Formatter) *StdoutWriter {
	return &StdoutWriter{
		Formatter: formatter,
	}
}

// Write writes metric to STDOUT
func (w *StdoutWriter) Write(metric Metric) error {
	fmt.Println(w.Formatter.Format(metric))
	return nil
}

// Close close the metric writer
func (w *StdoutWriter) Close() error {
	return nil
}
