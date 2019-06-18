package monitoring

import (
	"context"
	"time"

	consul "github.com/hashicorp/consul/api"
	"gopkg.in/yaml.v2"

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
		logger.Error.Println("unable to register monitoring service", err)
	}

	return mon
}

func (m *Monitoring) register() error {
	c, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		return err
	}
	agent := c.Agent()
	serviceDef := &consul.AgentServiceRegistration{
		Name: m.name,
		Check: &consul.AgentServiceCheck{
			TTL: m.ttl.String(),
		},
	}

	if err := agent.ServiceRegister(serviceDef); err != nil {
		return err
	}
	logger.Debug.Println("service registered into the Service Registry:", m.name)

	go func() {
		m.ping()
		ticker := time.NewTicker(m.ttl / 2)
		for range ticker.C {
			m.ping()
		}
	}()

	m.agent = agent
	m.kv = c.KV()

	return nil
}

func (m *Monitoring) unregister() {
	if m.agent != nil {
		m.agent.ServiceDeregister(m.name)
	}
}

func (m *Monitoring) ping() {
	if m.err != nil {
		if err := m.agent.FailTTL("service:"+m.name, m.err.Error()); err != nil {
			logger.Error.Println("unable to ping service registry", err)
		}
	} else if err := m.agent.PassTTL("service:"+m.name, ""); err != nil {
		logger.Error.Println("unable to ping service registry", err)
	}
}

func (m *Monitoring) getSDConfigKey() string {
	return m.name + "/monitors"
}

func (m *Monitoring) getSDConfig() ([]config.Monitor, error) {
	result := []config.Monitor{}
	if m.kv == nil {
		return result, nil
	}
	opts := &consul.QueryOptions{RequireConsistent: true}
	pair, meta, err := m.kv.Get(m.getSDConfigKey(), opts)
	if err != nil {
		return result, err
	}
	if pair == nil || meta == nil {
		return result, nil
	}
	if err = yaml.Unmarshal(pair.Value, &result); err != nil {
		return result, err
	}
	m.idx = meta.LastIndex
	return result, nil
}

// RestartOnSDConfigChange restart monitoring when Service Discover configuration change
func (m *Monitoring) RestartOnSDConfigChange() {
	if m.kv == nil {
		return
	}
	for {
		pair, meta, err := m.kv.Get(m.getSDConfigKey(), &consul.QueryOptions{
			WaitIndex: m.idx,
		})
		if err != nil {
			logger.Error.Println("read error from Service Discover KV store:", err)
			break
		}
		if pair == nil || meta == nil {
			// Query won’t be blocked if key not found
			time.Sleep(1 * time.Second)
		} else {
			logger.Debug.Println("configuration changed: reloading...")
			m.idx = meta.LastIndex
			m.restart()
		}
	}
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
	monitors = append(m.conf.Monitors, monitors...)
	// Now, create all of our monitors.
	for i := 0; i < len(monitors); i++ {
		// Get monitor configuration
		mConfig := monitors[i]
		if mConfig.Disable {
			continue
		}
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
	m.err = nil
	return nil
}

// Stop is charged to stop all monitors
func (m *Monitoring) Stop(ctx context.Context) error {
	logger.Debug.Println("stopping monitoring...")
	defer m.unregister()
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
