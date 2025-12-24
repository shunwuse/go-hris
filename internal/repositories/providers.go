package repositories

import (
	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	NewExampleRepository,
	NewUserRepository,
	NewRoleRepository,
	NewApprovalRepository,
	NewAuthRepository,
)
