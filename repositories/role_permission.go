package repositories

import (
	"github.com/shunwuse/go-hris/constants"
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/models"
)

type RolePermissionRepository struct {
	logger lib.Logger
	lib.Database

	rolePermissionMap map[constants.Role][]constants.Permission
}

func NewRolePermissionRepository() RolePermissionRepository {
	logger := lib.GetLogger()
	db := lib.GetDatabase()

	rolePermissionList := make([]models.RolePermission, 0)
	db.Preload("Role").Preload("Permission").Find(&rolePermissionList)

	rolePermissionMap := make(map[constants.Role][]constants.Permission)
	for _, rolePermission := range rolePermissionList {
		role := constants.Role(rolePermission.Role.Name)
		permission := constants.Permission(rolePermission.Permission.Description)

		if _, ok := rolePermissionMap[role]; !ok {
			rolePermissionMap[role] = make([]constants.Permission, 0)
		}

		rolePermissionMap[role] = append(rolePermissionMap[role], permission)
	}

	return RolePermissionRepository{
		logger:            logger,
		Database:          db,
		rolePermissionMap: rolePermissionMap,
	}
}

func (r RolePermissionRepository) GetPermissionsByRole(role constants.Role) []constants.Permission {
	return r.rolePermissionMap[role]
}
