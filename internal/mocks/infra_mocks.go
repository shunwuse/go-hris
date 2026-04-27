package mocks

import (
	"context"
	"time"

	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/stretchr/testify/mock"
)

// ========================================
// Transactor Mock
// ========================================

// MockTransactor is a mock implementation of ports.Transactor
type MockTransactor struct {
	mock.Mock
}

// WithTx simulates a transaction by directly executing the work function
func (m *MockTransactor) WithTx(ctx context.Context, work func(ctx context.Context) error) error {
	args := m.Called(ctx, work)
	// If the work function should be executed, execute it
	if args.Get(0) == nil {
		return work(ctx)
	}
	return args.Error(0)
}

// MockTransactorWithError returns an error without executing work
type MockTransactorWithError struct {
	mock.Mock
	Err error
}

func (m *MockTransactorWithError) WithTx(ctx context.Context, work func(ctx context.Context) error) error {
	return m.Err
}

// ========================================
// Token Service Mock
// ========================================

// MockTokenService is a mock implementation of ports.TokenService.
type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) GenerateAccessToken(ctx context.Context, user *domains.UserWithPermissions) (string, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) ValidateAccessToken(ctx context.Context, tokenString string) (*domains.Claims, error) {
	args := m.Called(ctx, tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domains.Claims), args.Error(1)
}

func (m *MockTokenService) BlacklistToken(ctx context.Context, jti string, expiration time.Duration) error {
	args := m.Called(ctx, jti, expiration)
	return args.Error(0)
}
