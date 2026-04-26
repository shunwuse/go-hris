package services_test

import (
	"context"
	"testing"

	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/errors"
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

type testApprovalServiceDependencies struct {
	logger             *logger.Logger
	approvalRepository *mocks.MockApprovalRepository
}

func setupTestApprovalServiceDependencies() *testApprovalServiceDependencies {
	return &testApprovalServiceDependencies{
		logger:             logger.NewNopLogger(),
		approvalRepository: new(mocks.MockApprovalRepository),
	}
}

// ========================================
// GetApprovals Tests
// ========================================

func TestApprovalService_GetApprovals(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		mockSetup   func(*mocks.MockApprovalRepository)
		checkResult func(*testing.T, []domains.Approval, error)
	}{
		{
			name: "Get all approvals success",
			ctx:  contextWithManager(),
			mockSetup: func(m *mocks.MockApprovalRepository) {
				approvals := []domains.Approval{
					createTestApproval(1, constants.ApprovalStatusPending),
					createTestApproval(2, constants.ApprovalStatusApproved),
				}
				m.On("FindAllWithRelations", mock.Anything).Return(approvals, nil)
			},
			checkResult: func(t *testing.T, result []domains.Approval, err error) {
				assert.NoError(t, err)
				assert.Len(t, result, 2)
			},
		},
		{
			name:      "No identity provided",
			ctx:       context.Background(),
			mockSetup: func(m *mocks.MockApprovalRepository) {},
			checkResult: func(t *testing.T, result []domains.Approval, err error) {
				assert.ErrorIs(t, err, errors.ErrUnauthorized)
				assert.Nil(t, result)
			},
		},
		{
			name: "No read permission",
			ctx: contextx.WithIdentity(context.Background(), &domains.Identity{
				UserID:      4,
				Username:    "guest",
				Permissions: constants.Permissions{},
			}),
			mockSetup: func(m *mocks.MockApprovalRepository) {},
			checkResult: func(t *testing.T, result []domains.Approval, err error) {
				assert.ErrorIs(t, err, errors.ErrForbidden)
				assert.Nil(t, result)
			},
		},
		{
			name: "Repository error",
			ctx:  contextWithManager(),
			mockSetup: func(m *mocks.MockApprovalRepository) {
				m.On("FindAllWithRelations", mock.Anything).Return(nil, errors.ErrInternalError)
			},
			checkResult: func(t *testing.T, result []domains.Approval, err error) {
				assert.ErrorIs(t, err, errors.ErrInternalError)
				assert.Nil(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := setupTestApprovalServiceDependencies()
			tt.mockSetup(deps.approvalRepository)

			svc := services.NewApprovalService(deps.logger, deps.approvalRepository)

			result, err := svc.GetApprovals(tt.ctx)

			tt.checkResult(t, result, err)

			deps.approvalRepository.AssertExpectations(t)
		})
	}
}

// ========================================
// GetApprovalsWithCursor Tests
// ========================================

func TestApprovalService_GetApprovalsWithCursor(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		query       domains.CursorQuery
		filter      domains.ApprovalFilter
		mockSetup   func(*mocks.MockApprovalRepository)
		checkResult func(*testing.T, *domains.CursorResult[domains.Approval], error)
	}{
		{
			name:   "Get approvals with cursor success",
			ctx:    contextWithManager(),
			query:  domains.CursorQuery{Limit: 10},
			filter: domains.ApprovalFilter{},
			mockSetup: func(m *mocks.MockApprovalRepository) {
				result := &domains.CursorResult[domains.Approval]{
					Items: []domains.Approval{
						createTestApproval(1, constants.ApprovalStatusPending),
					},
					HasMore: false,
				}
				m.On("FindAllWithCursor", mock.Anything, mock.Anything, mock.Anything).Return(result, nil)
			},
			checkResult: func(t *testing.T, result *domains.CursorResult[domains.Approval], err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Items, 1)
			},
		},
		{
			name:   "Get approvals with status filter success",
			ctx:    contextWithManager(),
			query:  domains.CursorQuery{Limit: 10},
			filter: domains.ApprovalFilter{Status: constants.ApprovalStatusPending},
			mockSetup: func(m *mocks.MockApprovalRepository) {
				result := &domains.CursorResult[domains.Approval]{
					Items:   []domains.Approval{},
					HasMore: false,
				}
				m.On("FindAllWithCursor", mock.Anything, mock.Anything, mock.Anything).Return(result, nil)
			},
			checkResult: func(t *testing.T, result *domains.CursorResult[domains.Approval], err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			},
		},
		{
			name:      "No identity provided",
			ctx:       context.Background(),
			mockSetup: func(m *mocks.MockApprovalRepository) {},
			checkResult: func(t *testing.T, result *domains.CursorResult[domains.Approval], err error) {
				assert.ErrorIs(t, err, errors.ErrUnauthorized)
				assert.Nil(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := setupTestApprovalServiceDependencies()
			tt.mockSetup(deps.approvalRepository)

			svc := services.NewApprovalService(deps.logger, deps.approvalRepository)

			result, err := svc.GetApprovalsWithCursor(tt.ctx, tt.query, tt.filter)

			tt.checkResult(t, result, err)

			deps.approvalRepository.AssertExpectations(t)
		})
	}
}

