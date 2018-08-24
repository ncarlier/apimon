package monitoring

import (
	"context"
	"sync"

	"github.com/ncarlier/apimon/pkg/config"
	"github.com/ncarlier/apimon/pkg/logger"
)

var stop = make(chan struct{})
var wg sync.WaitGroup

// StartMonitoring is charged to start all monitors.
func StartMonitoring(conf config.Config) error {
	// Now, create all of our monitors.
	for i := 0; i < len(conf.Monitors); i++ {
		// Get monitor configuration
		mConfig := conf.Monitors[i]
		// Apply global configuration
		mConfig.Healthcheck = config.MergeHealthcheckConfig(conf.Healthcheck, mConfig.Healthcheck)
		// Apply proxy configuration
		if mConfig.Proxy == "" {
			mConfig.Proxy = conf.Proxy
		}
		// Create new monitor
		worker, err := NewMonitor(i+1, mConfig, stop, &wg)
		if err != nil {
			logger.Error.Println("Unable to create monitor", err)
			continue
		}
		// Start the monitor
		worker.Start()
		wg.Add(1)
	}
	return nil
}

// StopMonitoring is charged to stop all monitors.
func StopMonitoring(ctx context.Context) error {
	close(stop)
	wg.Wait()
	return nil
}
