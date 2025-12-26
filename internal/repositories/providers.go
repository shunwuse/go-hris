package repositories

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewHealthRepository,
	NewUserRepository,
	NewRoleRepository,
	NewApprovalRepository,
	NewAuthRepository,
)
