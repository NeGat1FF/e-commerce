package utils

import (
	"fmt"

	// bcrypt is a library for hashing and comparing passwords.
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// ComparePasswords compares a hashed password with a plain password.
func ComparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
