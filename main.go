package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"github.com/ncarlier/apimon/pkg/config"
	"github.com/ncarlier/apimon/pkg/logger"
	"github.com/ncarlier/apimon/pkg/monitoring"
	"github.com/ncarlier/apimon/pkg/output"
)

// Version of the app
var Version = "snapshot"

var (
	healthy int32
	op      *output.Provider
)

var (
	version    = flag.Bool("version", false, "Print version")
	help       = flag.Bool("help", false, "Print this help screen")
	out        = flag.String("o", "", "Logging output file (default STDOUT)")
	debug      = flag.Bool("debug", false, "Activate debug logging level")
	verbose    = flag.Bool("verbose", false, "Activate verbose logging level")
	configFile = flag.String("c", "configuration.yml", "Configuration file")
)

func main() {
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *version {
		fmt.Printf(`apimon (%s)
Copyright (C) 2018 Nicolas Carlier. All rights reserved.

This work is licensed under the terms of the MIT license.  
For a copy, see <https://opensource.org/licenses/MIT>.
`, Version)
		os.Exit(0)
	}

	// Get logging level
	level := "info"
	if *debug {
		level = "debug"
	} else if *verbose {
		level = "verbose"
	}

	// Setup logger
	if *out != "" {
		f, err := os.OpenFile(*out, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			logger.Error.Panicln("unable to init the logger", err)
		}
		logger.Init(level, f, f)
		defer f.Close()
	} else {
		logger.Init(level, os.Stdout, os.Stderr)
	}

	logger.Debug.Println("starting APImon...")

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logger.Debug.Println("stoping APImon...")
		atomic.StoreInt32(&healthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := monitoring.StopMonitoring(ctx); err != nil {
			logger.Error.Fatalf("could not gracefully shutdown the daemon: %v\n", err)
		}
		if op != nil {
			op.Stop()
		}
		close(done)
	}()

	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	var confData []byte
	// Try to load the configuration...
	if fi.Mode()&os.ModeNamedPipe != 0 && fi.Size() > 0 {
		// ... form STDIN
		confData, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			logger.Error.Panicln("unable to load the configuration from stdin", err)
		}
	} else {
		// ... from file parameter
		logger.Debug.Printf("loading configuration: %s ...\n", *configFile)
		confData, err = ioutil.ReadFile(*configFile)
		if err != nil {
			logger.Error.Panicln("unable to load the configuration from file:", err)
		}
	}

	// Loading configuration...
	c, err := config.LoadConfig(confData)
	if err != nil {
		logger.Error.Panicln("unable to load the configuration", err)
	}

	// Create and start output provider
	op, err = output.NewOutputProvider(c.Output)
	if err != nil {
		logger.Error.Panicln(err)
	}
	op.Start()

	// Start all monitors...
	monitoring.StartMonitoring(*c)

	logger.Info.Println("APImon started")
	atomic.StoreInt32(&healthy, 1)

	<-done
	logger.Debug.Println("APImon stopped")
}
