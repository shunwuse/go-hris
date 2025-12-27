package repositories

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewMonitorRepository,
	NewUserRepository,
	NewRoleRepository,
	NewApprovalRepository,
	NewAuthRepository,
)
