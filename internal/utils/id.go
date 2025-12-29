package utils

import (
	"crypto/rand"
	"io"
	"sync"

	"github.com/oklog/ulid/v2"
)

var (
	entropyPool = &sync.Pool{
		New: func() any {
			return ulid.Monotonic(rand.Reader, 0)
		},
	}
)

// NewTraceID generates a new unique trace ID using ULID.
// It uses a sync.Pool of entropy sources for better performance and thread safety.
func NewTraceID() string {
	entropy := entropyPool.Get().(io.Reader)
	defer entropyPool.Put(entropy)

	return ulid.MustNew(ulid.Now(), entropy).String()
}
