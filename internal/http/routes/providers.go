package routes

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewHealthRoute,
	NewUserRoute,
	NewApprovalRoute,
	NewAuthRoute,
	NewMetricsRoute,

	NewRoutes,
)
