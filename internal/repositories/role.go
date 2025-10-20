package repositories

import (
	"context"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/infra"
)

type RoleRepository struct {
	logger infra.Logger
	infra.Database

	Roles []*entgen.Role
}

func NewRoleRepository(
	logger infra.Logger,
	db infra.Database,
) RoleRepository {
	// Initialize roles
	// var roles []models.Role
	// db.Find(&roles)
	roles, _ := db.Client.Role.
		Query().
		All(context.Background())

	return RoleRepository{
		logger:   logger,
		Database: db,
		Roles:    roles,
	}
}

func (r RoleRepository) getAllRoles(ctx context.Context) error {
	roles, err := r.Client.Role.
		Query().
		All(ctx)
	if err != nil {
		r.logger.Errorf("Error getting roles: %v", err)
		return err
	}

	r.Roles = roles

	return nil
}

func (r RoleRepository) GetRoleByName(ctx context.Context, name string) *entgen.Role {
	for _, role := range r.Roles {
		if role.Name == name {
			return role
		}
	}

	return nil
}

func (r RoleRepository) AddRole(ctx context.Context, role *domains.RoleCreate) error {
	_, err := r.Client.Role.
		Create().
		SetName(role.Name).
		Save(ctx)
	if err != nil {
		r.logger.Errorf("Error adding role: %v", err)
		return err
	}

	if err := r.getAllRoles(ctx); err != nil {
		return err
	}

	return nil
}
