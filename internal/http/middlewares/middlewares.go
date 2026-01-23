package middlewares

import (
	"github.com/go-chi/chi/v5"
)

type CommonMiddlewares []ICommonMiddleware

type ICommonMiddleware interface {
	Setup(router chi.Router)
}

func NewCommonMiddlewares(
	traceMiddleware *TraceMiddleware,
	idempotencyMiddleware *IdempotencyMiddleware,
	metricsMiddleware *MetricsMiddleware,
	requestLoggerMiddleware *RequestLoggerMiddleware,
	recoveryMiddleware *RecoveryMiddleware,
) CommonMiddlewares {
	return CommonMiddlewares{
		traceMiddleware,
		idempotencyMiddleware,
		metricsMiddleware,
		requestLoggerMiddleware,
		recoveryMiddleware,
	}
}

func (m CommonMiddlewares) Setup(router chi.Router) {
	NewCORSMiddleware().Setup(router) // setup CORS middleware

	// Setup built-in middlewares.
	// router.Use(middleware.Logger) // Replaced by requestLoggerMiddleware

	for _, middleware := range m {
		middleware.Setup(router)
	}
}
