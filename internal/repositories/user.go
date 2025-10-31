package repositories

import (
	"context"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/ent/entgen/user"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/infra"
	"go.uber.org/zap"
)

type UserRepository struct {
	logger *infra.Logger
	*infra.Database
}

func NewUserRepository(
	logger *infra.Logger,
	db *infra.Database,
) *UserRepository {
	return &UserRepository{
		logger:   logger,
		Database: db,
	}
}

func (r *UserRepository) FindAll(ctx context.Context) ([]*entgen.User, error) {
	users, err := r.Client.User.Query().
		All(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("failed to find all users", zap.Error(err))
		return nil, errors.ErrDatabaseError
	}

	return users, nil
}

func (r *UserRepository) Create(ctx context.Context, username string, name string) (*entgen.User, error) {
	u, err := r.Client.User.Create().
		SetUsername(username).
		SetName(name).
		Save(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("failed to create user", zap.Error(err))
		return nil, errors.ErrDatabaseError
	}

	return u, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*entgen.User, error) {
	u, err := r.Client.User.Query().
		WithPassword().
		WithRoles().
		Where(user.UsernameEQ(username)).
		Only(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("failed to find user by username", zap.Error(err), zap.String("username", username))
		if entgen.IsNotFound(err) {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrDatabaseError
	}

	return u, nil
}

func (r *UserRepository) Update(ctx context.Context, id uint, name string) error {
	if err := r.Client.User.Update().
		Where(user.IDEQ(id)).
		SetName(name).
		Exec(ctx); err != nil {
		r.logger.WithContext(ctx).Error("failed to update user", zap.Error(err), zap.Uint("user_id", id))
		return errors.ErrDatabaseError
	}

	return nil
}
