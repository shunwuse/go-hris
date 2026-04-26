package repository

import (
	"context"

	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
)

type RoleRepository interface {
	FindAllRoles(ctx context.Context) ([]domains.Role, error)
	FindRoleByName(ctx context.Context, name constants.Role) *domains.Role
	CreateRole(ctx context.Context, name constants.Role) (*domains.Role, error)
	GetPermissionsByRole(ctx context.Context, roleName constants.Role) constants.Permissions
}
