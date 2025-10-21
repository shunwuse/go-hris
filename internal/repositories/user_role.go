package repositories

import (
	"context"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/internal/infra"
)

type UserRoleRepository struct {
	infra.Database

	UserRoleMap []*entgen.UserRole
}

func NewUserRoleRepository(
	db infra.Database,
) UserRoleRepository {
	// Initialize user roles
	userRoles, _ := db.Client.UserRole.
		Query().
		All(context.Background())

	return UserRoleRepository{
		Database:    db,
		UserRoleMap: userRoles,
	}
}
