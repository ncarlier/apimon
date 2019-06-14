package logger

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	// Debug level
	Debug *log.Logger
	// Info level
	Info *log.Logger
	// Warning level
	Warning *log.Logger
	// Error level
	Error *log.Logger
)

func init() {
	Configure("warn", "")
}

// Configure logger
func Configure(level, output string) error {
	outWriter := os.Stdout
	errWriter := os.Stderr
	if output != "" {
		outWriter, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		errWriter = outWriter
	}

	var debugHandle, infoHandle, warnHandle, errorHandle io.Writer
	debugHandle = outWriter
	infoHandle = outWriter
	warnHandle = errWriter
	errorHandle = errWriter
	switch level {
	case "info":
		debugHandle = ioutil.Discard
	case "warn":
		debugHandle = ioutil.Discard
		infoHandle = ioutil.Discard
	case "error":
		debugHandle = ioutil.Discard
		infoHandle = ioutil.Discard
		warnHandle = ioutil.Discard
	}

	Debug = log.New(debugHandle, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(infoHandle, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(warnHandle, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(errorHandle, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	return nil
}