// ========================================
// GetApprovalByID Tests
// ========================================

func TestApprovalService_GetApprovalByID(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		approvalID  uint
		mockSetup   func(*mocks.MockApprovalRepository)
		checkResult func(*testing.T, *domains.Approval, error)
	}{
		{
			name:       "Get single approval success",
			ctx:        contextWithManager(),
			approvalID: 1,
			mockSetup: func(m *mocks.MockApprovalRepository) {
				approval := createTestApproval(1, constants.ApprovalStatusPending)
				m.On("FindByID", mock.Anything, uint(1)).Return(&approval, nil)
			},
			checkResult: func(t *testing.T, result *domains.Approval, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, uint(1), result.ID)
				assert.Equal(t, constants.ApprovalStatusPending, result.Status)
			},
		},
		{
			name:       "Approval not found",
			ctx:        contextWithManager(),
			approvalID: 999,
			mockSetup: func(m *mocks.MockApprovalRepository) {
				m.On("FindByID", mock.Anything, uint(999)).Return(nil, errors.ErrNotFound)
			},
			checkResult: func(t *testing.T, result *domains.Approval, err error) {
				assert.ErrorIs(t, err, errors.ErrNotFound)
				assert.Nil(t, result)
			},
		},
		{
			name:       "No identity provided",
			ctx:        context.Background(),
			approvalID: 1,
			mockSetup:  func(m *mocks.MockApprovalRepository) {},
			checkResult: func(t *testing.T, result *domains.Approval, err error) {
				assert.ErrorIs(t, err, errors.ErrUnauthorized)
				assert.Nil(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := setupTestApprovalServiceDependencies()
			tt.mockSetup(deps.approvalRepository)

			svc := services.NewApprovalService(deps.logger, deps.approvalRepository)

			result, err := svc.GetApprovalByID(tt.ctx, tt.approvalID)

			tt.checkResult(t, result, err)

			deps.approvalRepository.AssertExpectations(t)
		})
	}
}

// ========================================
// AddApproval Tests
// ========================================

func TestApprovalService_AddApproval(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		approval    *domains.ApprovalCreate
		mockSetup   func(*mocks.MockApprovalRepository)
		expectedErr error
	}{
		{
			name: "Add approval success",
			ctx:  contextWithManager(),
			approval: &domains.ApprovalCreate{
				Status: constants.ApprovalStatusPending,
			},
			mockSetup: func(m *mocks.MockApprovalRepository) {
				approval := createTestApproval(1, constants.ApprovalStatusPending)
				m.On("Create", mock.Anything, mock.MatchedBy(func(a *domains.ApprovalCreate) bool {
					return a.CreatorID == 2 // manager's user ID
				})).Return(&approval, nil)
			},
			expectedErr: nil,
		},
		{
			name: "No identity provided",
			ctx:  context.Background(),
			approval: &domains.ApprovalCreate{
				Status: constants.ApprovalStatusPending,
			},
			mockSetup:   func(m *mocks.MockApprovalRepository) {},
			expectedErr: errors.ErrUnauthorized,
		},
		{
			name: "No create permission",
			ctx: contextx.WithIdentity(context.Background(), &domains.Identity{
				UserID:      5,
				Username:    "viewer",
				Permissions: constants.Permissions{constants.PermissionReadApproval}, // only read
			}),
			approval: &domains.ApprovalCreate{
				Status: constants.ApprovalStatusPending,
			},
			mockSetup:   func(m *mocks.MockApprovalRepository) {},
			expectedErr: errors.ErrForbidden,
		},
		{
			name: "Repository error",
			ctx:  contextWithManager(),
			approval: &domains.ApprovalCreate{
				Status: constants.ApprovalStatusPending,
			},
			mockSetup: func(m *mocks.MockApprovalRepository) {
				m.On("Create", mock.Anything, mock.Anything).Return(nil, errors.ErrInternalError)
			},
			expectedErr: errors.ErrInternalError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := setupTestApprovalServiceDependencies()
			tt.mockSetup(deps.approvalRepository)

			svc := services.NewApprovalService(deps.logger, deps.approvalRepository)

			err := svc.AddApproval(tt.ctx, tt.approval)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			deps.approvalRepository.AssertExpectations(t)
		})
	}
}

