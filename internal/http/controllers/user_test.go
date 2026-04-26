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
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/dtos"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/http/controllers"
	"github.com/shunwuse/go-hris/internal/http/response"
	"github.com/shunwuse/go-hris/internal/mocks"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ========================================
// Test Setup
// ========================================

type testUserControllerDependencies struct {
	logger      *logger.Logger
	userService *mocks.MockUserService
}

func setupTestUserControllerDependencies() *testUserControllerDependencies {
	return &testUserControllerDependencies{
		logger:      logger.NewNopLogger(),
		userService: &mocks.MockUserService{},
	}
}

// ========================================
// GetUsers Tests
// ========================================

func TestUserController_GetUsers(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		mockSetup      func(*mocks.MockUserService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:        "Get users list success - default pagination",
			queryParams: "",
			mockSetup: func(m *mocks.MockUserService) {
				m.On("GetUsersWithOffset", mock.Anything, mock.Anything, mock.Anything).Return(&domains.OffsetResult[domains.User]{
					Items: []domains.User{
						createTestUser(1, "admin", "Admin User"),
						createTestUser(2, "user1", "Test User"),
					},
					TotalCount: 2,
					TotalPage:  1,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.OffsetListResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, 2, resp.Meta.Total)
				assert.Equal(t, 1, resp.Meta.LastPage)
			},
		},
		{
			name:        "Get users list success - with pagination params",
			queryParams: "page=2&per_page=5",
			mockSetup: func(m *mocks.MockUserService) {
				m.On("GetUsersWithOffset", mock.Anything, mock.MatchedBy(func(q domains.OffsetQuery) bool {
					return q.Page == 2 && q.PerPage == 5
				}), mock.Anything).Return(&domains.OffsetResult[domains.User]{
					Items:      []domains.User{},
					TotalCount: 10,
					TotalPage:  2,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.OffsetListResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, 10, resp.Meta.Total)
				assert.Equal(t, 2, resp.Meta.CurrentPage)
				assert.Equal(t, 5, resp.Meta.PerPage)
			},
		},
		{
			name:        "Get users list success - with filters",
			queryParams: "name=test&role=manager",
			mockSetup: func(m *mocks.MockUserService) {
				m.On("GetUsersWithOffset", mock.Anything, mock.Anything, mock.MatchedBy(func(f domains.UserFilter) bool {
					return f.Name == "test" && f.Role == constants.Manager
				})).Return(&domains.OffsetResult[domains.User]{
					Items:      []domains.User{createTestUser(3, "manager1", "Test Manager")},
					TotalCount: 1,
					TotalPage:  1,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.OffsetListResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, 1, resp.Meta.Total)
			},
		},
		{
			name:           "Get users list failure - invalid role",
			queryParams:    "role=invalid_role",
			mockSetup:      func(m *mocks.MockUserService) {},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeValidationFailed, resp.Error.Code)
			},
		},
		{
			name:        "Get users list failure - service error",
			queryParams: "",
			mockSetup: func(m *mocks.MockUserService) {
				m.On("GetUsersWithOffset", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.ErrInternalError)
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
			deps := setupTestUserControllerDependencies()
			tt.mockSetup(deps.userService)

			controller := controllers.NewUserController(deps.logger, deps.userService)

			url := "/users"
			if tt.queryParams != "" {
				url += "?" + tt.queryParams
			}
			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			controller.GetUsers(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)

			deps.userService.AssertExpectations(t)
		})
	}
}

// =======================================
// GetUser Tests
// =======================================

