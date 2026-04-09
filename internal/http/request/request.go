package request

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func init() {
	// Ignore unknown keys to prevent errors if extra query params are passed.
	decoder.IgnoreUnknownKeys(true)
}

// Decode decodes path, query, and JSON body.
// It follows the order: Path -> Query -> Body (JSON).
// Later values will override earlier ones if they share the same tag/field name.
func Decode(r *http.Request, dst any) error {
	if err := DecodePath(r, dst); err != nil {
		return err
	}

	if err := DecodeQuery(r, dst); err != nil {
		return err
	}

	if err := DecodeJSON(r, dst); err != nil {
		return err
	}

	return nil
}

// DecodeQuery decodes URL query parameters.
func DecodeQuery(r *http.Request, dst any) error {
	return decoder.Decode(dst, r.URL.Query())
}

// DecodePath decodes URL path parameters.
func DecodePath(r *http.Request, dst any) error {
	rctx := chi.RouteContext(r.Context())
	if rctx == nil {
		return nil
	}

	params := make(map[string][]string)
	for i, key := range rctx.URLParams.Keys {
		params[key] = []string{rctx.URLParams.Values[i]}
	}

	return decoder.Decode(dst, params)
}

// DecodeJSON decodes the JSON request body.
func DecodeJSON(r *http.Request, dst any) error {
	if r.Body == nil {
		return nil
	}

	// Drain and close the body to ensure connection reuse.
	defer func() {
		_, _ = io.Copy(io.Discard, r.Body)
		_ = r.Body.Close()
	}()

	err := json.NewDecoder(r.Body).Decode(dst)
	if err == io.EOF {
		return nil
	}

	return err
}
