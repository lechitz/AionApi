package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateJWTKey() (string, error) {
	key := make([]byte, 64)

	if _, err := rand.Read(key); err != nil {
		return "", err
	}

	stringBase64 := base64.StdEncoding.EncodeToString(key)
	fmt.Println("JWT_KEY: " + stringBase64)

	return stringBase64, nil
}
