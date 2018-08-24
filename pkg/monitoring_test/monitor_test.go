package monitoring_test

import (
	"sync"
	"testing"
	"time"

	"github.com/ncarlier/apimon/pkg/assert"
	"github.com/ncarlier/apimon/pkg/config"
	"github.com/ncarlier/apimon/pkg/monitoring"
)

var stop = make(chan struct{})
var wg sync.WaitGroup

func TestMonitorWithBadURLConfiguration(t *testing.T) {
	conf := &config.Monitor{
		URL: "foo",
		Healthcheck: config.Healthcheck{
			Rules: []config.Rule{
				config.Rule{Name: "code", Spec: "200"},
			},
		},
	}

	expected := "parse foo: invalid URI for request"
	_, err := monitoring.NewMonitor(0, *conf, stop, &wg)
	assert.NotNil(t, err, "Monitor creation should fail")
	assert.Equal(t, expected, err.Error(), "Unexpected error")
}

func TestMonitorWithDefaultConfiguration(t *testing.T) {
	expectedInterval := time.Duration(30) * time.Second
	expectedTimeout := time.Duration(5) * time.Second
	conf := &config.Monitor{
		URL: "http://foo",
	}

	monitor, err := monitoring.NewMonitor(0, *conf, stop, &wg)
	assert.Nil(t, err, "Monitor creation should not fail")
	assert.NotNil(t, monitor, "Monitor should be created")
	assert.Equal(t, expectedTimeout, monitor.Timeout, "Unexpected monitor timeout")
	assert.Equal(t, expectedInterval, monitor.Interval, "Unexpected monitor timeout")
	assert.Equal(t, 1, len(monitor.Validators), "Unexpected number of validators")
	assert.Equal(t, "code", monitor.Validators[0].Name, "Unexpected validator name")
	assert.Equal(t, "200", monitor.Validators[0].Spec, "Unexpected validator spec")
}

func TestMonitorWithAdjustedConfiguration(t *testing.T) {
	expectedInterval := time.Duration(2) * time.Second
	expectedTimeout := time.Duration(1900) * time.Millisecond
	conf := &config.Monitor{
		URL: "http://foo",
		Healthcheck: config.Healthcheck{
			Interval: "2s",
			Timeout:  "2s",
			Rules: []config.Rule{
				config.Rule{Name: "code", Spec: "200"},
			},
		},
	}

	monitor, err := monitoring.NewMonitor(0, *conf, stop, &wg)
	assert.Nil(t, err, "Monitor creation should not fail")
	assert.NotNil(t, monitor, "Monitor should be created")
	assert.Equal(t, expectedTimeout, monitor.Timeout, "Unexpected monitor timeout")
	assert.Equal(t, expectedInterval, monitor.Interval, "Unexpected monitor timeout")
}
