package writer

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ncarlier/apimon/pkg/model"
	"github.com/ncarlier/apimon/pkg/output/format"
)

// FileWriter writes metric to file
type FileWriter struct {
	Formatter format.Formatter
	File      *os.File
	Writer    *bufio.Writer
}

func newFileWriter(target string, formatter format.Formatter) (*FileWriter, error) {
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
func (w *FileWriter) Write(metric model.Metric) error {
	fmt.Fprintln(w.Writer, w.Formatter.Format(metric))
	return nil
}

// Close close the metric writer
func (w *FileWriter) Close() error {
	defer w.File.Close()
	return w.Writer.Flush()
}
