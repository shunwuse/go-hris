package middlewares

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewTraceMiddleware,
	NewRecoveryMiddleware,
	NewJWTMiddleware,

	NewCommonMiddlewares,
)
