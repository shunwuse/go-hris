package queries

import (
	"context"

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
			ID:           user.ID,
			Username:     user.Username,
			Name:         user.Name,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
			PasswordHash: user.PasswordHash,
			Roles:        user.Roles,
			Permissions:  r.collectPermissions(ctx, user.Roles),
		}, nil
	})
}

func (r *userIdentityReader) collectPermissions(ctx context.Context, roles []constants.Role) constants.Permissions {
	permissions := make(constants.Permissions, 0)

	for _, roleName := range roles {
		rolePermissions := r.roleRepository.GetPermissionsByRole(ctx, roleName)
		for _, permission := range rolePermissions {
			if !permissions.Contains(permission) {
				permissions = append(permissions, permission)
			}
		}
	}

	return permissions
}
