package repositories

import (
	"context"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/lib"
)

type UserRoleRepository struct {
	logger lib.Logger
	lib.Database

	UserRoleMap []*entgen.UserRole
}

func NewUserRoleRepository(
	logger lib.Logger,
	db lib.Database,
) UserRoleRepository {
	// Initialize user roles
	userRoles, _ := db.Client.UserRole.
		Query().
		All(context.Background())

	return UserRoleRepository{
		logger:      logger,
		Database:    db,
		UserRoleMap: userRoles,
	}
}
