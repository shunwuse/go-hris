package utils

import (
	"testing"
)

func TestEncrypt(t *testing.T) {
	password := "password"

	hashedPassword, err := HashPassword(password)

	t.Logf("hashed password: %s", hashedPassword)

	if err != nil {
		t.Errorf("failed to hash password: %v", err)
	}

	if hashedPassword == "" {
		t.Errorf("expected hashed password to be non-empty")
	}

	if !CheckPasswordHash(password, hashedPassword) {
		t.Errorf("expected password to match hashed password")
	}

	if CheckPasswordHash("wrongpassword", hashedPassword) {
		t.Errorf("expected password to not match hashed password")
	}
}
