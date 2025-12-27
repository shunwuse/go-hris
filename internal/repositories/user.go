package repositories

import (
	"context"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/ent/entgen/user"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/infra"
	"go.uber.org/zap"
)

type UserRepository struct {
	logger *infra.Logger
	*infra.Database
	cache *infra.Cache
}

func NewUserRepository(
	logger *infra.Logger,
	db *infra.Database,
	cache *infra.Cache,
) *UserRepository {
	return &UserRepository{
		logger:   logger,
		Database: db,
		cache:    cache,
	}
}

func (r *UserRepository) FindAll(ctx context.Context) ([]*entgen.User, error) {
	users, err := r.GetClient(ctx).User.Query().
		All(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("failed to find all users", zap.Error(err))
		return nil, errors.ErrDatabaseError
	}

	return users, nil
}

func (r *UserRepository) Create(ctx context.Context, username string, name string) (*entgen.User, error) {
	u, err := r.GetClient(ctx).User.Create().
		SetUsername(username).
		SetName(name).
		Save(ctx)
	if err != nil {
		if entgen.IsConstraintError(err) {
			return nil, errors.ErrConflict
		}
		r.logger.WithContext(ctx).Error("failed to create user", zap.Error(err))
		return nil, errors.ErrDatabaseError
	}

	return u, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*entgen.User, error) {
	u, err := r.GetClient(ctx).User.Query().
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

func (r *UserRepository) FindByID(ctx context.Context, id uint) (*domains.UserWithPermissions, error) {
	return infra.CacheGetOrSet(ctx, r.cache, constants.GetUserPermissionsKey(id), constants.CacheTTLUserPermissions, func() (*domains.UserWithPermissions, error) {
		u, err := r.GetClient(ctx).User.Query().
			WithPassword().
			WithRoles(func(rq *entgen.RoleQuery) {
				rq.WithRolePermission(func(rpq *entgen.RolePermissionQuery) {
					rpq.WithPermission()
				})
			}).
			Where(user.IDEQ(id)).
			Only(ctx)
		if err != nil {
			r.logger.WithContext(ctx).Error("failed to find user by id", zap.Error(err), zap.Uint("user_id", id))
			if entgen.IsNotFound(err) {
				return nil, errors.ErrNotFound
			}
			return nil, errors.ErrDatabaseError
		}

		// Aggregate permissions from roles.
		permissions := make(constants.Permissions, 0)
		for _, role := range u.Edges.Roles {
			for _, rolePermission := range role.Edges.RolePermission {
				if rolePermission.Edges.Permission != nil {
					permissions = append(permissions, constants.Permission(rolePermission.Edges.Permission.Description))
				}
			}
		}

		return &domains.UserWithPermissions{
			User:        u,
			Permissions: permissions,
		}, nil
	})
}

func (r *UserRepository) Update(ctx context.Context, id uint, name string) error {
	if err := r.GetClient(ctx).User.Update().
		Where(user.IDEQ(id)).
		SetName(name).
		Exec(ctx); err != nil {
		r.logger.WithContext(ctx).Error("failed to update user", zap.Error(err), zap.Uint("user_id", id))
		return errors.ErrDatabaseError
	}

	r.invalidateUserPermissionsCache(ctx, id)

	return nil
}

func (r *UserRepository) CreatePassword(ctx context.Context, hash string, owner *entgen.User) (*entgen.Password, error) {
	password, err := r.GetClient(ctx).Password.Create().
		SetHash(hash).
		SetOwner(owner).
		Save(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("failed to create password", zap.Error(err))
		return nil, errors.ErrDatabaseError
	}

	return password, nil
}

func (r *UserRepository) AssignRole(ctx context.Context, userID uint, roleID uint) (*entgen.UserRole, error) {
	userRole, err := r.GetClient(ctx).UserRole.Create().
		SetUserID(userID).
		SetRoleID(roleID).
		Save(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("failed to create user role association", zap.Error(err))
		return nil, errors.ErrDatabaseError
	}

	r.invalidateUserPermissionsCache(ctx, userID)

	return userRole, nil
}

func (r *UserRepository) invalidateUserPermissionsCache(ctx context.Context, userID uint) {
	r.cache.Client.Del(ctx, constants.GetUserPermissionsKey(userID))
}
