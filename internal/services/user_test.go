package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/infra/cache"
	"github.com/shunwuse/go-hris/internal/mocks"
	"github.com/shunwuse/go-hris/internal/pkg/contextx"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
	"github.com/shunwuse/go-hris/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ========================================
// Test Setup
// ========================================

// testUserServiceDependencies holds the common dependencies for UserService tests
type testUserServiceDependencies struct {
	logger         *logger.Logger
	miniRedis      *miniredis.Miniredis
	cache          *cache.Cache
	transactor     *mocks.MockTransactor
	userRepository *mocks.MockUserRepository
	roleRepository *mocks.MockRoleRepository
	userReader     *mocks.MockUserIdentityReader
}

func setupTestUserServiceDependencies(t *testing.T) *testUserServiceDependencies {
	miniRedis := miniredis.RunT(t)

	client := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})

	c := &cache.Cache{Client: client}

	return &testUserServiceDependencies{
		logger:         logger.NewNopLogger(),
		miniRedis:      miniRedis,
		cache:          c,
		transactor:     new(mocks.MockTransactor),
		userRepository: new(mocks.MockUserRepository),
		roleRepository: new(mocks.MockRoleRepository),
		userReader:     new(mocks.MockUserIdentityReader),
	}
}

// ========================================
// GetUsers Tests
// ========================================

func TestUserService_GetUsers(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		mockSetup     func(*mocks.MockUserRepository)
		expectedUsers []domains.User
		expectedErr   error
	}{
		{
			name: "Get all users success",
			ctx:  contextWithAdmin(),
			mockSetup: func(m *mocks.MockUserRepository) {
				users := []domains.User{
					createTestUserEntity(1, "admin", "Admin User"),
					createTestUserEntity(2, "staff", "Staff User"),
				}
				m.On("FindAll", mock.Anything).Return(users, nil)
			},
			expectedUsers: []domains.User{
				createTestUserEntity(1, "admin", "Admin User"),
				createTestUserEntity(2, "staff", "Staff User"),
			},
			expectedErr: nil,
		},
		{
			name:        "No identity provided",
			ctx:         context.Background(), // no identity
			mockSetup:   func(m *mocks.MockUserRepository) {},
			expectedErr: errors.ErrUnauthorized,
		},
		{
			name:        "No read permission",
			ctx:         contextWithStaff(constants.Permissions{}), // no permissions
			mockSetup:   func(m *mocks.MockUserRepository) {},
			expectedErr: errors.ErrForbidden,
		},
		{
			name: "Repository error",
			ctx:  contextWithAdmin(),
			mockSetup: func(m *mocks.MockUserRepository) {
				m.On("FindAll", mock.Anything).Return(nil, errors.ErrInternalError)
			},
			expectedErr: errors.ErrInternalError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := setupTestUserServiceDependencies(t)
			tt.mockSetup(deps.userRepository)

			svc := services.NewUserService(
				deps.logger,
				deps.transactor,
				deps.userRepository,
				deps.roleRepository,
				deps.userReader,
			)

			users, err := svc.GetUsers(tt.ctx)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
				assert.Nil(t, users)
			} else {
				assert.NoError(t, err)
				assert.Len(t, users, len(tt.expectedUsers))
			}

			deps.userRepository.AssertExpectations(t)
		})
	}
}

// ========================================
// GetUsersWithOffset Tests
// ========================================

