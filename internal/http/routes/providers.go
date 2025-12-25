package routes

import (
	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	NewHealthRoute,
	NewUserRoute,
	NewApprovalRoute,
	NewAuthRoute,

	NewRoutes,
)
