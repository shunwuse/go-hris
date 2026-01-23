package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shunwuse/go-hris/internal/http/controllers"
	"github.com/shunwuse/go-hris/internal/infra/logger"
)

// MonitorRoute struct
type MonitorRoute struct {
	logger            *logger.Logger
	monitorController *controllers.MonitorController
}

func NewMonitorRoute(
	log *logger.Logger,
	monitorController *controllers.MonitorController,
) *MonitorRoute {
	return &MonitorRoute{
		logger:            log,
		monitorController: monitorController,
	}
}

func (r *MonitorRoute) Setup(router chi.Router) {
	r.logger.Info("setting up monitor routes")

	router.Get("/health", r.monitorController.HealthCheck)

	router.Handle("/metrics", promhttp.Handler())
}
