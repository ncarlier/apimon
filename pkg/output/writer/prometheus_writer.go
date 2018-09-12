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

var healthCheckSummary = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Name:       "http_health_check",
		Help:       "HTTP health check status.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
	[]string{"name", "status", "reason"},
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
	prometheus.MustRegister(healthCheckSummary)
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
	status := "0"
	if metric.Status == "UP" {
		status = "1"
	}
	duration := float64(metric.Duration / time.Millisecond)
	if metric.Error != "" {
		reason := strings.SplitN(metric.Error, ":", 2)[0]
		reason = strings.ToLower(reason)
		healthCheckSummary.With(prometheus.Labels{
			"name":   metric.Name,
			"status": status,
			"reason": reason,
		}).Observe(duration)
	} else {
		healthCheckSummary.With(prometheus.Labels{
			"name":   metric.Name,
			"status": status,
			"reason": "",
		}).Observe(duration)
	}
	return nil
}

// Close close the metric writer
func (w *PrometheusWriter) Close() error {
	logger.Debug.Printf("stopping HTTP server (%s)...\n", w.srv.Addr)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return w.srv.Shutdown(ctx)
}
