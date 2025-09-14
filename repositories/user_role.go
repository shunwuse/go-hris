package repositories

import (
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/models"
)

type UserRoleRepository struct {
	logger lib.Logger
	lib.Database

	UserRoleMap []models.UserRole
}

func NewUserRoleRepository(
	logger lib.Logger,
	db lib.Database,
) UserRoleRepository {
	// Initialize user roles
	var userRoles []models.UserRole
	db.Find(&userRoles)

	return UserRoleRepository{
		logger:      logger,
		Database:    db,
		UserRoleMap: userRoles,
	}
}