func TestUserService_GetUsersWithOffset(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		query       domains.OffsetQuery
		filter      domains.UserFilter
		mockSetup   func(*mocks.MockUserRepository)
		checkResult func(*testing.T, *domains.OffsetResult[domains.User], error)
	}{
		{
			name: "Get users with pagination success",
			ctx:  contextWithAdmin(),
			query: domains.OffsetQuery{
				Page:    1,
				PerPage: 10,
			},
			filter: domains.UserFilter{},
			mockSetup: func(m *mocks.MockUserRepository) {
				result := &domains.OffsetResult[domains.User]{
					Items: []domains.User{
						createTestUserEntity(1, "admin", "Admin User"),
					},
					TotalCount: 1,
				}
				m.On("FindAllWithOffset", mock.Anything, mock.Anything, mock.Anything).Return(result, nil)
			},
			checkResult: func(t *testing.T, result *domains.OffsetResult[domains.User], err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Items, 1)
				assert.Equal(t, 1, result.TotalCount)
			},
		},
		{
			name:   "Get users with filter success",
			ctx:    contextWithAdmin(),
			query:  domains.OffsetQuery{Page: 1, PerPage: 10},
			filter: domains.UserFilter{Name: "test", Role: constants.Staff},
			mockSetup: func(m *mocks.MockUserRepository) {
				result := &domains.OffsetResult[domains.User]{
					Items:      []domains.User{},
					TotalCount: 0,
				}
				m.On("FindAllWithOffset", mock.Anything, mock.Anything, mock.Anything).Return(result, nil)
			},
			checkResult: func(t *testing.T, result *domains.OffsetResult[domains.User], err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Empty(t, result.Items)
			},
		},
		{
			name:      "No identity provided",
			ctx:       context.Background(),
			mockSetup: func(m *mocks.MockUserRepository) {},
			checkResult: func(t *testing.T, result *domains.OffsetResult[domains.User], err error) {
				assert.ErrorIs(t, err, errors.ErrUnauthorized)
				assert.Nil(t, result)
			},
		},
		{
			name:      "No read permission",
			ctx:       contextWithStaff(constants.Permissions{}),
			mockSetup: func(m *mocks.MockUserRepository) {},
			checkResult: func(t *testing.T, result *domains.OffsetResult[domains.User], err error) {
				assert.ErrorIs(t, err, errors.ErrForbidden)
				assert.Nil(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := setupTestUserServiceDependencies(t)
			tt.mockSetup(deps.userRepository)

			svc := services.NewUserService(
				deps.logger,
				deps.transactor,
				deps.userRepository,
				deps.roleRepository,
				deps.userReader,
			)

			result, err := svc.GetUsersWithOffset(tt.ctx, tt.query, tt.filter)

			tt.checkResult(t, result, err)

			deps.userRepository.AssertExpectations(t)
		})
	}
}

// ========================================
// GetUserByID Tests
// ========================================

func TestUserService_GetUserByID(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		userID      uint
		mockSetup   func(*mocks.MockUserIdentityReader)
		checkResult func(*testing.T, *domains.UserWithPermissions, error)
	}{
		{
			name:   "Get user success",
			ctx:    contextWithAdmin(),
			userID: 1,
			mockSetup: func(reader *mocks.MockUserIdentityReader) {
				user := createTestUserEntity(1, "admin", "Admin User")
				result := &domains.UserWithPermissions{
					ID:          user.ID,
					Username:    user.Username,
					Name:        user.Name,
					CreatedAt:   user.CreatedAt,
					UpdatedAt:   user.UpdatedAt,
					Permissions: constants.Permissions{"user:read"},
				}
				reader.On("GetUserWithPermissionsByID", mock.Anything, uint(1)).Return(result, nil)
			},
			checkResult: func(t *testing.T, result *domains.UserWithPermissions, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, uint(1), result.ID)
				assert.Equal(t, "admin", result.Username)
				assert.Contains(t, result.Permissions, constants.Permission("user:read"))
			},
		},
		{
			name:   "User not found",
			ctx:    contextWithAdmin(),
			userID: 999,
			mockSetup: func(reader *mocks.MockUserIdentityReader) {
				reader.On("GetUserWithPermissionsByID", mock.Anything, uint(999)).Return(nil, errors.ErrNotFound)
			},
			checkResult: func(t *testing.T, result *domains.UserWithPermissions, err error) {
				assert.ErrorIs(t, err, errors.ErrNotFound)
				assert.Nil(t, result)
			},
		},
		{
			name:      "No identity provided",
			ctx:       context.Background(),
			userID:    1,
			mockSetup: func(reader *mocks.MockUserIdentityReader) {},
			checkResult: func(t *testing.T, result *domains.UserWithPermissions, err error) {
				assert.ErrorIs(t, err, errors.ErrUnauthorized)
				assert.Nil(t, result)
			},
		},
		{
			name:      "No read permission",
			ctx:       contextWithStaff(constants.Permissions{}),
			userID:    1,
			mockSetup: func(reader *mocks.MockUserIdentityReader) {},
			checkResult: func(t *testing.T, result *domains.UserWithPermissions, err error) {
				assert.ErrorIs(t, err, errors.ErrForbidden)
				assert.Nil(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := setupTestUserServiceDependencies(t)
			tt.mockSetup(deps.userReader)

			deps.miniRedis.FlushAll() // Clear cache for each test

			svc := services.NewUserService(
				deps.logger,
				deps.transactor,
				deps.userRepository,
				deps.roleRepository,
				deps.userReader,
			)

			result, err := svc.GetUserByID(tt.ctx, tt.userID)

			tt.checkResult(t, result, err)

			deps.userReader.AssertExpectations(t)
		})
	}
}

