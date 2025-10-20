package routes

import (
	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	NewExampleRoute,
	NewUserRoute,
	NewApprovalRoute,
	NewSwaggerRoute,

	NewRoutes,
)
