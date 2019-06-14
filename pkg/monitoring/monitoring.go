package monitoring

import (
	"context"

	"github.com/ncarlier/apimon/pkg/config"
	"github.com/ncarlier/apimon/pkg/logger"
)

// Monitoring structure
type Monitoring struct {
	conf     config.Config
	monitors []*Monitor
}

// NewMonitoring creates new monitoring instance
func NewMonitoring(conf config.Config) *Monitoring {
	return &Monitoring{
		conf: conf,
	}
}

// Start is charged to start all monitors
func (m *Monitoring) Start() error {
	// Now, create all of our monitors.
	for i := 0; i < len(m.conf.Monitors); i++ {
		// Get monitor configuration
		mConfig := m.conf.Monitors[i]
		// Apply global configuration
		mConfig.Healthcheck = config.MergeHealthcheckConfig(m.conf.Healthcheck, mConfig.Healthcheck)
		// Apply proxy configuration
		if mConfig.Proxy == "" {
			mConfig.Proxy = m.conf.Proxy
		}
		// Create new monitor
		monitor, err := NewMonitor(i+1, mConfig)
		if err != nil {
			logger.Error.Println("unable to create monitor", err)
			continue
		}
		// Start the monitor
		monitor.Start()
		m.monitors = append(m.monitors, monitor)
	}
	return nil
}

// Stop is charged to stop all monitors
func (m *Monitoring) Stop(ctx context.Context) error {
	c := make(chan struct{})
	go func() {
		defer close(c)
		for _, monitor := range m.monitors {
			monitor.Stop()
		}
	}()

	select {
	case <-c:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
