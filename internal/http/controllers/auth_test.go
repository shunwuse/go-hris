package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/dtos"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/http/controllers"
	"github.com/shunwuse/go-hris/internal/http/response"
	"github.com/shunwuse/go-hris/internal/mocks"
	"github.com/shunwuse/go-hris/internal/pkg/contextx"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ========================================
// Test Setup
// ========================================

type testAuthControllerDependencies struct {
	logger      *logger.Logger
	authService *mocks.MockAuthService
}

func setupTestAuthControllerDependencies() *testAuthControllerDependencies {
	return &testAuthControllerDependencies{
		logger:      logger.NewNopLogger(),
		authService: new(mocks.MockAuthService),
	}
}

// ========================================
// Login Tests
// ========================================

func TestAuthController_Login(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    any
		mockSetup      func(*mocks.MockAuthService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Login success",
			requestBody: dtos.UserLogin{
				Username: "admin",
				Password: "password123",
			},
			mockSetup: func(m *mocks.MockAuthService) {
				m.On("Login", mock.Anything, "admin", "password123").Return(&domains.LoginResult{
					Username:     "admin",
					Roles:        []constants.Role{constants.Admin},
					AccessToken:  "test-access-token",
					RefreshToken: "test-refresh-token",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp dtos.UserLoginResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, "admin", resp.Username)
				assert.Contains(t, resp.Roles, constants.Admin)
				assert.NotEmpty(t, resp.AccessToken)
				assert.NotEmpty(t, resp.RefreshToken)
			},
		},
		{
			name:           "Login failure - invalid JSON",
			requestBody:    "invalid json",
			mockSetup:      func(m *mocks.MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeInvalidInput, resp.Error.Code)
			},
		},
		{
			name: "Login failure - missing username",
			requestBody: dtos.UserLogin{
				Username: "",
				Password: "password123",
			},
			mockSetup:      func(m *mocks.MockAuthService) {},
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
			name: "Login failure - missing password",
			requestBody: dtos.UserLogin{
				Username: "admin",
				Password: "",
			},
			mockSetup:      func(m *mocks.MockAuthService) {},
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
			name: "Login failure - missing username and password",
			requestBody: dtos.UserLogin{
				Username: "",
				Password: "",
			},
			mockSetup:      func(m *mocks.MockAuthService) {},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeValidationFailed, resp.Error.Code)
				assert.Contains(t, resp.Error.Details, "username")
				assert.Contains(t, resp.Error.Details, "password")
			},
		},
		{
			name: "Login failure - invalid credentials",
			requestBody: dtos.UserLogin{
				Username: "admin",
				Password: "wrongpassword",
			},
			mockSetup: func(m *mocks.MockAuthService) {
				m.On("Login", mock.Anything, "admin", "wrongpassword").Return(nil, errors.ErrInvalidCredentials)
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeInvalidCredentials, resp.Error.Code)
			},
		},
		{
			name: "Login failure - user not found",
			requestBody: dtos.UserLogin{
				Username: "nonexistent",
				Password: "password123",
			},
			mockSetup: func(m *mocks.MockAuthService) {
				m.On("Login", mock.Anything, "nonexistent", "password123").Return(nil, errors.ErrNotFound)
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
			deps := setupTestAuthControllerDependencies()
			tt.mockSetup(deps.authService)

			controller := controllers.NewAuthController(deps.logger, deps.authService)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			controller.Login(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)

			deps.authService.AssertExpectations(t)
		})
	}
}

// ========================================
// RefreshToken Tests
// ========================================

func TestAuthController_RefreshToken(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    any
		mockSetup      func(*mocks.MockAuthService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Refresh success",
			requestBody: dtos.RefreshRequest{
				RefreshToken: "valid-refresh-token",
			},
			mockSetup: func(m *mocks.MockAuthService) {
				m.On("RefreshAccessToken", mock.Anything, "valid-refresh-token").Return(&domains.TokenPair{
					AccessToken:  "new-access-token",
					RefreshToken: "new-refresh-token",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp dtos.RefreshResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, "new-access-token", resp.AccessToken)
				assert.Equal(t, "new-refresh-token", resp.RefreshToken)
			},
		},
		{
			name:           "Refresh failure - invalid JSON",
			requestBody:    "invalid json",
			mockSetup:      func(m *mocks.MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeInvalidInput, resp.Error.Code)
			},
		},
		{
			name: "Refresh failure - missing refresh_token",
			requestBody: dtos.RefreshRequest{
				RefreshToken: "",
			},
			mockSetup:      func(m *mocks.MockAuthService) {},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeValidationFailed, resp.Error.Code)
				assert.Contains(t, resp.Error.Details, "refresh_token")
			},
		},
		{
			name: "Refresh failure - invalid token",
			requestBody: dtos.RefreshRequest{
				RefreshToken: "invalid-token",
			},
			mockSetup: func(m *mocks.MockAuthService) {
				m.On("RefreshAccessToken", mock.Anything, "invalid-token").Return(nil, errors.ErrTokenInvalid)
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeTokenInvalid, resp.Error.Code)
			},
		},
		{
			name: "Refresh failure - token expired",
			requestBody: dtos.RefreshRequest{
				RefreshToken: "expired-token",
			},
			mockSetup: func(m *mocks.MockAuthService) {
				m.On("RefreshAccessToken", mock.Anything, "expired-token").Return(nil, errors.ErrTokenExpired)
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeTokenExpired, resp.Error.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := setupTestAuthControllerDependencies()
			tt.mockSetup(deps.authService)

			controller := controllers.NewAuthController(deps.logger, deps.authService)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			controller.RefreshToken(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)

			deps.authService.AssertExpectations(t)
		})
	}
}

