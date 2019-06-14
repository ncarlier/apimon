package monitoring

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/ncarlier/apimon/pkg/config"
	"github.com/ncarlier/apimon/pkg/logger"
	"github.com/ncarlier/apimon/pkg/model"
	"github.com/ncarlier/apimon/pkg/output"
	"github.com/ncarlier/apimon/pkg/rule"
)

var defaultDuration = time.Duration(30) * time.Second
var defaultTimeout = time.Duration(5) * time.Second

// Monitor is a go routine in charge of the monitoring work.
type Monitor struct {
	ID         int
	Alias      string
	URL        url.URL
	Headers    []string
	Client     *http.Client
	Interval   time.Duration
	Timeout    time.Duration
	Validators []rule.Validator
	Ticker     *time.Ticker
}

// NewMonitor create a new monitor
func NewMonitor(id int, conf config.Monitor) (*Monitor, error) {
	// Parse the interval...
	interval, err := time.ParseDuration(conf.Healthcheck.Interval)
	if err != nil {
		logger.Warning.Printf("unable to parse healthcheck interval: '%s'. Using default: %s", conf.Healthcheck.Interval, defaultDuration)
		interval = defaultDuration
	}

	// Parse the timeout...
	timeout, err := time.ParseDuration(conf.Healthcheck.Timeout)
	if err != nil {
		logger.Warning.Printf("unable to parse timeout: '%s'. Using default: %s", conf.Healthcheck.Timeout, defaultTimeout)
		timeout = defaultTimeout
	}
	if timeout >= interval {
		logger.Warning.Printf("timeout can't be longer than the interval: %s > %s. Adjusting timeout.", timeout, interval)
		timeout = interval - time.Duration(100)*time.Millisecond
	}

	// Parse the URL
	u, err := url.ParseRequestURI(conf.URL)
	if err != nil {
		logger.Error.Printf("unable to parse URL: '%s'", conf.URL)
		return nil, err
	}

	// Create HTTP client
	client := &http.Client{}
	transport := &http.Transport{}
	if conf.Proxy != "" {
		proxyURL, err := url.ParseRequestURI(conf.Proxy)
		if err != nil {
			logger.Error.Printf("unable to parse proxy URL: '%s'", conf.Proxy)
			return nil, err
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	}
	if conf.Unsafe {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	client.Transport = transport

	// Parse validators
	validators, err := rule.CreateValidatorPipeline(conf.Healthcheck.Rules)
	if err != nil {
		logger.Error.Printf("unable to parse healthcheck rules: '%s'", conf.Healthcheck.Rules)
		return nil, err
	}

	// Create, and return the monitor.
	monitor := Monitor{
		ID:         id,
		Alias:      conf.Alias,
		URL:        *u,
		Client:     client,
		Interval:   interval,
		Timeout:    timeout,
		Validators: validators,
		Headers:    conf.Headers,
	}
	logger.Debug.Printf("monitor created: %s\n", monitor)

	return &monitor, nil
}

// String to string convertion
func (m Monitor) String() string {
	return fmt.Sprintf(
		"{id: %d, alias: \"%s\", url: \"%s\", interval: \"%s\", timeout: \"%s\"}",
		m.ID,
		m.Alias,
		m.URL.String(),
		m.Interval,
		m.Timeout)
}

// Start the monitor
func (m *Monitor) Start() {
	logger.Debug.Printf("starting monitor %s#%d...\n", m.Alias, m.ID)
	m.Ticker = time.NewTicker(m.Interval)
	go func() {
		for range m.Ticker.C {
			var name string
			if m.Alias != "" {
				name = m.Alias
			} else {
				name = m.URL.String()
			}
			_metric := &model.Metric{
				Name:      name,
				Status:    "UP",
				Timestamp: time.Now(),
			}
			var err error
			_metric.Duration, err = m.Validate()
			if err != nil {
				_metric.Status = "DOWN"
				_metric.Error = err.Error()
			}
			logger.Debug.Printf("monitor %s#%d: %s\n", m.Alias, m.ID, _metric)
			output.Queue <- *_metric
		}
	}()
}

// Stop the monitor
func (m *Monitor) Stop() {
	m.Ticker.Stop()
	logger.Debug.Printf("monitor %s#%d stopped\n", m.Alias, m.ID)
}

// Validate the monitor endpoint by aplying all validators
func (m *Monitor) Validate() (time.Duration, error) {
	start := time.Now()
	ctx, cancel := context.WithCancel(context.TODO())
	timer := time.AfterFunc(m.Timeout, func() {
		cancel()
	})

	req, err := http.NewRequest("GET", m.URL.String(), nil)
	if err != nil {
		return time.Since(start), fmt.Errorf("PREPARE_REQUEST: %s", err)
	}
	req = req.WithContext(ctx)

	for _, header := range m.Headers {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) == 2 {
			req.Header.Add(parts[0], parts[1])
		}
	}

	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "APImonPinger/1.0")
	}

	resp, err := m.Client.Do(req)
	if err != nil {
		matched, _ := regexp.MatchString("context canceled", err.Error())
		if matched {
			return time.Since(start), fmt.Errorf("TIMEOUT: %s", err)
		}
		return time.Since(start), fmt.Errorf("REQUEST: %s", err)
	}
	defer resp.Body.Close()
	timer.Stop()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return time.Since(start), fmt.Errorf("BODY: %s", err)
	}

	for _, validator := range m.Validators {
		if err = validator.Validate(string(body), resp); err != nil {
			return time.Since(start), fmt.Errorf("RULE_%s: %s", strings.ToUpper(validator.Name()), err)
		}
	}

	return time.Since(start), nil
}
