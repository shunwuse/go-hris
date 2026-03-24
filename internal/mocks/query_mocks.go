package mocks

import (
	"context"

	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/stretchr/testify/mock"
)

// ========================================
// Query Reader Mock
// ========================================

// MockUserIdentityReader is a mock implementation of query.UserIdentityReader.
type MockUserIdentityReader struct {
	mock.Mock
}

func (m *MockUserIdentityReader) GetUserWithPermissionsByUsername(ctx context.Context, username string) (*domains.UserWithPermissions, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domains.UserWithPermissions), args.Error(1)
}

func (m *MockUserIdentityReader) GetUserWithPermissionsByID(ctx context.Context, userID uint) (*domains.UserWithPermissions, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domains.UserWithPermissions), args.Error(1)
}
