package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/dtos"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/http/controllers"
	"github.com/shunwuse/go-hris/internal/http/response"
	"github.com/shunwuse/go-hris/internal/infra/logger"
	"github.com/shunwuse/go-hris/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ========================================
// Test Setup
// ========================================

type testApprovalControllerDependencies struct {
	logger          *logger.Logger
	approvalService *mocks.MockApprovalService
}

func setupTestApprovalControllerDependencies() *testApprovalControllerDependencies {
	return &testApprovalControllerDependencies{
		logger:          logger.NewNopLogger(),
		approvalService: new(mocks.MockApprovalService),
	}
}

// ========================================
// GetApprovals Tests
// ========================================

func TestApprovalController_GetApprovals(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		mockSetup      func(*mocks.MockApprovalService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:        "Get approvals list success - default params",
			queryParams: "",
			mockSetup: func(m *mocks.MockApprovalService) {
				approverName := "Approver"
				m.On("GetApprovalsWithCursor", mock.Anything, mock.Anything, mock.Anything).Return(&domains.CursorResult[*entgen.Approval]{
					Items: []*entgen.Approval{
						createTestApproval(1, "Creator 1", constants.ApprovalStatusPending, nil),
						createTestApproval(2, "Creator 2", constants.ApprovalStatusApproved, &approverName),
					},
					NextCursor: "",
					HasMore:    false,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.CursorListResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.False(t, resp.Meta.HasMore)

				// Check data is array
				data, ok := resp.Data.([]interface{})
				assert.True(t, ok)
				assert.Len(t, data, 2)
			},
		},
		{
			name:        "Get approvals list success - with pagination params",
			queryParams: "cursor=abc123&limit=5",
			mockSetup: func(m *mocks.MockApprovalService) {
				m.On("GetApprovalsWithCursor", mock.Anything, mock.MatchedBy(func(q domains.CursorQuery) bool {
					return q.Cursor == "abc123" && q.Limit == 5
				}), mock.Anything).Return(&domains.CursorResult[*entgen.Approval]{
					Items:      []*entgen.Approval{},
					NextCursor: "next123",
					HasMore:    true,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.CursorListResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.True(t, resp.Meta.HasMore)
				assert.Equal(t, "next123", resp.Meta.NextCursor)
			},
		},
		{
			name:        "Get approvals list success - filter by status pending",
			queryParams: "status=PENDING",
			mockSetup: func(m *mocks.MockApprovalService) {
				m.On("GetApprovalsWithCursor", mock.Anything, mock.Anything, mock.MatchedBy(func(f domains.ApprovalFilter) bool {
					return f.Status == constants.ApprovalStatusPending
				})).Return(&domains.CursorResult[*entgen.Approval]{
					Items: []*entgen.Approval{
						createTestApproval(1, "Creator", constants.ApprovalStatusPending, nil),
					},
					NextCursor: "",
					HasMore:    false,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.CursorListResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
			},
		},
		{
			name:        "Get approvals list success - filter by status approved",
			queryParams: "status=APPROVED",
			mockSetup: func(m *mocks.MockApprovalService) {
				m.On("GetApprovalsWithCursor", mock.Anything, mock.Anything, mock.MatchedBy(func(f domains.ApprovalFilter) bool {
					return f.Status == constants.ApprovalStatusApproved
				})).Return(&domains.CursorResult[*entgen.Approval]{
					Items:      []*entgen.Approval{},
					NextCursor: "",
					HasMore:    false,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.CursorListResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
			},
		},
		{
			name:        "Get approvals list success - filter by status rejected",
			queryParams: "status=REJECTED",
			mockSetup: func(m *mocks.MockApprovalService) {
				m.On("GetApprovalsWithCursor", mock.Anything, mock.Anything, mock.MatchedBy(func(f domains.ApprovalFilter) bool {
					return f.Status == constants.ApprovalStatusRejected
				})).Return(&domains.CursorResult[*entgen.Approval]{
					Items:      []*entgen.Approval{},
					NextCursor: "",
					HasMore:    false,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.CursorListResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
			},
		},
		{
			name:           "Get approvals list failure - invalid status",
			queryParams:    "status=invalid_status",
			mockSetup:      func(m *mocks.MockApprovalService) {},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeValidationFailed, resp.Error.Code)
				assert.Contains(t, resp.Error.Details, "status")
			},
		},
		{
			name:        "Get approvals list failure - service error",
			queryParams: "",
			mockSetup: func(m *mocks.MockApprovalService) {
				m.On("GetApprovalsWithCursor", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.ErrInternalError)
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeInternalError, resp.Error.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := setupTestApprovalControllerDependencies()
			tt.mockSetup(deps.approvalService)

			controller := controllers.NewApprovalController(deps.logger, deps.approvalService)

			url := "/approvals"
			if tt.queryParams != "" {
				url += "?" + tt.queryParams
			}
			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			controller.GetApprovals(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)

			deps.approvalService.AssertExpectations(t)
		})
	}
}

// ========================================
// GetApproval Tests
// ========================================

func TestApprovalController_GetApproval(t *testing.T) {
	tests := []struct {
		name           string
		approvalID     string
		mockSetup      func(*mocks.MockApprovalService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:       "Get single approval success - pending",
			approvalID: "1",
			mockSetup: func(m *mocks.MockApprovalService) {
				m.On("GetApprovalByID", mock.Anything, uint(1)).Return(
					createTestApproval(1, "Creator Name", constants.ApprovalStatusPending, nil), nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp dtos.ApprovalResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, uint(1), resp.ID)
				assert.Equal(t, "Creator Name", resp.CreatorName)
				assert.Equal(t, constants.ApprovalStatusPending, resp.Status)
				assert.Nil(t, resp.ApproverName)
			},
		},
		{
			name:       "Get single approval success - approved with approver",
			approvalID: "2",
			mockSetup: func(m *mocks.MockApprovalService) {
				approverName := "Approver Name"
				m.On("GetApprovalByID", mock.Anything, uint(2)).Return(
					createTestApproval(2, "Creator", constants.ApprovalStatusApproved, &approverName), nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp dtos.ApprovalResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, uint(2), resp.ID)
				assert.Equal(t, constants.ApprovalStatusApproved, resp.Status)
				assert.NotNil(t, resp.ApproverName)
				assert.Equal(t, "Approver Name", *resp.ApproverName)
			},
		},
		{
			name:           "Get single approval failure - id is 0",
			approvalID:     "0",
			mockSetup:      func(m *mocks.MockApprovalService) {},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeValidationFailed, resp.Error.Code)
				assert.Contains(t, resp.Error.Details, "id")
			},
		},
		{
			name:           "Get single approval failure - invalid id format",
			approvalID:     "abc",
			mockSetup:      func(m *mocks.MockApprovalService) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeInvalidInput, resp.Error.Code)
			},
		},
		{
			name:       "Get single approval failure - approval not found",
			approvalID: "999",
			mockSetup: func(m *mocks.MockApprovalService) {
				m.On("GetApprovalByID", mock.Anything, uint(999)).Return(nil, errors.ErrNotFound)
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeNotFound, resp.Error.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := setupTestApprovalControllerDependencies()
			tt.mockSetup(deps.approvalService)

			controller := controllers.NewApprovalController(deps.logger, deps.approvalService)

			req := httptest.NewRequest(http.MethodGet, "/approvals/"+tt.approvalID, nil)
			req = withApprovalURLParams(req, map[string]string{"id": tt.approvalID})
			w := httptest.NewRecorder()

			controller.GetApproval(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)

			deps.approvalService.AssertExpectations(t)
		})
	}
}

// ========================================
// AddApproval Tests
// ========================================

func TestApprovalController_AddApproval(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*mocks.MockApprovalService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Add approval success",
			mockSetup: func(m *mocks.MockApprovalService) {
				m.On("AddApproval", mock.Anything, mock.MatchedBy(func(a *domains.ApprovalCreate) bool {
					return a.Status == constants.ApprovalStatusPending
				})).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Contains(t, w.Body.String(), "approval added successfully")
			},
		},
		{
			name: "Add approval failure - service error",
			mockSetup: func(m *mocks.MockApprovalService) {
				m.On("AddApproval", mock.Anything, mock.Anything).Return(errors.ErrInternalError)
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeInternalError, resp.Error.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := setupTestApprovalControllerDependencies()
			tt.mockSetup(deps.approvalService)

			controller := controllers.NewApprovalController(deps.logger, deps.approvalService)

			req := httptest.NewRequest(http.MethodPost, "/approvals", nil)
			w := httptest.NewRecorder()

			controller.AddApproval(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)

			deps.approvalService.AssertExpectations(t)
		})
	}
}

// ========================================
// ActionApproval Tests
// ========================================

func TestApprovalController_ActionApproval(t *testing.T) {
	tests := []struct {
		name           string
		approvalID     string
		requestBody    any
		mockSetup      func(*mocks.MockApprovalService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:       "Approve approval success",
			approvalID: "1",
			requestBody: dtos.ApprovalAction{
				Action: constants.ApprovalStatusApproved,
			},
			mockSetup: func(m *mocks.MockApprovalService) {
				m.On("ActionApproval", mock.Anything, uint(1), constants.ApprovalStatusApproved).Return(nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Contains(t, w.Body.String(), "approval actioned successfully")
			},
		},
		{
			name:       "Reject approval success",
			approvalID: "1",
			requestBody: dtos.ApprovalAction{
				Action: constants.ApprovalStatusRejected,
			},
			mockSetup: func(m *mocks.MockApprovalService) {
				m.On("ActionApproval", mock.Anything, uint(1), constants.ApprovalStatusRejected).Return(nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Contains(t, w.Body.String(), "approval actioned successfully")
			},
		},
		{
			name:           "Approval action failure - invalid id",
			approvalID:     "0",
			requestBody:    dtos.ApprovalAction{Action: constants.ApprovalStatusApproved},
			mockSetup:      func(m *mocks.MockApprovalService) {},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeValidationFailed, resp.Error.Code)
			},
		},
		{
			name:           "Approval action failure - invalid JSON",
			approvalID:     "1",
			requestBody:    "invalid json",
			mockSetup:      func(m *mocks.MockApprovalService) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeInvalidInput, resp.Error.Code)
			},
		},
		{
			name:       "Approval action failure - invalid action (pending)",
			approvalID: "1",
			requestBody: dtos.ApprovalAction{
				Action: constants.ApprovalStatusPending, // Invalid - can only be approved or rejected
			},
			mockSetup:      func(m *mocks.MockApprovalService) {},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeValidationFailed, resp.Error.Code)
				assert.Contains(t, resp.Error.Details, "action")
			},
		},
		{
			name:       "Approval action failure - empty action",
			approvalID: "1",
			requestBody: dtos.ApprovalAction{
				Action: "",
			},
			mockSetup:      func(m *mocks.MockApprovalService) {},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeValidationFailed, resp.Error.Code)
				assert.Contains(t, resp.Error.Details, "action")
			},
		},
		{
			name:       "Approval action failure - approval not found",
			approvalID: "999",
			requestBody: dtos.ApprovalAction{
				Action: constants.ApprovalStatusApproved,
			},
			mockSetup: func(m *mocks.MockApprovalService) {
				m.On("ActionApproval", mock.Anything, uint(999), constants.ApprovalStatusApproved).Return(errors.ErrNotFound)
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeNotFound, resp.Error.Code)
			},
		},
		{
			name:       "Approval action failure - insufficient permissions",
			approvalID: "1",
			requestBody: dtos.ApprovalAction{
				Action: constants.ApprovalStatusApproved,
			},
			mockSetup: func(m *mocks.MockApprovalService) {
				m.On("ActionApproval", mock.Anything, uint(1), constants.ApprovalStatusApproved).Return(errors.ErrForbidden)
			},
			expectedStatus: http.StatusForbidden,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeForbidden, resp.Error.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := setupTestApprovalControllerDependencies()
			tt.mockSetup(deps.approvalService)

			controller := controllers.NewApprovalController(deps.logger, deps.approvalService)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/approvals/"+tt.approvalID+"/action", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req = withApprovalURLParams(req, map[string]string{"id": tt.approvalID})
			w := httptest.NewRecorder()

			controller.ActionApproval(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)

			deps.approvalService.AssertExpectations(t)
		})
	}
}

// ========================================
// Helper Functions
// ========================================

// Helper function to create a test approval with edges
func createTestApproval(id uint, creatorName string, status constants.ApprovalStatus, approverName *string) *entgen.Approval {
	now := time.Now()
	approval := &entgen.Approval{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Status:    status,
		CreatorID: 1,
	}
	approval.Edges.Creator = &entgen.User{
		ID:   1,
		Name: creatorName,
	}
	if approverName != nil {
		approval.Edges.Approver = &entgen.User{
			ID:   2,
			Name: *approverName,
		}
		approverId := uint(2)
		approval.ApproverID = &approverId
	}
	return approval
}

// Helper function to setup chi router context with URL params for approval tests
func withApprovalURLParams(r *http.Request, params map[string]string) *http.Request {
	rctx := chi.NewRouteContext()
	for key, val := range params {
		rctx.URLParams.Add(key, val)
	}
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}
