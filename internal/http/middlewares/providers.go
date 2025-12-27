package middlewares

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewTraceMiddleware,
	NewMetricsMiddleware,
	NewRequestLoggerMiddleware,
	NewRecoveryMiddleware,
	NewJWTMiddleware,

	NewCommonMiddlewares,
)
