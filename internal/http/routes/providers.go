package routes

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewMonitorRoute,
	NewUserRoute,
	NewApprovalRoute,
	NewAuthRoute,

	NewRoutes,
)
