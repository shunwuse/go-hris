package middlewares

import (
	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	NewJWTMiddleware,
	NewDBTrxMiddleware,

	NewCommonMiddlewares,
)
