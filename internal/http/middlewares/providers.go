package middlewares

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewTraceMiddleware,
	NewIdempotencyMiddleware,
	NewMetricsMiddleware,
	NewRequestLoggerMiddleware,
	NewRecoveryMiddleware,
	NewExceptionMiddleware,
	NewJWTMiddleware,
	NewProfilerMiddleware,

	NewCommonMiddlewares,
)
