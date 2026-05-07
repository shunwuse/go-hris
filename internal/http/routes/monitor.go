package routes

import (
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/shunwuse/go-hris/internal/http/controllers"
	"github.com/shunwuse/go-hris/internal/http/middlewares"
	infrmetrics "github.com/shunwuse/go-hris/internal/infra/metrics"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
)

type MonitorRoute struct {
	logger             *logger.Logger
	metrics            *infrmetrics.Metrics
	profilerMiddleware *middlewares.ProfilerMiddleware
	monitorController  *controllers.MonitorController
}

func NewMonitorRoute(
	log *logger.Logger,
	metrics *infrmetrics.Metrics,
	profMiddleware *middlewares.ProfilerMiddleware,
	monitorController *controllers.MonitorController,
) *MonitorRoute {
	return &MonitorRoute{
		logger:             log,
		metrics:            metrics,
		profilerMiddleware: profMiddleware,
		monitorController:  monitorController,
	}
}

func (r *MonitorRoute) Setup(router chi.Router) {
	r.logger.Info("setting up monitor routes")

	router.Group(func(monitorRouter chi.Router) {
		monitorRouter.Use(r.profilerMiddleware.Handler())

		monitorRouter.Get("/health", r.monitorController.HealthCheck)

		monitorRouter.Handle("/metrics", r.metrics.Handler())

		// Mount pprof routes.
		monitorRouter.Route("/debug/pprof", func(pprofRouter chi.Router) {
			pprofRouter.HandleFunc("/", pprof.Index)
			pprofRouter.HandleFunc("/cmdline", pprof.Cmdline)
			pprofRouter.HandleFunc("/profile", pprof.Profile)
			pprofRouter.HandleFunc("/symbol", pprof.Symbol)
			pprofRouter.HandleFunc("/trace", pprof.Trace)

			// Also handle sub-paths like /heap, /goroutine, etc.
			pprofRouter.Handle("/{name}", http.HandlerFunc(pprof.Index))
		})
	})
}
