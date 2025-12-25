package services

import (
	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	NewHealthService,
	NewUserService,
	NewAuthService,
	NewApprovalService,
)