func TestUserController_GetUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockSetup      func(*mocks.MockUserService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "Get single user success",
			userID: "1",
			mockSetup: func(m *mocks.MockUserService) {
				m.On("GetUserByID", mock.Anything, uint(1)).Return(createTestUserWithPermissions(1, "admin", "Admin User"), nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp dtos.UserResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, uint(1), resp.ID)
				assert.Equal(t, "admin", resp.Username)
				assert.Equal(t, "Admin User", resp.Name)
			},
		},
		{
			name:           "Get single user failure - id is 0",
			userID:         "0",
			mockSetup:      func(m *mocks.MockUserService) {},
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
			name:           "Get single user failure - invalid id format",
			userID:         "abc",
			mockSetup:      func(m *mocks.MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeInvalidInput, resp.Error.Code)
			},
		},
		{
			name:   "Get single user failure - user not found",
			userID: "999",
			mockSetup: func(m *mocks.MockUserService) {
				m.On("GetUserByID", mock.Anything, uint(999)).Return(nil, errors.ErrNotFound)
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
			deps := setupTestUserControllerDependencies()
			tt.mockSetup(deps.userService)

			controller := controllers.NewUserController(deps.logger, deps.userService)

			req := httptest.NewRequest(http.MethodGet, "/users/"+tt.userID, nil)
			req = withURLParams(req, map[string]string{"id": tt.userID})
			w := httptest.NewRecorder()

			controller.GetUser(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)

			deps.userService.AssertExpectations(t)
		})
	}
}

// =======================================
// CreateUser Tests
// =======================================

func TestUserController_CreateUser(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    any
		mockSetup      func(*mocks.MockUserService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Create user success",
			requestBody: dtos.UserCreate{
				Username: "newuser",
				Password: "password123",
				Name:     "New User",
				Role:     constants.Staff,
			},
			mockSetup: func(m *mocks.MockUserService) {
				m.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *domains.UserCreate) bool {
					return u.Username == "newuser" && u.Name == "New User" && u.Password == "password123"
				}), constants.Staff).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Contains(t, w.Body.String(), "user created successfully")
			},
		},
		{
			name:           "Create user failure - invalid JSON",
			requestBody:    "invalid json",
			mockSetup:      func(m *mocks.MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeInvalidInput, resp.Error.Code)
			},
		},
		{
			name: "Create user failure - missing username",
			requestBody: dtos.UserCreate{
				Username: "",
				Password: "password123",
				Name:     "New User",
				Role:     constants.Staff,
			},
			mockSetup:      func(m *mocks.MockUserService) {},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeValidationFailed, resp.Error.Code)
				assert.Contains(t, resp.Error.Details, "username")
			},
		},
		{
			name: "Create user failure - password too short",
			requestBody: dtos.UserCreate{
				Username: "newuser",
				Password: "12345", // less than 6 chars
				Name:     "New User",
				Role:     constants.Staff,
			},
			mockSetup:      func(m *mocks.MockUserService) {},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeValidationFailed, resp.Error.Code)
				assert.Contains(t, resp.Error.Details, "password")
			},
		},
		{
			name: "Create user failure - missing name",
			requestBody: dtos.UserCreate{
				Username: "newuser",
				Password: "password123",
				Name:     "",
				Role:     constants.Staff,
			},
			mockSetup:      func(m *mocks.MockUserService) {},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeValidationFailed, resp.Error.Code)
				assert.Contains(t, resp.Error.Details, "name")
			},
		},
		{
			name: "Create user failure - missing role",
			requestBody: dtos.UserCreate{
				Username: "newuser",
				Password: "password123",
				Name:     "New User",
				Role:     "",
			},
			mockSetup:      func(m *mocks.MockUserService) {},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeValidationFailed, resp.Error.Code)
				assert.Contains(t, resp.Error.Details, "role")
			},
		},
		{
			name: "Create user failure - invalid role (admin)",
			requestBody: dtos.UserCreate{
				Username: "newuser",
				Password: "password123",
				Name:     "New User",
				Role:     constants.Admin, // Cannot create admin user
			},
			mockSetup:      func(m *mocks.MockUserService) {},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeValidationFailed, resp.Error.Code)
				assert.Contains(t, resp.Error.Details, "role")
			},
		},
		{
			name: "Create user failure - username exists",
			requestBody: dtos.UserCreate{
				Username: "existinguser",
				Password: "password123",
				Name:     "Existing User",
				Role:     constants.Staff,
			},
			mockSetup: func(m *mocks.MockUserService) {
				m.On("CreateUser", mock.Anything, mock.Anything, constants.Staff).Return(errors.ErrAlreadyExists)
			},
			expectedStatus: http.StatusConflict,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeAlreadyExists, resp.Error.Code)
			},
		},
		{
			name: "Create user failure - multiple validation errors",
			requestBody: dtos.UserCreate{
				Username: "",
				Password: "",
				Name:     "",
				Role:     "",
			},
			mockSetup:      func(m *mocks.MockUserService) {},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeValidationFailed, resp.Error.Code)
				assert.Contains(t, resp.Error.Details, "username")
				assert.Contains(t, resp.Error.Details, "password")
				assert.Contains(t, resp.Error.Details, "name")
				assert.Contains(t, resp.Error.Details, "role")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := setupTestUserControllerDependencies()
			tt.mockSetup(deps.userService)

			controller := controllers.NewUserController(deps.logger, deps.userService)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			controller.CreateUser(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)

			deps.userService.AssertExpectations(t)
		})
	}
}

