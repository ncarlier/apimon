package writer

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ncarlier/apimon/pkg/logger"
	"github.com/ncarlier/apimon/pkg/model"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var healthCheckOKCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_health_check_ok",
		Help: "HTTP health check OK.",
	},
	[]string{"name"},
)
var healthCheckKOCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_health_check_ko",
		Help: "HTTP health check kO.",
	},
	[]string{"name", "reason"},
)
var healthCheckResponseTimeGauge = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "http_health_check_response_time",
		Help: "HTTP health check response time.",
	},
	[]string{"name"},
)

// PrometheusWriter writes data to Prometheus endpoint
type PrometheusWriter struct {
	srv *http.Server
}

func newPrometheusWriter(uri string) (*PrometheusWriter, error) {
	u, err := url.ParseRequestURI(uri)
	if err != nil || u.Scheme != "http" {
		return nil, fmt.Errorf("invalid listen URL: %s", uri)
	}
	prometheus.MustRegister(healthCheckOKCounter)
	prometheus.MustRegister(healthCheckKOCounter)
	prometheus.MustRegister(healthCheckResponseTimeGauge)
	srv := &http.Server{Addr: u.Hostname() + ":" + u.Port()}
	http.Handle(u.Path, promhttp.Handler())
	go func() {
		logger.Debug.Printf("starting HTTP server (%s)...\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error.Panicln("unable to create Prometheus server endpoint:", err)
		}
	}()
	return &PrometheusWriter{
		srv: srv,
	}, nil
}

// Write writes metric to Prometheus
func (w *PrometheusWriter) Write(metric model.Metric) error {
	duration := float64(metric.Duration / time.Millisecond)
	reason := ""
	if metric.Error != "" {
		reason = strings.SplitN(metric.Error, ":", 2)[0]
		reason = strings.ToLower(reason)
		healthCheckKOCounter.With(prometheus.Labels{
			"name":   metric.Name,
			"reason": reason,
		}).Inc()
	} else {
		healthCheckOKCounter.With(prometheus.Labels{
			"name": metric.Name,
		}).Inc()
	}
	healthCheckResponseTimeGauge.With(prometheus.Labels{
		"name": metric.Name,
	}).Set(duration)

	return nil
}

// Close close the metric writer
func (w *PrometheusWriter) Close() error {
	logger.Debug.Printf("stopping HTTP server (%s)...\n", w.srv.Addr)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return w.srv.Shutdown(ctx)
}
