package repositories

import (
	"context"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/infra"
	"go.uber.org/zap"
)

type UserRoleRepository struct {
	logger *infra.Logger
	*infra.Database

	UserRoleMap []*entgen.UserRole
}

func NewUserRoleRepository(
	logger *infra.Logger,
	db *infra.Database,
) *UserRoleRepository {
	// Initialize user roles.
	userRoles, _ := db.Client.UserRole.
		Query().
		All(context.Background())

	return &UserRoleRepository{
		logger:      logger,
		Database:    db,
		UserRoleMap: userRoles,
	}
}

func (r *UserRoleRepository) Create(ctx context.Context, userID uint, roleID uint) (*entgen.UserRole, error) {
	userRole, err := r.Client.UserRole.Create().
		SetUserID(userID).
		SetRoleID(roleID).
		Save(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("failed to create user role association", zap.Error(err))
		return nil, errors.ErrDatabaseError
	}

	return userRole, nil
}
