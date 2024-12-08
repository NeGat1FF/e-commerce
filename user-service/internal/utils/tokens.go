package utils

import (
	"crypto/sha256"
	"fmt"
)

// Hash token to store in database.
func HashToken(token string) string {
	// Hash token using SHA-256
	hash := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", hash)
}
