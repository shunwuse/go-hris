package request

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type decodeJSONPayload struct {
	Name string `json:"name"`
}

func TestDecodeJSON_DecodesWhenContentLengthUnknown(t *testing.T) {
	req := httptest.NewRequest("POST", "/test", strings.NewReader(`{"name":"alice"}`))
	req.ContentLength = -1 // Simulate unknown request body length (for example, chunked transport).

	var payload decodeJSONPayload
	err := DecodeJSON(req, &payload)
	if err != nil {
		t.Fatalf("DecodeJSON returned error: %v", err)
	}

	if payload.Name != "alice" {
		t.Fatalf("expected payload.Name to be %q, got %q", "alice", payload.Name)
	}
}

func TestDecodeJSON_EmptyBodyReturnsNil(t *testing.T) {
	req := httptest.NewRequest("POST", "/test", http.NoBody)
	req.ContentLength = 0

	var payload decodeJSONPayload
	err := DecodeJSON(req, &payload)
	if err != nil {
		t.Fatalf("DecodeJSON returned error for empty body: %v", err)
	}
}
