package security

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateJWTKey generates a random 64-byte JWT key and returns it as a base64-encoded string or an error.
func GenerateJWTKey() (string, error) {
	key := make([]byte, 64)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}

	stringBase64 := base64.StdEncoding.EncodeToString(key)
	return stringBase64, nil
}
