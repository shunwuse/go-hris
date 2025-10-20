package repositories

import (
	"context"

	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/infra"
)

type RolePermissionRepository struct {
	logger infra.Logger
	infra.Database

	rolePermissionMap map[constants.Role]constants.Permissions
}

func NewRolePermissionRepository(
	logger infra.Logger,
	db infra.Database,
) RolePermissionRepository {
	roles, _ := db.Client.Role.
		Query().
		WithPermissions().
		All(context.Background())

	rolePermissionMap := make(map[constants.Role]constants.Permissions)
	for _, role := range roles {
		roleKey := constants.Role(role.Name)
		permissions := make(constants.Permissions, 0, len(role.Edges.Permissions))

		for _, p := range role.Edges.Permissions {
			permissions = append(permissions, constants.Permission(p.Description))
		}

		rolePermissionMap[roleKey] = permissions
	}

	return RolePermissionRepository{
		logger:            logger,
		Database:          db,
		rolePermissionMap: rolePermissionMap,
	}
}

func (r RolePermissionRepository) GetPermissionsByRole(ctx context.Context, role constants.Role) constants.Permissions {
	return r.rolePermissionMap[role]
}
