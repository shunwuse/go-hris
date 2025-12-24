package repositories

import (
	"context"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/infra"
	"go.uber.org/zap"
)

type RoleRepository struct {
	logger *infra.Logger
	*infra.Database

	Roles             []*entgen.Role
	rolePermissionMap map[constants.Role]constants.Permissions
}

func NewRoleRepository(
	logger *infra.Logger,
	db *infra.Database,
) *RoleRepository {
	roles, _ := db.Client.Role.
		Query().
		WithPermissions().
		All(context.Background())

	rolePermissionMap := make(map[constants.Role]constants.Permissions)
	for _, role := range roles {
		roleKey := constants.Role(role.Name)
		permissions := make(constants.Permissions, 0, len(role.Edges.Permissions))

		for _, p := range role.Edges.Permissions {
			permissions = append(permissions, constants.Permission(p.Description))
		}

		rolePermissionMap[roleKey] = permissions
	}

	return &RoleRepository{
		logger:            logger,
		Database:          db,
		Roles:             roles,
		rolePermissionMap: rolePermissionMap,
	}
}

func (r *RoleRepository) FindAll(ctx context.Context) error {
	roles, err := r.Client.Role.Query().
		All(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("failed to find all roles", zap.Error(err))
		return errors.ErrDatabaseError
	}

	r.Roles = roles // update cached roles

	return nil
}

func (r *RoleRepository) FindByName(ctx context.Context, name string) *entgen.Role {
	for _, role := range r.Roles {
		if role.Name == name {
			return role
		}
	}

	return nil
}

func (r *RoleRepository) Create(ctx context.Context, name string) (*entgen.Role, error) {
	role, err := r.Client.Role.Create().
		SetName(name).
		Save(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("failed to create role", zap.Error(err))
		return nil, errors.ErrDatabaseError
	}

	if err := r.FindAll(ctx); err != nil {
		return nil, err
	}

	return role, nil
}

func (r *RoleRepository) GetPermissionsByRole(ctx context.Context, role constants.Role) constants.Permissions {
	return r.rolePermissionMap[role]
}
