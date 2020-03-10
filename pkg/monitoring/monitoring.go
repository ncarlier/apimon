package monitoring

import (
	"context"
	"time"

	consul "github.com/hashicorp/consul/api"

	"github.com/ncarlier/apimon/pkg/config"
	"github.com/ncarlier/apimon/pkg/logger"
)

// Monitoring structure
type Monitoring struct {
	name     string
	conf     config.Config
	monitors []*Monitor
	agent    *consul.Agent
	kv       *consul.KV
	ttl      time.Duration
	idx      uint64
	err      error
}

// NewMonitoring creates new monitoring instance
func NewMonitoring(conf config.Config) *Monitoring {
	mon := &Monitoring{
		name: "apimon",
		conf: conf,
		ttl:  30 * time.Second,
	}
	// Register service into SD registry
	if err := mon.register(); err != nil {
		logger.Warning.Println("unable to register service to SD registry", err)
	}

	return mon
}

// Start is charged to start all monitors
func (m *Monitoring) Start() error {
	logger.Debug.Println("starting monitoring...")
	// Get configuration from Service Discovery K/V store
	monitors, err := m.getSDConfig()
	if err != nil {
		m.err = err
		return err
	}
	monitors = append(monitors, m.getFilesConfig()...)
	monitors = append(monitors, m.conf.Monitors...)
	// Now, create all of our monitors.
	for i := 0; i < len(monitors); i++ {
		// Get monitor configuration
		mConfig := monitors[i]
		if mConfig.Disable {
			continue
		}
		// Apply global configuration
		mConfig.Healthcheck = config.MergeHealthcheckConfig(m.conf.Healthcheck, mConfig.Healthcheck)
		mConfig.Labels = config.MergeLabelsConfig(m.conf.Labels, mConfig.Labels)
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
	m.err = nil
	return nil
}

// Stop is charged to stop all monitors
func (m *Monitoring) Stop(ctx context.Context) error {
	logger.Debug.Println("stopping monitoring...")
	defer m.deregister()
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

func (m *Monitoring) restart() error {
	for _, monitor := range m.monitors {
		monitor.Stop()
	}
	m.monitors = m.monitors[:0]
	return m.Start()
}
