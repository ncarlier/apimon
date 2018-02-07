package config

import (
	"gopkg.in/yaml.v2"
)

// Output configuration of the output channel
type Output struct {
	Target string `yaml:"target"`
	Format string `yaml:"format"`
}

// Healthcheck configuration of a healthcheck
type Healthcheck struct {
	Interval string `yaml:"interval"`
	Timeout  string `yaml:"timeout"`
	Rules    string `yaml:"rules"`
}

// Monitor configuration of a monitor
type Monitor struct {
	Alias       string      `yaml:"alias"`
	URL         string      `yaml:"url"`
	Headers     []string    `yaml:"headers"`
	Healthcheck Healthcheck `yaml:"healthcheck"`
}

// Config configuration structure
type Config struct {
	Output      Output      `yaml:"output"`
	Healthcheck Healthcheck `yaml:"healthcheck"`
	Monitors    []Monitor   `yaml:"monitors"`
}

// LoadConfig create new configuration object form a YAML file
func LoadConfig(data []byte) (*Config, error) {
	config := Config{}
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
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
	if result.Rules == "" {
		result.Rules = a.Rules
	}
	return result
}
