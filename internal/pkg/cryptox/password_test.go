package cryptox

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncrypt(t *testing.T) {
	password := "password"

	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)

	t.Logf("hashed password: %s", hashedPassword)

	assert.True(t, CheckPasswordHash(password, hashedPassword))
	assert.False(t, CheckPasswordHash("wrongpassword", hashedPassword))
}
