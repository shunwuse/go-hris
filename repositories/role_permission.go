package repositories

import (
	"context"

	"github.com/shunwuse/go-hris/constants"
	"github.com/shunwuse/go-hris/lib"
)

type RolePermissionRepository struct {
	logger lib.Logger
	lib.Database

	rolePermissionMap map[constants.Role]constants.Permissions
}

func NewRolePermissionRepository(
	logger lib.Logger,
	db lib.Database,
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
