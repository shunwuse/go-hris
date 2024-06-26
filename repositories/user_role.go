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

func NewUserRoleRepository() UserRoleRepository {
	logger := lib.GetLogger()
	db := lib.GetDatabase()

	// Initialize user roles
	var userRoles []models.UserRole
	db.Find(&userRoles)

	return UserRoleRepository{
		logger:      logger,
		Database:    db,
		UserRoleMap: userRoles,
	}
}
