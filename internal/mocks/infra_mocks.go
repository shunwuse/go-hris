package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// ========================================
// Transactor Mock
// ========================================

// MockTransactor is a mock implementation of database.Transactor
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
