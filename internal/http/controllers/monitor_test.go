package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/dtos"
	"github.com/shunwuse/go-hris/internal/http/controllers"
	"github.com/shunwuse/go-hris/internal/mocks"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ========================================
// Test Setup
// ========================================

type monitorControllerTestDependencies struct {
	logger         *logger.Logger
	monitorService *mocks.MockMonitorService
}

func setupTestMonitorControllerDependencies() *monitorControllerTestDependencies {
	return &monitorControllerTestDependencies{
		logger:         logger.NewNopLogger(),
		monitorService: new(mocks.MockMonitorService),
	}
}

// ========================================
// HealthCheck Tests
// ========================================

func TestMonitorController_HealthCheck(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*mocks.MockMonitorService)
		expectedStatus int
		expectedHealth dtos.HealthResponse
	}{
		{
			name: "Health check success - all services up",
			mockSetup: func(m *mocks.MockMonitorService) {
				m.On("HealthCheck", mock.Anything).Return(&domains.Health{
					Status: constants.StatusUp,
					Components: domains.HealthComponents{
						Database: constants.StatusUp,
						Redis:    constants.StatusUp,
					},
					Info: domains.HealthInfo{
						Version:     "1.0.0",
						Environment: "test",
						Uptime:      "1h0m0s",
						InstanceID:  "test-instance",
						Hostname:    "localhost",
						GoVersion:   "go1.21",
					},
				})
			},
			expectedStatus: http.StatusOK,
			expectedHealth: dtos.HealthResponse{
				Status: constants.StatusUp,
				Components: dtos.HealthComponentsResponse{
					Database: constants.StatusUp,
					Redis:    constants.StatusUp,
				},
				Info: dtos.HealthInfoResponse{
					Version:     "1.0.0",
					Environment: "test",
					Uptime:      "1h0m0s",
					InstanceID:  "test-instance",
					Hostname:    "localhost",
					GoVersion:   "go1.21",
				},
			},
		},
		{
			name: "Health check failure - database down",
			mockSetup: func(m *mocks.MockMonitorService) {
				m.On("HealthCheck", mock.Anything).Return(&domains.Health{
					Status: constants.StatusDown,
					Components: domains.HealthComponents{
						Database: constants.StatusDown,
						Redis:    constants.StatusUp,
					},
					Info: domains.HealthInfo{
						Version:     "1.0.0",
						Environment: "test",
						Uptime:      "1h0m0s",
						InstanceID:  "test-instance",
						Hostname:    "localhost",
						GoVersion:   "go1.21",
					},
				})
			},
			expectedStatus: http.StatusServiceUnavailable,
			expectedHealth: dtos.HealthResponse{
				Status: constants.StatusDown,
				Components: dtos.HealthComponentsResponse{
					Database: constants.StatusDown,
					Redis:    constants.StatusUp,
				},
				Info: dtos.HealthInfoResponse{
					Version:     "1.0.0",
					Environment: "test",
					Uptime:      "1h0m0s",
					InstanceID:  "test-instance",
					Hostname:    "localhost",
					GoVersion:   "go1.21",
				},
			},
		},
		{
			name: "Health check failure - redis down",
			mockSetup: func(m *mocks.MockMonitorService) {
				m.On("HealthCheck", mock.Anything).Return(&domains.Health{
					Status: constants.StatusDown,
					Components: domains.HealthComponents{
						Database: constants.StatusUp,
						Redis:    constants.StatusDown,
					},
					Info: domains.HealthInfo{
						Version:     "1.0.0",
						Environment: "test",
						Uptime:      "30m0s",
						InstanceID:  "test-instance-2",
						Hostname:    "localhost",
						GoVersion:   "go1.21",
					},
				})
			},
			expectedStatus: http.StatusServiceUnavailable,
			expectedHealth: dtos.HealthResponse{
				Status: constants.StatusDown,
				Components: dtos.HealthComponentsResponse{
					Database: constants.StatusUp,
					Redis:    constants.StatusDown,
				},
				Info: dtos.HealthInfoResponse{
					Version:     "1.0.0",
					Environment: "test",
					Uptime:      "30m0s",
					InstanceID:  "test-instance-2",
					Hostname:    "localhost",
					GoVersion:   "go1.21",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := setupTestMonitorControllerDependencies()
			tt.mockSetup(deps.monitorService)

			controller := controllers.NewMonitorController(deps.logger, deps.monitorService)

			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			w := httptest.NewRecorder()

			controller.HealthCheck(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response dtos.HealthResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedHealth, response)

			deps.monitorService.AssertExpectations(t)
		})
	}
}
