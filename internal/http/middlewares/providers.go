package middlewares

import (
	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	NewTraceMiddleware,
	NewJWTMiddleware,
	NewDBTrxMiddleware,

	NewCommonMiddlewares,
)
