package request

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type decodeJSONPayload struct {
	Name string `json:"name"`
}

type decodePayload struct {
	ID   int    `json:"id" schema:"id"`
	Name string `json:"name" schema:"name"`
}

func TestDecodeJSON_DecodesWhenContentLengthUnknown(t *testing.T) {
	req := httptest.NewRequest("POST", "/test", strings.NewReader(`{"name":"alice"}`))
	req.ContentLength = -1 // Simulate unknown request body length (for example, chunked transport).

	var payload decodeJSONPayload
	err := DecodeJSON(req, &payload)
	require.NoError(t, err)
	assert.Equal(t, "alice", payload.Name)
}

func TestDecodeJSON_EmptyBodyReturnsNil(t *testing.T) {
	req := httptest.NewRequest("POST", "/test", http.NoBody)
	req.ContentLength = 0

	var payload decodeJSONPayload
	err := DecodeJSON(req, &payload)
	require.NoError(t, err)
}

func TestDecodeJSON_InvalidJSONReturnsError(t *testing.T) {
	req := httptest.NewRequest("POST", "/test", strings.NewReader(`{"name":`))

	var payload decodeJSONPayload
	err := DecodeJSON(req, &payload)
	require.Error(t, err)
	assert.False(t, errors.Is(err, http.ErrBodyNotAllowed))
}

func TestDecodeQuery(t *testing.T) {
	req := httptest.NewRequest("GET", "/users?id=12&name=alice", nil)

	var payload decodePayload
	err := DecodeQuery(req, &payload)
	require.NoError(t, err)
	assert.Equal(t, 12, payload.ID)
	assert.Equal(t, "alice", payload.Name)
}

func TestDecodePath(t *testing.T) {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "99")

	req := httptest.NewRequest("GET", "/users/99", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	var payload decodePayload
	err := DecodePath(req, &payload)
	require.NoError(t, err)
	assert.Equal(t, 99, payload.ID)
}

func TestDecode_PathQueryBodyOrder(t *testing.T) {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	rctx.URLParams.Add("name", "from-path")

	req := httptest.NewRequest("POST", "/users/1?name=from-query", strings.NewReader(`{"name":"from-body"}`))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	var payload decodePayload
	err := Decode(req, &payload)
	require.NoError(t, err)

	assert.Equal(t, 1, payload.ID)
	assert.Equal(t, "from-body", payload.Name)
}

func TestDecodePath_WithoutRouteContext(t *testing.T) {
	req := httptest.NewRequest("GET", "/users/1", nil)

	var payload decodePayload
	err := DecodePath(req, &payload)
	require.NoError(t, err)
}
