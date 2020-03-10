package monitoring

import (
	"net/url"
	"os"
	"strconv"
	"time"

	consul "github.com/hashicorp/consul/api"
	"gopkg.in/yaml.v2"

	"github.com/ncarlier/apimon/pkg/config"
	"github.com/ncarlier/apimon/pkg/logger"
)

func (m *Monitoring) register() error {
	c, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		return err
	}

	port := 0
	address, err := os.Hostname()
	if err != nil {
		address = "127.0.0.1"
	}
	if m.conf.Output.Format == "prometheus" {
		u, err := url.ParseRequestURI(m.conf.Output.Target)
		if err == nil {
			port, _ = strconv.Atoi(u.Port())
		}
	}

	agent := c.Agent()
	serviceDef := &consul.AgentServiceRegistration{
		Name:    m.name,
		Port:    port,
		Address: address,
		Check: &consul.AgentServiceCheck{
			TTL: m.ttl.String(),
		},
	}

	logger.Debug.Println("registering service Service Registry:", serviceDef.Name, serviceDef.Address, serviceDef.Port)
	if err := agent.ServiceRegister(serviceDef); err != nil {
		return err
	}
	logger.Info.Println("service registered into the Service Registry:", serviceDef.Name)

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

func (m *Monitoring) deregister() {
	if m.agent != nil {
		m.agent.ServiceDeregister(m.name)
		m.kv = nil
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
	for {
		if m.kv == nil {
			// SD provider not register. Try again...
			time.Sleep(5 * time.Second)
			m.register()
			return
		}
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
