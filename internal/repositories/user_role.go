package repositories

import (
	"context"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/internal/infra"
)

type UserRoleRepository struct {
	logger infra.Logger
	infra.Database

	UserRoleMap []*entgen.UserRole
}

func NewUserRoleRepository(
	logger infra.Logger,
	db infra.Database,
) UserRoleRepository {
	// Initialize user roles.
	userRoles, _ := db.Client.UserRole.
		Query().
		All(context.Background())

	return UserRoleRepository{
		logger:      logger,
		Database:    db,
		UserRoleMap: userRoles,
	}
}
