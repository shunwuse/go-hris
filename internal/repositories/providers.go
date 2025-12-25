package repositories

import (
	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	NewHealthRepository,
	NewUserRepository,
	NewRoleRepository,
	NewApprovalRepository,
	NewAuthRepository,
)
