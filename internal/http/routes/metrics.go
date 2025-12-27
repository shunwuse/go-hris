package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsRoute struct{}

func NewMetricsRoute() *MetricsRoute {
	return &MetricsRoute{}
}

func (r *MetricsRoute) Setup(router chi.Router) {
	router.Handle("/metrics", promhttp.Handler())
}
