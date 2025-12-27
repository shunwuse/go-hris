package middlewares

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/shunwuse/go-hris/internal/infra"
)

type MetricsMiddleware struct {
	metrics *infra.Metrics
}

func NewMetricsMiddleware(metrics *infra.Metrics) *MetricsMiddleware {
	return &MetricsMiddleware{
		metrics: metrics,
	}
}

func (m *MetricsMiddleware) Handler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			duration := time.Since(start).Seconds()
			status := strconv.Itoa(ww.Status())
			path := m.getPattern(r)

			m.metrics.HttpRequestsTotal.WithLabelValues(r.Method, path, status).Inc()
			m.metrics.HttpRequestDuration.WithLabelValues(r.Method, path).Observe(duration)
		})
	}
}

func (m *MetricsMiddleware) getPattern(r *http.Request) string {
	rctx := chi.RouteContext(r.Context())
	if pattern := rctx.RoutePattern(); pattern != "" {
		return pattern
	}
	return r.URL.Path
}

func (m *MetricsMiddleware) Setup(router chi.Router) {
	router.Use(m.Handler())
}
