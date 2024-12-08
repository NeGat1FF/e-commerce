package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// Generate random string with length l and expiry time in hours.
func GenerateToken(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be greater than 0")
	}

	// Create a byte slice with half the size of the desired length (2 hex chars = 1 byte)
	tokenBytes := make([]byte, length/2+length%2)

	// Fill the slice with random bytes
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	// Convert the bytes to a hex string and trim to the desired length
	token := hex.EncodeToString(tokenBytes)[:length]
	return token, nil
}