// ========================================
// ActionApproval Tests
// ========================================

func TestApprovalService_ActionApproval(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		approvalID  uint
		action      constants.ApprovalStatus
		mockSetup   func(*mocks.MockApprovalRepository)
		expectedErr error
	}{
		{
			name:       "Approve approval success",
			ctx:        contextWithManager(),
			approvalID: 1,
			action:     constants.ApprovalStatusApproved,
			mockSetup: func(m *mocks.MockApprovalRepository) {
				m.On("UpdateStatusByID", mock.Anything, uint(1), constants.ApprovalStatusApproved, uint(2)).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:       "Reject approval success",
			ctx:        contextWithManager(),
			approvalID: 1,
			action:     constants.ApprovalStatusRejected,
			mockSetup: func(m *mocks.MockApprovalRepository) {
				m.On("UpdateStatusByID", mock.Anything, uint(1), constants.ApprovalStatusRejected, uint(2)).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:        "Invalid approval action (PENDING)",
			ctx:         contextWithManager(),
			approvalID:  1,
			action:      constants.ApprovalStatusPending,
			mockSetup:   func(m *mocks.MockApprovalRepository) {},
			expectedErr: errors.ErrValidationFailed,
		},
		{
			name:        "Invalid approval action (empty)",
			ctx:         contextWithManager(),
			approvalID:  1,
			action:      constants.ApprovalStatus(""),
			mockSetup:   func(m *mocks.MockApprovalRepository) {},
			expectedErr: errors.ErrValidationFailed,
		},
		{
			name:        "No identity provided",
			ctx:         context.Background(),
			approvalID:  1,
			action:      constants.ApprovalStatusApproved,
			mockSetup:   func(m *mocks.MockApprovalRepository) {},
			expectedErr: errors.ErrUnauthorized,
		},
		{
			name:        "No approval action permission (read/create only)",
			ctx:         contextWithStaffApproval(),
			approvalID:  1,
			action:      constants.ApprovalStatusApproved,
			mockSetup:   func(m *mocks.MockApprovalRepository) {},
			expectedErr: errors.ErrForbidden,
		},
		{
			name:       "Approval not found",
			ctx:        contextWithManager(),
			approvalID: 999,
			action:     constants.ApprovalStatusApproved,
			mockSetup: func(m *mocks.MockApprovalRepository) {
				m.On("UpdateStatusByID", mock.Anything, uint(999), constants.ApprovalStatusApproved, uint(2)).Return(errors.ErrNotFound)
			},
			expectedErr: errors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := setupTestApprovalServiceDependencies()
			tt.mockSetup(deps.approvalRepository)

			svc := services.NewApprovalService(deps.logger, deps.approvalRepository)

			err := svc.ActionApproval(tt.ctx, tt.approvalID, tt.action)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			deps.approvalRepository.AssertExpectations(t)
		})
	}
}

// ========================================
// Helper Functions
// ========================================

// Helper to create context with manager identity (has approval permissions)
func contextWithManager() context.Context {
	return contextx.WithIdentity(context.Background(), &domains.Identity{
		UserID:   2,
		Username: "manager",
		Roles:    []constants.Role{constants.Manager},
		Permissions: constants.Permissions{
			constants.PermissionReadApproval,
			constants.PermissionCreateApproval,
			constants.PermissionActionApproval,
		},
	})
}

// Helper to create context with staff (limited permissions - no approval action)
func contextWithStaffApproval() context.Context {
	return contextx.WithIdentity(context.Background(), &domains.Identity{
		UserID:      3,
		Username:    "staff",
		Roles:       []constants.Role{constants.Staff},
		Permissions: constants.Permissions{constants.PermissionReadApproval, constants.PermissionCreateApproval},
	})
}

// Helper to create test approval entity
func createTestApproval(id uint, status constants.ApprovalStatus) domains.Approval {
	return domains.Approval{
		ID:           id,
		Status:       status,
		CreatorID:    2,
		CreatorName:  "Creator",
		ApproverID:   nil,
		ApproverName: nil,
	}
}
