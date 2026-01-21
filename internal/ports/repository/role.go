package repository

import (
	"context"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/internal/constants"
)

type RoleRepository interface {
	FindAllRoles(ctx context.Context) ([]*entgen.Role, error)
	FindRoleByName(ctx context.Context, name constants.Role) *entgen.Role
	CreateRole(ctx context.Context, name constants.Role) (*entgen.Role, error)
	GetPermissionsByRole(ctx context.Context, roleName constants.Role) constants.Permissions
}