// =======================================
// UpdateUser Tests
// =======================================

func TestUserController_UpdateUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		requestBody    any
		mockSetup      func(*mocks.MockUserService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "Update user success",
			userID: "1",
			requestBody: dtos.UserUpdate{
				Name: "Updated Name",
			},
			mockSetup: func(m *mocks.MockUserService) {
				m.On("UpdateUser", mock.Anything, mock.MatchedBy(func(u *domains.UserUpdate) bool {
					return u.ID == 1 && u.Name == "Updated Name"
				})).Return(nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Contains(t, w.Body.String(), "user updated successfully")
			},
		},
		{
			name:           "Update user failure - invalid id",
			userID:         "0",
			requestBody:    dtos.UserUpdate{Name: "Updated Name"},
			mockSetup:      func(m *mocks.MockUserService) {},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeValidationFailed, resp.Error.Code)
			},
		},
		{
			name:   "Update user failure - missing name",
			userID: "1",
			requestBody: dtos.UserUpdate{
				Name: "",
			},
			mockSetup:      func(m *mocks.MockUserService) {},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeValidationFailed, resp.Error.Code)
				assert.Contains(t, resp.Error.Details, "name")
			},
		},
		{
			name:           "Update user failure - invalid JSON",
			userID:         "1",
			requestBody:    "invalid json",
			mockSetup:      func(m *mocks.MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeInvalidInput, resp.Error.Code)
			},
		},
		{
			name:   "Update user failure - user not found",
			userID: "999",
			requestBody: dtos.UserUpdate{
				Name: "Updated Name",
			},
			mockSetup: func(m *mocks.MockUserService) {
				m.On("UpdateUser", mock.Anything, mock.Anything).Return(errors.ErrNotFound)
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
			deps := setupTestUserControllerDependencies()
			tt.mockSetup(deps.userService)

			controller := controllers.NewUserController(deps.logger, deps.userService)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/users/"+tt.userID, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req = withURLParams(req, map[string]string{"id": tt.userID})
			w := httptest.NewRecorder()

			controller.UpdateUser(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)

			deps.userService.AssertExpectations(t)
		})
	}
}

// ========================================
// Helper Functions
// ========================================

// Helper function to create a test user (for list operations)
func createTestUser(id uint, username, name string) domains.User {
	now := time.Now()
	return domains.User{
		ID:        id,
		Username:  username,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Helper function to create a test user with permissions (for GetUser operations)
func createTestUserWithPermissions(id uint, username, name string) *domains.UserWithPermissions {
	user := createTestUser(id, username, name)
	return &domains.UserWithPermissions{
		ID:          user.ID,
		Username:    user.Username,
		Name:        user.Name,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Permissions: constants.Permissions{"user:read", "user:write"},
	}
}

// Helper function to setup chi router context with URL params
func withURLParams(r *http.Request, params map[string]string) *http.Request {
	rctx := chi.NewRouteContext()
	for key, val := range params {
		rctx.URLParams.Add(key, val)
	}
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}
