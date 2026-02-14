package routes

import (
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shunwuse/go-hris/internal/http/controllers"
	"github.com/shunwuse/go-hris/internal/http/middlewares"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
)

type MonitorRoute struct {
	logger             *logger.Logger
	profilerMiddleware *middlewares.ProfilerMiddleware
	monitorController  *controllers.MonitorController
}

func NewMonitorRoute(
	log *logger.Logger,
	profMiddleware *middlewares.ProfilerMiddleware,
	monitorController *controllers.MonitorController,
) *MonitorRoute {
	return &MonitorRoute{
		logger:             log,
		profilerMiddleware: profMiddleware,
		monitorController:  monitorController,
	}
}

func (r *MonitorRoute) Setup(router chi.Router) {
	r.logger.Info("setting up monitor routes")

	router.Get("/health", r.monitorController.HealthCheck)

	router.Handle("/metrics", promhttp.Handler())

	// Mount pprof routes.
	router.Route("/debug/pprof", func(router chi.Router) {
		router.Use(r.profilerMiddleware.Handler())

		router.HandleFunc("/", pprof.Index)
		router.HandleFunc("/cmdline", pprof.Cmdline)
		router.HandleFunc("/profile", pprof.Profile)
		router.HandleFunc("/symbol", pprof.Symbol)
		router.HandleFunc("/trace", pprof.Trace)

		// Also handle sub-paths like /heap, /goroutine, etc.
		router.Handle("/{name}", http.HandlerFunc(pprof.Index))
	})
}
