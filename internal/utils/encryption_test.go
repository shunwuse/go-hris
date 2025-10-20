package utils

import (
	"testing"
)

func TestEncrypt(t *testing.T) {
	// create a password
	password := "password"

	// hash the password
	hashedPassword, err := HashPassword(password)

	t.Logf("Hasded password: %s", hashedPassword)

	// check if there is an error
	if err != nil {
		t.Errorf("Error hashing password: %v", err)
	}

	// check if the hashed password is empty
	if hashedPassword == "" {
		t.Errorf("Expected hashed password to be non-empty")
	}

	// check if the password matches the hashed password
	if !CheckPasswordHash(password, hashedPassword) {
		t.Errorf("Expected password to match hashed password")
	}

	// check if the password does not match the hashed password
	if CheckPasswordHash("wrongpassword", hashedPassword) {
		t.Errorf("Expected password to not match hashed password")
	}
}
