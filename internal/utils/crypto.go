package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// SHA256Hex returns a SHA-256 hash of the input string.
// The hash is returned as a 64-character hexadecimal string (0-9, a-f).
func SHA256Hex(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}
