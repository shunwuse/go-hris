package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/http/response"
	"github.com/shunwuse/go-hris/internal/infra/app"
	"github.com/shunwuse/go-hris/internal/infra/metrics"
	"github.com/shunwuse/go-hris/internal/mocks"
	"github.com/shunwuse/go-hris/internal/pkg/contextx"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
	"github.com/shunwuse/go-hris/internal/ports/infra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type captureAlerter struct {
	messages []infra.Message
}

func (a *captureAlerter) Send(_ context.Context, msg infra.Message) error {
	a.messages = append(a.messages, msg)
	return nil
}

type fakeCommonMiddleware struct {
	called int
}

func (m *fakeCommonMiddleware) Setup(_ chi.Router) {
	m.called++
}

func TestCORSMiddleware_Handler(t *testing.T) {
	mw := NewCORSMiddleware()
	handler := mw.Handler()(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodOptions, "/hello", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.NotEmpty(t, rr.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", rr.Header().Get("Access-Control-Allow-Credentials"))
}

func TestProfilerMiddleware_Handler(t *testing.T) {
	tests := []struct {
		name         string
		configToken  string
		headerToken  string
		expectedCode int
		nextCalled   bool
	}{
		{
			name:         "missing configured token returns forbidden",
			configToken:  "",
			headerToken:  "any",
			expectedCode: http.StatusForbidden,
			nextCalled:   false,
		},
		{
			name:         "invalid header token returns unauthorized",
			configToken:  "secret",
			headerToken:  "wrong",
			expectedCode: http.StatusUnauthorized,
			nextCalled:   false,
		},
		{
			name:         "valid token passes through",
			configToken:  "secret",
			headerToken:  "secret",
			expectedCode: http.StatusNoContent,
			nextCalled:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw := NewProfilerMiddleware(&app.ServiceConfig{ProfilerToken: tt.configToken})
			called := false

			handler := mw.Handler()(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				called = true
				w.WriteHeader(http.StatusNoContent)
			}))

			req := httptest.NewRequest(http.MethodGet, "/debug/pprof", nil)
			req.Header.Set("X-Profiler-Token", tt.headerToken)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			assert.Equal(t, tt.nextCalled, called)
		})
	}
}

func TestTraceMiddleware_Handler(t *testing.T) {
	log := logger.NewNopLogger()
	mw := NewTraceMiddleware(log)

	tests := []struct {
		name         string
		inputTraceID string
	}{
		{name: "uses incoming trace id", inputTraceID: "trace-123"},
		{name: "generates trace id", inputTraceID: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var traceInContext string
			var spanInContext string

			handler := mw.Handler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				traceInContext = contextx.GetTraceID(r.Context())
				spanInContext = contextx.GetSpanID(r.Context())
				w.WriteHeader(http.StatusNoContent)
			}))

			req := httptest.NewRequest(http.MethodGet, "/trace", nil)
			if tt.inputTraceID != "" {
				req.Header.Set("X-Trace-Id", tt.inputTraceID)
			}
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusNoContent, rr.Code)
			assert.NotEmpty(t, rr.Header().Get("X-Span-Id"))
			assert.NotEmpty(t, traceInContext)
			assert.NotEmpty(t, spanInContext)

			if tt.inputTraceID != "" {
				assert.Equal(t, tt.inputTraceID, rr.Header().Get("X-Trace-Id"))
				assert.Equal(t, tt.inputTraceID, traceInContext)
			}
		})
	}
}

