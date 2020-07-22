package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Output configuration
type Output struct {
	Target string `yaml:"target"`
	Format string `yaml:"format"`
}

// Rule configuration
type Rule struct {
	Name string `yaml:"name"`
	Spec string `yaml:"spec"`
}

// Healthcheck configuration
type Healthcheck struct {
	Interval string `yaml:"interval"`
	Timeout  string `yaml:"timeout"`
	Rules    []Rule `yaml:"rules"`
}

// TLS configuration
type TLS struct {
	Unsafe         bool   `yaml:"unsafe"`
	ClientCertFile string `yaml:"client_cert_file"`
	ClientKeyFile  string `yaml:"client_key_file"`
	CACertFile     string `yaml:"ca_cert_file"`
}

// Monitor configuration
type Monitor struct {
	Alias       string            `yaml:"alias"`
	Disable     bool              `yaml:"disable"`
	URL         string            `yaml:"url"`
	Method      string            `yaml:"method"`
	Headers     []string          `yaml:"headers"`
	Body        string            `yaml:"body"`
	Healthcheck Healthcheck       `yaml:"healthcheck"`
	Proxy       string            `yaml:"proxy"`
	TLS         TLS               `yaml:"tls"`
	Labels      map[string]string `yaml:"labels"`
}

// Config is the base configuration structure
type Config struct {
	Output        Output            `yaml:"output"`
	Healthcheck   Healthcheck       `yaml:"healthcheck"`
	Proxy         string            `yaml:"proxy"`
	Labels        map[string]string `yaml:"labels"`
	Monitors      []Monitor         `yaml:"monitors"`
	MonitorsFiles []string          `yaml:"monitors_files"`
}

// NewConfig create new configuration from bytes
func NewConfig(data []byte) (*Config, error) {
	data = []byte(os.ExpandEnv(string(data)))
	config := Config{}
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		err = fmt.Errorf("unable to read configuration: %s", err.Error())
		return nil, err
	}
	return &config, nil
}

// Load create new configuration object form a YAML source
func Load(configFilename string) (*Config, error) {

	// Try to load the configuration from STDIN...
	fi, err := os.Stdin.Stat()
	if err != nil {
		err = fmt.Errorf("unable to load configuration from STDIN: %s", err.Error())
		return nil, err
	}
	if fi.Mode()&os.ModeNamedPipe != 0 && fi.Size() > 0 {
		data, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			err = fmt.Errorf("unable to load configuration from STDIN: %s", err.Error())
			return nil, err
		}
		return NewConfig(data)
	}

	// Try to load configuration from file...
	data, err := ioutil.ReadFile(configFilename)
	if err != nil {
		err = fmt.Errorf("unable to load configuration from file (%s): %s", configFilename, err.Error())
		return nil, err
	}
	return NewConfig(data)
}

// MergeHealthcheckConfig merge a healthcheck configuration with another
func MergeHealthcheckConfig(a, b Healthcheck) Healthcheck {
	result := b
	if result.Interval == "" {
		result.Interval = a.Interval
	}
	if result.Timeout == "" {
		result.Timeout = a.Timeout
	}
	if len(result.Rules) == 0 {
		result.Rules = a.Rules
	}
	return result
}

// MergeLabelsConfig merge a label configuration with another
func MergeLabelsConfig(a, b map[string]string) map[string]string {
	result := make(map[string]string)
	for k, v := range a {
		result[k] = v
	}
	for k, v := range b {
		result[k] = v
	}
	return result
}
