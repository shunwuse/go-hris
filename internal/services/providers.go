package services

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewMonitorService,
	NewUserService,
	NewAuthService,
	NewApprovalService,
)
