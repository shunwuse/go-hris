package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	registry *prometheus.Registry

	HttpRequestsTotal   *prometheus.CounterVec
	HttpRequestDuration *prometheus.HistogramVec
}

func New() *Metrics {
	registry := prometheus.NewRegistry()
	factory := promauto.With(registry)

	metrics := &Metrics{
		registry: registry,
		HttpRequestsTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		HttpRequestDuration: factory.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
	}

	registry.MustRegister(collectors.NewGoCollector())
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	return metrics
}

// Handler exposes the dedicated Prometheus registry for HTTP scraping.
func (m *Metrics) Handler() http.Handler {
	if m == nil || m.registry == nil {
		return promhttp.Handler()
	}

	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}
