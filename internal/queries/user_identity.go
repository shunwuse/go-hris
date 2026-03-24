package queries

import (
	"context"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/infra/cache"
	"github.com/shunwuse/go-hris/internal/ports/query"
	"github.com/shunwuse/go-hris/internal/ports/repository"
)

type userIdentityReader struct {
	cache          *cache.Cache
	userRepository repository.UserRepository
	roleRepository repository.RoleRepository
}

func NewUserIdentityReader(
	c *cache.Cache,
	userRepository repository.UserRepository,
	roleRepository repository.RoleRepository,
) query.UserIdentityReader {
	return &userIdentityReader{
		cache:          c,
		userRepository: userRepository,
		roleRepository: roleRepository,
	}
}

func (r *userIdentityReader) GetUserWithPermissionsByUsername(ctx context.Context, username string) (*domains.UserWithPermissions, error) {
	user, err := r.userRepository.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return r.GetUserWithPermissionsByID(ctx, user.ID)
}

func (r *userIdentityReader) GetUserWithPermissionsByID(ctx context.Context, userID uint) (*domains.UserWithPermissions, error) {
	return cache.Fetch(ctx, r.cache, constants.GetUserPermissionsKey(userID), constants.CacheTTLUser, func() (*domains.UserWithPermissions, error) {
		user, err := r.userRepository.FindByID(ctx, userID)
		if err != nil {
			return nil, err
		}

		return &domains.UserWithPermissions{
			User:        user,
			Permissions: r.collectPermissions(ctx, user),
		}, nil
	})
}

func (r *userIdentityReader) collectPermissions(ctx context.Context, user *entgen.User) constants.Permissions {
	permissions := make(constants.Permissions, 0)

	for _, role := range user.Edges.Roles {
		rolePermissions := r.roleRepository.GetPermissionsByRole(ctx, role.Name)
		for _, permission := range rolePermissions {
			if !permissions.Contains(permission) {
				permissions = append(permissions, permission)
			}
		}
	}

	return permissions
}
