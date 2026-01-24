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
	metricsMiddleware *MetricsMiddleware,
	requestLoggerMiddleware *RequestLoggerMiddleware,
	recoveryMiddleware *RecoveryMiddleware,
	idempotencyMiddleware *IdempotencyMiddleware,
	exceptionMiddleware *ExceptionMiddleware,
) CommonMiddlewares {
	return CommonMiddlewares{
		traceMiddleware,         // Trace: Request identification
		metricsMiddleware,       // Metrics: Performance telemetry
		requestLoggerMiddleware, // Logger: Traffic observation
		recoveryMiddleware,      // Recovery: Panic protection & safety
		idempotencyMiddleware,   // Idempotency: Duplicate prevention
		exceptionMiddleware,     // Exception: Business error alerting
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