func TestJWTMiddleware_Handler(t *testing.T) {
	log := logger.NewNopLogger()

	tests := []struct {
		name           string
		authorization  string
		mockSetup      func(*mocks.MockAuthService)
		expectedStatus int
		nextCalled     bool
	}{
		{
			name:           "missing authorization header",
			authorization:  "",
			mockSetup:      func(_ *mocks.MockAuthService) {},
			expectedStatus: http.StatusUnauthorized,
			nextCalled:     false,
		},
		{
			name:           "invalid header format",
			authorization:  "Bearer token extra",
			mockSetup:      func(_ *mocks.MockAuthService) {},
			expectedStatus: http.StatusUnauthorized,
			nextCalled:     false,
		},
		{
			name:           "invalid scheme",
			authorization:  "Token abc",
			mockSetup:      func(_ *mocks.MockAuthService) {},
			expectedStatus: http.StatusUnauthorized,
			nextCalled:     false,
		},
		{
			name:          "token validation failed",
			authorization: "Bearer invalid",
			mockSetup: func(s *mocks.MockAuthService) {
				s.On("ValidateAccessToken", mock.Anything, "invalid").Return(nil, errors.ErrTokenInvalid).Once()
			},
			expectedStatus: http.StatusUnauthorized,
			nextCalled:     false,
		},
		{
			name:          "token validation success",
			authorization: "Bearer valid",
			mockSetup: func(s *mocks.MockAuthService) {
				claims := &domains.Claims{
					JTI:       "jti-1",
					ExpiresAt: time.Now().Add(time.Hour),
					Identity: domains.Identity{
						UserID:   1,
						Username: "admin",
						Roles:    []constants.Role{constants.Admin},
					},
				}
				s.On("ValidateAccessToken", mock.Anything, "valid").Return(claims, nil).Once()
			},
			expectedStatus: http.StatusNoContent,
			nextCalled:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authSvc := new(mocks.MockAuthService)
			tt.mockSetup(authSvc)

			mw := NewJWTMiddleware(log, authSvc)
			called := false

			handler := mw.Handler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				claims, ok := contextx.GetClaims(r.Context())
				require.True(t, ok)
				require.NotNil(t, claims)
				w.WriteHeader(http.StatusNoContent)
			}))

			req := httptest.NewRequest(http.MethodGet, "/secure", nil)
			if tt.authorization != "" {
				req.Header.Set("Authorization", tt.authorization)
			}
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.nextCalled, called)
			authSvc.AssertExpectations(t)
		})
	}
}

func TestExceptionMiddleware_Handler(t *testing.T) {
	t.Run("sends alert on 5xx", func(t *testing.T) {
		alerter := &captureAlerter{}
		mw := NewExceptionMiddleware(alerter)

		handler := mw.Handler()(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))

		req := httptest.NewRequest(http.MethodGet, "/boom", nil)
		req = req.WithContext(contextx.WithTraceID(req.Context(), "trace-500"))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		require.Len(t, alerter.messages, 1)
		msg := alerter.messages[0]
		assert.Equal(t, infra.LevelError, msg.Level)
		assert.Equal(t, "trace-500", msg.TraceID)
		assert.Equal(t, "HTTP 500 Error", msg.Title)
		assert.Contains(t, msg.Content, "Method: GET")
	})

	t.Run("does not send alert on non-5xx", func(t *testing.T) {
		alerter := &captureAlerter{}
		mw := NewExceptionMiddleware(alerter)

		handler := mw.Handler()(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Empty(t, alerter.messages)
	})
}

func TestRecoveryMiddleware_Handler(t *testing.T) {
	alerter := &captureAlerter{}
	mw := NewRecoveryMiddleware(logger.NewNopLogger(), alerter)

	handler := mw.Handler()(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		panic("boom")
	}))

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	req = req.WithContext(contextx.WithTraceID(req.Context(), "trace-panic"))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "close", rr.Header().Get("Connection"))

	var body response.ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, errors.CodeInternalError, body.Error.Code)

	require.Len(t, alerter.messages, 1)
	msg := alerter.messages[0]
	assert.Equal(t, infra.LevelCritical, msg.Level)
	assert.Equal(t, "trace-panic", msg.TraceID)
	assert.Equal(t, "Panic Recovered", msg.Title)
	assert.NotEmpty(t, msg.StackTrace)
}

func TestRecoveryMiddleware_UpgradeConnection(t *testing.T) {
	alerter := &captureAlerter{}
	mw := NewRecoveryMiddleware(logger.NewNopLogger(), alerter)

	handler := mw.Handler()(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		panic("upgrade panic")
	}))

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	req.Header.Set("Connection", "Upgrade")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.NotEqual(t, "close", rr.Header().Get("Connection"))
}

func TestRecoveryMiddleware_AbortHandlerPanics(t *testing.T) {
	alerter := &captureAlerter{}
	mw := NewRecoveryMiddleware(logger.NewNopLogger(), alerter)

	handler := mw.Handler()(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		panic(http.ErrAbortHandler)
	}))

	req := httptest.NewRequest(http.MethodGet, "/abort", nil)
	rr := httptest.NewRecorder()

	assert.Panics(t, func() {
		handler.ServeHTTP(rr, req)
	})
}

