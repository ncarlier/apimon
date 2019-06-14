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

// Monitor configuration
type Monitor struct {
	Alias       string      `yaml:"alias"`
	URL         string      `yaml:"url"`
	Headers     []string    `yaml:"headers"`
	Healthcheck Healthcheck `yaml:"healthcheck"`
	Proxy       string      `yaml:"proxy"`
	Unsafe      bool        `yaml:"unsafe"`
}

// Config is the base configuration structure
type Config struct {
	Output      Output      `yaml:"output"`
	Healthcheck Healthcheck `yaml:"healthcheck"`
	Proxy       string      `yaml:"proxy"`
	MonitorsURI string      `yaml:"monitors_uri"`
	Monitors    []Monitor   `yaml:"monitors"`
}

func newConfig(data []byte) (*Config, error) {
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
		return newConfig(data)
	}

	// Try to load configuration from file...
	data, err := ioutil.ReadFile(configFilename)
	if err != nil {
		err = fmt.Errorf("unable to load configuration from file (%s): %s", configFilename, err.Error())
		return nil, err
	}
	return newConfig(data)
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
