package metrics

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewReturnsIndependentRegistries(t *testing.T) {
	first := New()
	second := New()

	require.NotNil(t, first)
	require.NotNil(t, second)
	assert.NotSame(t, first.registry, second.registry)
}

func TestHandlerExposesRecordedMetrics(t *testing.T) {
	metricSet := New()
	metricSet.HttpRequestsTotal.WithLabelValues("GET", "/health", "200").Inc()
	metricSet.HttpRequestDuration.WithLabelValues("GET", "/health").Observe(0.25)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)

	metricSet.Handler().ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)
	body := recorder.Body.String()
	assert.Contains(t, body, "http_requests_total")
	assert.Contains(t, body, "http_request_duration_seconds")
	assert.Contains(t, body, "go_goroutines")
	assert.Contains(t, body, "process_cpu_seconds_total")
	assert.Contains(t, body, "http_requests_total{method=\"GET\",path=\"/health\",status=\"200\"}")
	assert.Contains(t, body, "http_request_duration_seconds_bucket{method=\"GET\",path=\"/health\",le=\"+")
	assert.True(t, strings.Contains(body, "http_request_duration_seconds_sum{method=\"GET\",path=\"/health\"}") || strings.Contains(body, "http_request_duration_seconds_sum "))
}