// ========================================
// Logout Tests
// ========================================

func TestAuthController_Logout(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    any
		setupContext   func(context.Context) context.Context
		mockSetup      func(*mocks.MockAuthService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Logout success - with refresh_token",
			requestBody: dtos.LogoutRequest{
				RefreshToken: "valid-refresh-token",
			},
			setupContext: func(ctx context.Context) context.Context {
				claims := &domains.Claims{
					JTI:       "test-jti",
					ExpiresAt: time.Now().Add(time.Hour),
					Identity: domains.Identity{
						UserID:   1,
						Username: "admin",
					},
				}
				return contextx.WithClaims(ctx, claims)
			},
			mockSetup: func(m *mocks.MockAuthService) {
				m.On("RevokeRefreshToken", mock.Anything, "valid-refresh-token").Return(nil)
				m.On("BlacklistToken", mock.Anything, "test-jti", mock.AnythingOfType("time.Duration")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Contains(t, w.Body.String(), "logged out successfully")
			},
		},
		{
			name:        "Logout success - without refresh_token",
			requestBody: dtos.LogoutRequest{},
			setupContext: func(ctx context.Context) context.Context {
				claims := &domains.Claims{
					JTI:       "test-jti",
					ExpiresAt: time.Now().Add(time.Hour),
					Identity: domains.Identity{
						UserID:   1,
						Username: "admin",
					},
				}
				return contextx.WithClaims(ctx, claims)
			},
			mockSetup: func(m *mocks.MockAuthService) {
				m.On("BlacklistToken", mock.Anything, "test-jti", mock.AnythingOfType("time.Duration")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Contains(t, w.Body.String(), "logged out successfully")
			},
		},
		{
			name: "Logout success - no claims (edge case)",
			requestBody: dtos.LogoutRequest{
				RefreshToken: "valid-refresh-token",
			},
			setupContext: func(ctx context.Context) context.Context {
				return ctx // No claims set
			},
			mockSetup: func(m *mocks.MockAuthService) {
				m.On("RevokeRefreshToken", mock.Anything, "valid-refresh-token").Return(nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Contains(t, w.Body.String(), "logged out successfully")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := setupTestAuthControllerDependencies()
			tt.mockSetup(deps.authService)

			controller := controllers.NewAuthController(deps.logger, deps.authService)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/auth/logout", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			if tt.setupContext != nil {
				req = req.WithContext(tt.setupContext(req.Context()))
			}

			w := httptest.NewRecorder()

			controller.Logout(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)

			deps.authService.AssertExpectations(t)
		})
	}
}

// ========================================
// LogoutAll Tests
// ========================================

func TestAuthController_LogoutAll(t *testing.T) {
	tests := []struct {
		name           string
		setupContext   func(context.Context) context.Context
		mockSetup      func(*mocks.MockAuthService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Logout all devices success",
			setupContext: func(ctx context.Context) context.Context {
				claims := &domains.Claims{
					JTI:       "test-jti",
					ExpiresAt: time.Now().Add(time.Hour),
					Identity: domains.Identity{
						UserID:   1,
						Username: "admin",
					},
				}
				return contextx.WithClaims(ctx, claims)
			},
			mockSetup: func(m *mocks.MockAuthService) {
				m.On("RevokeAllUserTokens", mock.Anything, uint(1)).Return(nil)
				m.On("BlacklistToken", mock.Anything, "test-jti", mock.AnythingOfType("time.Duration")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Contains(t, w.Body.String(), "logged out from all devices successfully")
			},
		},
		{
			name: "Logout all devices failure - no claims",
			setupContext: func(ctx context.Context) context.Context {
				return ctx // No claims set
			},
			mockSetup:      func(m *mocks.MockAuthService) {},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp response.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, errors.CodeUnauthorized, resp.Error.Code)
			},
		},
		{
			name: "Logout all devices failure - service error",
			setupContext: func(ctx context.Context) context.Context {
				claims := &domains.Claims{
					JTI:       "test-jti",
					ExpiresAt: time.Now().Add(time.Hour),
					Identity: domains.Identity{
						UserID:   1,
						Username: "admin",
					},
				}
				return contextx.WithClaims(ctx, claims)
			},
			mockSetup: func(m *mocks.MockAuthService) {
				m.On("RevokeAllUserTokens", mock.Anything, uint(1)).Return(errors.ErrInternalError)
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
			deps := setupTestAuthControllerDependencies()
			tt.mockSetup(deps.authService)

			controller := controllers.NewAuthController(deps.logger, deps.authService)

			req := httptest.NewRequest(http.MethodPost, "/auth/logout-all", nil)

			if tt.setupContext != nil {
				req = req.WithContext(tt.setupContext(req.Context()))
			}

			w := httptest.NewRecorder()

			controller.LogoutAll(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)

			deps.authService.AssertExpectations(t)
		})
	}
}
