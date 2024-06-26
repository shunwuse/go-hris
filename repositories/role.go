package repositories

import (
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/models"
)

type RoleRepository struct {
	logger lib.Logger
	lib.Database

	Roles []models.Role
}

func NewRoleRepository() RoleRepository {
	logger := lib.GetLogger()
	db := lib.GetDatabase()

	// Initialize roles
	var roles []models.Role
	db.Find(&roles)

	return RoleRepository{
		logger:   logger,
		Database: db,
		Roles:    roles,
	}
}

func (r RoleRepository) getAllRoles() error {
	result := r.Find(&r.Roles)
	if result.Error != nil {
		r.logger.Errorf("Error getting roles: %v", result.Error)
		return result.Error
	}

	return nil
}

func (r RoleRepository) GetRoleByName(name string) *models.Role {
	for _, role := range r.Roles {
		if role.Name == name {
			return &role
		}
	}

	return nil
}

func (r RoleRepository) AddRole(role *models.Role) error {
	result := r.Create(role)
	if result.Error != nil {
		r.logger.Errorf("Error adding role: %v", result.Error)
		return result.Error
	}

	if err := r.getAllRoles(); err != nil {
		return err
	}

	return nil
}
