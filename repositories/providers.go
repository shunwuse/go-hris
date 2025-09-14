package repositories

import (
	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	NewExampleRepository,
	NewUserRepository,
	NewRoleRepository,
	NewUserRoleRepository,
	NewApprovalRepository,
	NewRolePermissionRepository,
)
