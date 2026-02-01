package mocks

import (
	"context"
	"time"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/stretchr/testify/mock"
)

// ========================================
// Monitor Repository Mock
// ========================================

// MockMonitorRepository is a mock implementation of repository.MonitorRepository
type MockMonitorRepository struct {
	mock.Mock
}

func (m *MockMonitorRepository) CheckDatabase(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

func (m *MockMonitorRepository) CheckRedis(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

// ========================================
// User Repository Mock
// ========================================

// MockUserRepository is a mock implementation of repository.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindAll(ctx context.Context) ([]*entgen.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entgen.User), args.Error(1)
}

func (m *MockUserRepository) FindAllWithOffset(ctx context.Context, query domains.OffsetQuery, filter domains.UserFilter) (*domains.OffsetResult[*entgen.User], error) {
	args := m.Called(ctx, query, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domains.OffsetResult[*entgen.User]), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*entgen.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entgen.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint) (*entgen.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entgen.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, username string, name string) (*entgen.User, error) {
	args := m.Called(ctx, username, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entgen.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, id uint, name string) error {
	args := m.Called(ctx, id, name)
	return args.Error(0)
}

func (m *MockUserRepository) CreatePassword(ctx context.Context, hash string, owner *entgen.User) (*entgen.Password, error) {
	args := m.Called(ctx, hash, owner)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entgen.Password), args.Error(1)
}

func (m *MockUserRepository) AssignRole(ctx context.Context, userID uint, roleID uint) (*entgen.UserRole, error) {
	args := m.Called(ctx, userID, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entgen.UserRole), args.Error(1)
}

// ========================================
// Auth Repository Mock
// ========================================

// MockAuthRepository is a mock implementation of repository.AuthRepository
type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) CreateRefreshToken(ctx context.Context, tokenHash string, userID uint, expiresAt time.Time) (*entgen.RefreshToken, error) {
	args := m.Called(ctx, tokenHash, userID, expiresAt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entgen.RefreshToken), args.Error(1)
}

func (m *MockAuthRepository) FindRefreshTokenByTokenHash(ctx context.Context, tokenHash string) (*entgen.RefreshToken, error) {
	args := m.Called(ctx, tokenHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entgen.RefreshToken), args.Error(1)
}

func (m *MockAuthRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	args := m.Called(ctx, tokenHash)
	return args.Error(0)
}

func (m *MockAuthRepository) RevokeAllRefreshTokensForUser(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockAuthRepository) DeleteExpiredTokens(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

// ========================================
// Role Repository Mock
// ========================================

// MockRoleRepository is a mock implementation of repository.RoleRepository
type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) FindAllRoles(ctx context.Context) ([]*entgen.Role, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entgen.Role), args.Error(1)
}

func (m *MockRoleRepository) FindRoleByName(ctx context.Context, name constants.Role) *entgen.Role {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*entgen.Role)
}

func (m *MockRoleRepository) CreateRole(ctx context.Context, name constants.Role) (*entgen.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entgen.Role), args.Error(1)
}

func (m *MockRoleRepository) GetPermissionsByRole(ctx context.Context, roleName constants.Role) constants.Permissions {
	args := m.Called(ctx, roleName)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(constants.Permissions)
}

// MockApprovalRepository is a mock implementation of repository.ApprovalRepository
type MockApprovalRepository struct {
	mock.Mock
}

func (m *MockApprovalRepository) FindAllWithRelations(ctx context.Context) ([]*entgen.Approval, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entgen.Approval), args.Error(1)
}

func (m *MockApprovalRepository) FindAllWithCursor(ctx context.Context, query domains.CursorQuery, filter domains.ApprovalFilter) (*domains.CursorResult[*entgen.Approval], error) {
	args := m.Called(ctx, query, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domains.CursorResult[*entgen.Approval]), args.Error(1)
}

func (m *MockApprovalRepository) FindByID(ctx context.Context, id uint) (*entgen.Approval, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entgen.Approval), args.Error(1)
}

func (m *MockApprovalRepository) Create(ctx context.Context, approval *domains.ApprovalCreate) (*entgen.Approval, error) {
	args := m.Called(ctx, approval)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entgen.Approval), args.Error(1)
}

func (m *MockApprovalRepository) UpdateStatusByID(ctx context.Context, id uint, status constants.ApprovalStatus, approverID uint) error {
	args := m.Called(ctx, id, status, approverID)
	return args.Error(0)
}
