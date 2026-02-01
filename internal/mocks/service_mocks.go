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
// Monitor Service Mock
// ========================================

// MockMonitorService is a mock implementation of service.MonitorService
type MockMonitorService struct {
	mock.Mock
}

func (m *MockMonitorService) HealthCheck(ctx context.Context) *domains.Health {
	args := m.Called(ctx)
	return args.Get(0).(*domains.Health)
}

// ========================================
// User Service Mock
// ========================================

// MockUserService is a mock implementation of service.UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUsers(ctx context.Context) ([]*entgen.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entgen.User), args.Error(1)
}

func (m *MockUserService) GetUsersWithOffset(ctx context.Context, query domains.OffsetQuery, filter domains.UserFilter) (*domains.OffsetResult[*entgen.User], error) {
	args := m.Called(ctx, query, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domains.OffsetResult[*entgen.User]), args.Error(1)
}

func (m *MockUserService) GetUserByUsername(ctx context.Context, username string) (*domains.UserWithPermissions, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domains.UserWithPermissions), args.Error(1)
}

func (m *MockUserService) GetUserByID(ctx context.Context, id uint) (*domains.UserWithPermissions, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domains.UserWithPermissions), args.Error(1)
}

func (m *MockUserService) CreateUser(ctx context.Context, user *domains.UserCreate, role constants.Role) error {
	args := m.Called(ctx, user, role)
	return args.Error(0)
}

func (m *MockUserService) UpdateUser(ctx context.Context, user *domains.UserUpdate) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// ========================================
// Auth Service Mock
// ========================================

// MockAuthService is a mock implementation of service.AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(ctx context.Context, username string, password string) (*domains.LoginResult, error) {
	args := m.Called(ctx, username, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domains.LoginResult), args.Error(1)
}

func (m *MockAuthService) GenerateAccessToken(ctx context.Context, user *domains.UserWithPermissions) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) ValidateAccessToken(ctx context.Context, tokenString string) (*domains.Claims, error) {
	args := m.Called(ctx, tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domains.Claims), args.Error(1)
}

func (m *MockAuthService) GenerateRefreshToken(ctx context.Context, user *domains.UserWithPermissions) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (*domains.TokenPair, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domains.TokenPair), args.Error(1)
}

func (m *MockAuthService) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	args := m.Called(ctx, refreshToken)
	return args.Error(0)
}

func (m *MockAuthService) RevokeAllUserTokens(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockAuthService) BlacklistToken(ctx context.Context, jti string, expiration time.Duration) error {
	args := m.Called(ctx, jti, expiration)
	return args.Error(0)
}

func (m *MockAuthService) CleanupExpiredTokens(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

// ========================================
// Approval Service Mock
// ========================================

// MockApprovalService is a mock implementation of service.ApprovalService
type MockApprovalService struct {
	mock.Mock
}

func (m *MockApprovalService) GetApprovals(ctx context.Context) ([]*entgen.Approval, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entgen.Approval), args.Error(1)
}

func (m *MockApprovalService) GetApprovalsWithCursor(ctx context.Context, query domains.CursorQuery, filter domains.ApprovalFilter) (*domains.CursorResult[*entgen.Approval], error) {
	args := m.Called(ctx, query, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domains.CursorResult[*entgen.Approval]), args.Error(1)
}

func (m *MockApprovalService) GetApprovalByID(ctx context.Context, id uint) (*entgen.Approval, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entgen.Approval), args.Error(1)
}

func (m *MockApprovalService) AddApproval(ctx context.Context, approval *domains.ApprovalCreate) error {
	args := m.Called(ctx, approval)
	return args.Error(0)
}

func (m *MockApprovalService) ActionApproval(ctx context.Context, approvalID uint, action constants.ApprovalStatus) error {
	args := m.Called(ctx, approvalID, action)
	return args.Error(0)
}
