package constants

import (
	"fmt"
	"time"
)

const (
	// Cache TTLs.
	CacheTTLAllRoles        = 1 * time.Hour
	CacheTTLRolePermissions = 1 * time.Hour
	CacheTTLUserPermissions = 30 * time.Minute
)

// Cache Key Helpers.
func GetAllRolesKey() string {
	return "roles:all"
}

func GetRolePermissionsKey(role Role) string {
	return fmt.Sprintf("role:%s:permissions", role)
}

func GetUserPermissionsKey(userID uint) string {
	return fmt.Sprintf("user:%d:permissions", userID)
}
