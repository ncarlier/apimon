package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/ncarlier/apimon/pkg/config"
	"github.com/ncarlier/apimon/pkg/logger"
	"github.com/ncarlier/apimon/pkg/monitoring"
	"github.com/ncarlier/apimon/pkg/output"
)

var (
	healthy int32
	op      *output.Provider
)

var (
	version    = flag.Bool("version", false, "Print version")
	help       = flag.Bool("help", false, "Print this help screen")
	out        = flag.String("o", "", "Logging output file (default STDOUT)")
	debug      = flag.Bool("vv", false, "Activate debug logging level")
	verbose    = flag.Bool("v", false, "Activate verbose logging level")
	configFile = flag.String("c", "./apimon.yml", "Configuration file (can also be provided with STDIN)")
)

func main() {
	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *version {
		printVersion()
		os.Exit(0)
	}

	// Get logging level
	level := "warn"
	if *debug {
		level = "debug"
	} else if *verbose {
		level = "info"
	}

	// Setup logger
	if err := logger.Configure(level, *out); err != nil {
		log.Panicln("unable to init the logger", err)
	}

	logger.Debug.Println("starting APImon...")

	// Loading configuration...
	c, err := config.Load(*configFile)
	if err != nil {
		logger.Error.Panicln("unable to load the configuration", err)
	}

	// Create and start output provider...
	op, err = output.NewOutputProvider(c.Output)
	if err != nil {
		logger.Error.Panicln(err)
	}
	op.Start()

	// Create and start monitoring...
	mon := monitoring.NewMonitoring(*c)
	mon.Start()
	go mon.RestartOnSDConfigChange()

	// Graceful shutdown
	go func() {
		<-quit
		logger.Debug.Println("stoping APImon...")
		atomic.StoreInt32(&healthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := mon.Stop(ctx); err != nil {
			logger.Error.Fatalf("could not gracefully shutdown the daemon: %v\n", err)
		}
		if op != nil {
			op.Stop()
		}
		close(done)
	}()

	logger.Info.Println("APImon started")
	atomic.StoreInt32(&healthy, 1)

	<-done
	logger.Debug.Println("APImon stopped")
}
