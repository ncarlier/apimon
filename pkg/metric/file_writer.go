package metric

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// FileWriter writes metric to file
type FileWriter struct {
	Formatter Formatter
	File      *os.File
	Writer    *bufio.Writer
}

func newFileWriter(target string, formatter Formatter) (*FileWriter, error) {
	file, err := os.Create(strings.TrimPrefix(target, "file://"))
	if err != nil {
		return nil, err
	}
	writer := bufio.NewWriter(file)

	return &FileWriter{
		Formatter: formatter,
		File:      file,
		Writer:    writer,
	}, nil
}

// Write writes metric to file
func (w *FileWriter) Write(metric Metric) error {
	fmt.Fprintln(w.Writer, w.Formatter.Format(metric))
	return nil
}

// Close close the metric writer
func (w *FileWriter) Close() error {
	defer w.File.Close()
	return w.Writer.Flush()
}
