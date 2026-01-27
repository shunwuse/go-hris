package utils

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// EncodeCursor encodes multiple values into a single opaque base64 string.
// Example: EncodeCursor(123) or EncodeCursor("2023-01-01", 123).
func EncodeCursor(parts ...any) string {
	strParts := make([]string, len(parts))
	for idx, part := range parts {
		strParts[idx] = fmt.Sprintf("%v", part)
	}

	combined := strings.Join(strParts, ":")

	return base64.StdEncoding.EncodeToString([]byte(combined))
}

// DecodeCursor decodes an opaque base64 string into its constituent parts.
func DecodeCursor(cursor string) ([]string, error) {
	if cursor == "" {
		return nil, nil
	}

	decoded, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(decoded), ":"), nil
}
