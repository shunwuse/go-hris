package repositories

import (
	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	NewExampleRepository,
	NewUserRepository,
	NewPasswordRepository,
	NewRoleRepository,
	NewUserRoleRepository,
	NewApprovalRepository,
	NewRolePermissionRepository,
	NewRefreshTokenRepository,
)
