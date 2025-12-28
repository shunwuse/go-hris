package constants

import (
	"fmt"
	"time"
)

const (
	// Cache TTLs.
	CacheTTLAllRoles        = 1 * time.Hour
	CacheTTLUser            = 30 * time.Minute
	CacheTTLRolePermissions = 1 * time.Hour
)

// Cache Key Helpers.
func GetBlacklistKey(jti string) string {
	return fmt.Sprintf("blacklist:%s", jti)
}

func GetUserKey(userID uint) string {
	return fmt.Sprintf("user:%d", userID)
}

func GetUserPermissionsKey(userID uint) string {
	return fmt.Sprintf("user:%d:permissions", userID)
}

func GetAllRolesKey() string {
	return "roles:all"
}

func GetRolePermissionsKey(role Role) string {
	return fmt.Sprintf("role:%s:permissions", role)
}
