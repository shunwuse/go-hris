package queries

import (
	"context"
	"testing"

	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/infra/cache"
	"github.com/shunwuse/go-hris/internal/mocks"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupUserIdentityReader(t *testing.T) (*mocks.MockUserRepository, *mocks.MockRoleRepository, *userIdentityReader) {
	t.Helper()

	c := cache.New(&cache.Config{UseMiniredis: true}, logger.NewNopLogger())
	t.Cleanup(func() {
		_ = c.Close()
	})

	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)

	reader := &userIdentityReader{
		cache:          c,
		userRepository: userRepo,
		roleRepository: roleRepo,
	}

	return userRepo, roleRepo, reader
}

func TestUserIdentityReader_GetUserWithPermissionsByUsername(t *testing.T) {
	userRepo, roleRepo, reader := setupUserIdentityReader(t)
	ctx := context.Background()

	baseUser := &domains.User{
		ID:       10,
		Username: "alice",
	}

	userWithRoles := &domains.UserWithRoles{
		ID:       10,
		Username: "alice",
		Roles: []constants.Role{
			constants.Admin,
			constants.Manager,
		},
	}

	userRepo.On("FindByUsername", mock.Anything, "alice").Return(baseUser, nil).Once()
	userRepo.On("FindByID", mock.Anything, uint(10)).Return(userWithRoles, nil).Once()
	roleRepo.On("GetPermissionsByRole", mock.Anything, constants.Admin).
		Return(constants.Permissions{constants.PermissionReadUser, constants.PermissionCreateUser}).Once()
	roleRepo.On("GetPermissionsByRole", mock.Anything, constants.Manager).
		Return(constants.Permissions{constants.PermissionReadUser, constants.PermissionUpdateUser}).Once()

	result, err := reader.GetUserWithPermissionsByUsername(ctx, "alice")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, uint(10), result.ID)
	assert.Equal(t, "alice", result.Username)

	assert.True(t, result.Permissions.Contains(constants.PermissionReadUser))
	assert.True(t, result.Permissions.Contains(constants.PermissionCreateUser))
	assert.True(t, result.Permissions.Contains(constants.PermissionUpdateUser))
	assert.Len(t, result.Permissions, 3)

	userRepo.AssertExpectations(t)
	roleRepo.AssertExpectations(t)
}

func TestUserIdentityReader_GetUserWithPermissionsByID_UsesCache(t *testing.T) {
	userRepo, roleRepo, reader := setupUserIdentityReader(t)
	ctx := context.Background()

	userWithRoles := &domains.UserWithRoles{
		ID:       30,
		Username: "cache-user",
		Roles: []constants.Role{
			constants.Staff,
		},
	}

	userRepo.On("FindByID", mock.Anything, uint(30)).Return(userWithRoles, nil).Once()
	roleRepo.On("GetPermissionsByRole", mock.Anything, constants.Staff).
		Return(constants.Permissions{constants.PermissionReadApproval}).Once()

	first, err := reader.GetUserWithPermissionsByID(ctx, 30)
	require.NoError(t, err)
	require.NotNil(t, first)

	second, err := reader.GetUserWithPermissionsByID(ctx, 30)
	require.NoError(t, err)
	require.NotNil(t, second)

	assert.Equal(t, first.ID, second.ID)
	assert.Equal(t, first.Permissions, second.Permissions)

	userRepo.AssertExpectations(t)
	roleRepo.AssertExpectations(t)
}

func TestUserIdentityReader_GetUserWithPermissionsByID_ReturnsError(t *testing.T) {
	userRepo, roleRepo, reader := setupUserIdentityReader(t)
	ctx := context.Background()

	userRepo.On("FindByID", mock.Anything, uint(404)).Return(nil, errors.ErrNotFound).Once()

	result, err := reader.GetUserWithPermissionsByID(ctx, 404)
	require.Error(t, err)
	assert.ErrorIs(t, err, errors.ErrNotFound)
	assert.Nil(t, result)

	userRepo.AssertExpectations(t)
	roleRepo.AssertNotCalled(t, "GetPermissionsByRole", mock.Anything, mock.Anything)
}
