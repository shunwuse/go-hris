package pagination

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// Encode encodes a value into an opaque base64-encoded JSON string.
func Encode[T any](data T) (string, error) {
	raw, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal cursor: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(raw), nil
}

// Decode decodes an opaque base64-encoded JSON string into a value of type T.
// Returns a pointer to the decoded value. If the cursor is empty, returns nil error but nil value.
func Decode[T any](cursor string) (*T, error) {
	if cursor == "" {
		return nil, nil
	}

	decoded, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return nil, fmt.Errorf("failed to decode cursor: %w", err)
	}

	var payload T
	if err := json.Unmarshal(decoded, &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cursor: %w", err)
	}

	return &payload, nil
}
