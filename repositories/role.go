package repositories

import (
	"context"

	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/models"
)

type RoleRepository struct {
	logger lib.Logger
	lib.Database

	Roles []models.Role
}

func NewRoleRepository(
	logger lib.Logger,
	db lib.Database,
) RoleRepository {
	// Initialize roles
	var roles []models.Role
	db.Find(&roles)

	return RoleRepository{
		logger:   logger,
		Database: db,
		Roles:    roles,
	}
}

func (r RoleRepository) getAllRoles(ctx context.Context) error {
	result := r.Find(&r.Roles)
	if result.Error != nil {
		r.logger.Errorf("Error getting roles: %v", result.Error)
		return result.Error
	}

	return nil
}

func (r RoleRepository) GetRoleByName(ctx context.Context, name string) *models.Role {
	for _, role := range r.Roles {
		if role.Name == name {
			return &role
		}
	}

	return nil
}

func (r RoleRepository) AddRole(ctx context.Context, role *models.Role) error {
	result := r.Create(role)
	if result.Error != nil {
		r.logger.Errorf("Error adding role: %v", result.Error)
		return result.Error
	}

	if err := r.getAllRoles(ctx); err != nil {
		return err
	}

	return nil
}