func TestMetricsMiddleware_Handler(t *testing.T) {
	metricSet := &metrics.Metrics{
		HttpRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{Name: "test_http_requests_total", Help: "test"},
			[]string{"method", "path", "status"},
		),
		HttpRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{Name: "test_http_request_duration_seconds", Help: "test"},
			[]string{"method", "path"},
		),
	}

	mw := NewMetricsMiddleware(metricSet)

	handler := mw.Handler()(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	rctx := chi.NewRouteContext()
	rctx.RoutePatterns = append(rctx.RoutePatterns, "/users/{id}")

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestMetricsMiddleware_GetPatternFallback(t *testing.T) {
	metricSet := &metrics.Metrics{
		HttpRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{Name: "test_http_requests_total_fallback", Help: "test"},
			[]string{"method", "path", "status"},
		),
		HttpRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{Name: "test_http_request_duration_seconds_fallback", Help: "test"},
			[]string{"method", "path"},
		),
	}

	mw := NewMetricsMiddleware(metricSet)

	rctx := chi.NewRouteContext()
	req := httptest.NewRequest(http.MethodGet, "/fallback", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	assert.Equal(t, "/fallback", mw.getPattern(req))
}

func TestRequestLoggerMiddleware_Handler(t *testing.T) {
	mw := NewRequestLoggerMiddleware(logger.NewNopLogger())

	statuses := []int{http.StatusOK, http.StatusBadRequest, http.StatusInternalServerError}
	for _, status := range statuses {
		t.Run(http.StatusText(status), func(t *testing.T) {
			handler := mw.Handler()(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(status)
			}))

			req := httptest.NewRequest(http.MethodGet, "/log", nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, status, rr.Code)
		})
	}
}

func TestIdempotencyMiddleware_SkipsSafeMethods(t *testing.T) {
	mw := NewIdempotencyMiddleware(
		&app.ServiceConfig{IdempotencyExpireMinutes: 1},
		nil,
		logger.NewNopLogger(),
	)

	methods := []string{http.MethodGet, http.MethodHead, http.MethodOptions}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			called := false
			handler := mw.HandlerWithTTL(time.Minute)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				called = true
				w.WriteHeader(http.StatusNoContent)
			}))

			req := httptest.NewRequest(method, "/resource", nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.True(t, called)
			assert.Equal(t, http.StatusNoContent, rr.Code)
		})
	}
}

func TestIdempotencyMiddleware_SkipsWhenTraceIDMissing(t *testing.T) {
	mw := NewIdempotencyMiddleware(
		&app.ServiceConfig{IdempotencyExpireMinutes: 1},
		nil,
		logger.NewNopLogger(),
	)

	called := false
	handler := mw.HandlerWithTTL(time.Minute)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodPost, "/resource", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestCommonMiddlewares_Setup(t *testing.T) {
	first := &fakeCommonMiddleware{}
	second := &fakeCommonMiddleware{}

	common := CommonMiddlewares{first, second}
	router := chi.NewRouter()
	common.Setup(router)

	router.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	assert.Equal(t, 1, first.called)
	assert.Equal(t, 1, second.called)

	req := httptest.NewRequest(http.MethodOptions, "/ping", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "GET")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.NotEmpty(t, rr.Header().Get("Access-Control-Allow-Origin"))
}

func TestNewCommonMiddlewares_Order(t *testing.T) {
	metricSet := &metrics.Metrics{
		HttpRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{Name: "test_common_middleware_order_http_requests_total", Help: "test"},
			[]string{"method", "path", "status"},
		),
		HttpRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{Name: "test_common_middleware_order_http_request_duration_seconds", Help: "test"},
			[]string{"method", "path"},
		),
	}

	trace := NewTraceMiddleware(logger.NewNopLogger())
	metricsMiddleware := NewMetricsMiddleware(metricSet)
	requestLoggerMiddleware := NewRequestLoggerMiddleware(logger.NewNopLogger())
	recoveryMiddleware := NewRecoveryMiddleware(logger.NewNopLogger(), &captureAlerter{})
	idempotencyMiddleware := NewIdempotencyMiddleware(
		&app.ServiceConfig{IdempotencyExpireMinutes: 1},
		nil,
		logger.NewNopLogger(),
	)
	exceptionMiddleware := NewExceptionMiddleware(&captureAlerter{})

	common := NewCommonMiddlewares(
		trace,
		metricsMiddleware,
		requestLoggerMiddleware,
		recoveryMiddleware,
		idempotencyMiddleware,
		exceptionMiddleware,
	)

	require.Len(t, common, 6)
	assert.Same(t, trace, common[0])
	assert.Same(t, metricsMiddleware, common[1])
	assert.Same(t, requestLoggerMiddleware, common[2])
	assert.Same(t, recoveryMiddleware, common[3])
	assert.Same(t, idempotencyMiddleware, common[4])
	assert.Same(t, exceptionMiddleware, common[5])
}
