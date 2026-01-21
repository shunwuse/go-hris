package repositories

import (
	"context"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/ent/entgen/role"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/infra"
	"github.com/shunwuse/go-hris/internal/ports/repository"
	"go.uber.org/zap"
)

type RoleRepository struct {
	logger *infra.Logger
	*infra.Database
	cache *infra.Cache
}

func NewRoleRepository(
	logger *infra.Logger,
	db *infra.Database,
	cache *infra.Cache,
) repository.RoleRepository {
	return &RoleRepository{
		logger:   logger,
		Database: db,
		cache:    cache,
	}
}

func (r *RoleRepository) FindAllRoles(ctx context.Context) ([]*entgen.Role, error) {
	return infra.CacheGetOrSet(ctx, r.cache, constants.GetAllRolesKey(), constants.CacheTTLAllRoles, func() ([]*entgen.Role, error) {
		roles, err := r.GetClient(ctx).Role.Query().
			All(ctx)
		if err != nil {
			r.logger.WithContext(ctx).Error("failed to find all roles", zap.Error(err))
			return nil, errors.ErrDatabaseError
		}
		return roles, nil
	})
}

func (r *RoleRepository) FindRoleByName(ctx context.Context, name constants.Role) *entgen.Role {
	roles, err := r.FindAllRoles(ctx)
	if err != nil {
		return nil
	}

	for _, role := range roles {
		if role.Name == name {
			return role
		}
	}

	return nil
}

func (r *RoleRepository) CreateRole(ctx context.Context, name constants.Role) (*entgen.Role, error) {
	role, err := r.GetClient(ctx).Role.Create().
		SetName(name).
		Save(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("failed to create role", zap.Error(err))
		return nil, errors.ErrDatabaseError
	}

	r.invalidateRolesCache(ctx)

	return role, nil
}

func (r *RoleRepository) GetPermissionsByRole(ctx context.Context, roleName constants.Role) constants.Permissions {
	permissions, _ := infra.CacheGetOrSet(ctx, r.cache, constants.GetRolePermissionsKey(roleName), constants.CacheTTLRolePermissions, func() (constants.Permissions, error) {
		roleData, err := r.GetClient(ctx).Role.Query().
			Where(role.NameEQ(roleName)).
			WithPermissions().
			Only(ctx)
		if err != nil {
			r.logger.WithContext(ctx).Error("failed to get permissions for role", zap.Error(err), zap.String("role", string(roleName)))
			return nil, err
		}

		permissions := make(constants.Permissions, len(roleData.Edges.Permissions))
		for idx, perm := range roleData.Edges.Permissions {
			permissions[idx] = perm.Description
		}

		return permissions, nil
	})

	return permissions
}

func (r *RoleRepository) invalidateRolesCache(ctx context.Context) {
	r.cache.Client.Del(ctx, constants.GetAllRolesKey())
}
