package repositories

import (
	"context"
	"log/slog"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/infra"
)

type RoleRepository struct {
	infra.Database

	Roles []*entgen.Role
}

func NewRoleRepository(
	db infra.Database,
) RoleRepository {
	// Initialize roles
	// var roles []models.Role
	// db.Find(&roles)
	roles, _ := db.Client.Role.
		Query().
		All(context.Background())

	return RoleRepository{
		Database: db,
		Roles:    roles,
	}
}

func (r RoleRepository) getAllRoles(ctx context.Context) error {
	roles, err := r.Client.Role.
		Query().
		All(ctx)
	if err != nil {
		slog.Error("Error getting roles", "error", err)
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
		slog.Error("Error adding role", "error", err)
		return err
	}

	if err := r.getAllRoles(ctx); err != nil {
		return err
	}

	return nil
}
