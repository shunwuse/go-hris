package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	apperrors "github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/http/middlewares"
	"github.com/shunwuse/go-hris/internal/http/response"
	"github.com/shunwuse/go-hris/internal/http/routes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testRoute struct{}

func (testRoute) Setup(router chi.Router) {
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}

func newTestRouter() *Router {
	return New(middlewares.CommonMiddlewares{}, routes.Routes{testRoute{}})
}

func TestServeHTTP_NotFoundReturnsJSONEnvelope(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/missing", nil)

	newTestRouter().ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "application/json; charset=utf-8", rr.Header().Get("Content-Type"))

	var body response.ErrorResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	assert.Equal(t, apperrors.CodeNotFound, body.Error.Code)
	assert.Equal(t, "resource not found", body.Error.Message)
}