// ========================================
// CreateUser Tests
// ========================================

func TestUserService_CreateUser(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		user        *domains.UserCreate
		role        constants.Role
		mockSetup   func(*mocks.MockUserRepository, *mocks.MockRoleRepository, *mocks.MockTransactor)
		expectedErr error
	}{
		{
			name: "Create user success",
			ctx:  contextWithAdmin(),
			user: &domains.UserCreate{
				Username: "NewUser",
				Name:     "New User",
				Password: "password123",
			},
			role: constants.Staff,
			mockSetup: func(userRepo *mocks.MockUserRepository, roleRepo *mocks.MockRoleRepository, tx *mocks.MockTransactor) {
				user := createTestUserEntity(3, "newuser", "New User")
				role := &domains.Role{ID: 2, Name: constants.Staff}

				userRepo.On("Create", mock.Anything, "newuser", "New User").Return(&user, nil)
				userRepo.On("CreatePassword", mock.Anything, mock.AnythingOfType("string"), user.ID).Return(nil)
				roleRepo.On("FindRoleByName", mock.Anything, constants.Staff).Return(role)
				userRepo.On("AssignRole", mock.Anything, uint(3), uint(2)).Return(nil)
				tx.On("WithTx", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "Cannot create admin role user",
			ctx:  contextWithAdmin(),
			user: &domains.UserCreate{
				Username: "admin2",
				Name:     "Admin 2",
				Password: "password123",
			},
			role: constants.Admin,
			mockSetup: func(userRepo *mocks.MockUserRepository, roleRepo *mocks.MockRoleRepository, tx *mocks.MockTransactor) {
			},
			expectedErr: errors.ErrOperationNotAllowed,
		},
		{
			name: "No identity provided",
			ctx:  context.Background(),
			user: &domains.UserCreate{
				Username: "test",
				Name:     "Test",
				Password: "password123",
			},
			role: constants.Staff,
			mockSetup: func(userRepo *mocks.MockUserRepository, roleRepo *mocks.MockRoleRepository, tx *mocks.MockTransactor) {
			},
			expectedErr: errors.ErrUnauthorized,
		},
		{
			name: "No create permission",
			ctx:  contextWithStaff(constants.Permissions{constants.PermissionReadUser}), // only read
			user: &domains.UserCreate{
				Username: "test",
				Name:     "Test",
				Password: "password123",
			},
			role: constants.Staff,
			mockSetup: func(userRepo *mocks.MockUserRepository, roleRepo *mocks.MockRoleRepository, tx *mocks.MockTransactor) {
			},
			expectedErr: errors.ErrForbidden,
		},
		{
			name: "Username already exists",
			ctx:  contextWithAdmin(),
			user: &domains.UserCreate{
				Username: "existing",
				Name:     "Existing User",
				Password: "password123",
			},
			role: constants.Staff,
			mockSetup: func(userRepo *mocks.MockUserRepository, roleRepo *mocks.MockRoleRepository, tx *mocks.MockTransactor) {
				userRepo.On("Create", mock.Anything, "existing", "Existing User").Return(nil, errors.ErrAlreadyExists)
				tx.On("WithTx", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
			},
			expectedErr: errors.ErrAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := setupTestUserServiceDependencies(t)
			tt.mockSetup(deps.userRepository, deps.roleRepository, deps.transactor)

			svc := services.NewUserService(
				deps.logger,
				deps.transactor,
				deps.userRepository,
				deps.roleRepository,
				deps.userReader,
			)

			err := svc.CreateUser(tt.ctx, tt.user, tt.role)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			deps.userRepository.AssertExpectations(t)
			deps.roleRepository.AssertExpectations(t)
		})
	}
}

// ========================================
// UpdateUser Tests
// ========================================

func TestUserService_UpdateUser(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		update      *domains.UserUpdate
		mockSetup   func(*mocks.MockUserRepository)
		expectedErr error
	}{
		{
			name: "Update user success",
			ctx:  contextWithAdmin(),
			update: &domains.UserUpdate{
				ID:   1,
				Name: "Updated Name",
			},
			mockSetup: func(m *mocks.MockUserRepository) {
				m.On("Update", mock.Anything, uint(1), "Updated Name").Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "User not found",
			ctx:  contextWithAdmin(),
			update: &domains.UserUpdate{
				ID:   999,
				Name: "Updated Name",
			},
			mockSetup: func(m *mocks.MockUserRepository) {
				m.On("Update", mock.Anything, uint(999), "Updated Name").Return(errors.ErrNotFound)
			},
			expectedErr: errors.ErrNotFound,
		},
		{
			name: "No identity provided",
			ctx:  context.Background(),
			update: &domains.UserUpdate{
				ID:   1,
				Name: "Updated Name",
			},
			mockSetup:   func(m *mocks.MockUserRepository) {},
			expectedErr: errors.ErrUnauthorized,
		},
		{
			name: "No update permission",
			ctx:  contextWithStaff(constants.Permissions{constants.PermissionReadUser}),
			update: &domains.UserUpdate{
				ID:   1,
				Name: "Updated Name",
			},
			mockSetup:   func(m *mocks.MockUserRepository) {},
			expectedErr: errors.ErrForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := setupTestUserServiceDependencies(t)
			tt.mockSetup(deps.userRepository)

			svc := services.NewUserService(
				deps.logger,
				deps.transactor,
				deps.userRepository,
				deps.roleRepository,
				deps.userReader,
			)

			err := svc.UpdateUser(tt.ctx, tt.update)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			deps.userRepository.AssertExpectations(t)
		})
	}
}

// ========================================
// Helper Functions
// ========================================

// Helper to create context with admin identity (has all permissions)
func contextWithAdmin() context.Context {
	return contextx.WithIdentity(context.Background(), &domains.Identity{
		UserID:      1,
		Username:    "admin",
		Roles:       []constants.Role{constants.Admin},
		Permissions: constants.AllPermissions,
	})
}

// Helper to create context with staff identity (limited permissions)
func contextWithStaff(permissions constants.Permissions) context.Context {
	return contextx.WithIdentity(context.Background(), &domains.Identity{
		UserID:      2,
		Username:    "staff",
		Roles:       []constants.Role{constants.Staff},
		Permissions: permissions,
	})
}

// Helper to create test user entity
func createTestUserEntity(id uint, username string, name string) domains.User {
	now := time.Now()
	return domains.User{
		ID:        id,
		Username:  username,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
