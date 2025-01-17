package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// //You need to run this script and copy the generated key to the .env file
// GenerateJWTKey generates a random key for JWT
func GenerateJWTKey() (string, error) {
	key := make([]byte, 64)

	if _, err := rand.Read(key); err != nil {
		return "", err
	}

	stringBase64 := base64.StdEncoding.EncodeToString(key)
	fmt.Println("JWT_KEY: " + stringBase64)

	return stringBase64, nil
}
