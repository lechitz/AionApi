package middlewares

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
)

// Script to generate a random key for JWT. You need to run this script and copy the generated key to the .env file
func GenerateJWTKey() {
	key := make([]byte, 64)

	if _, err := rand.Read(key); err != nil {
		log.Fatal(err)
	}

	stringBase64 := base64.StdEncoding.EncodeToString(key)
	fmt.Println("JWT_KEY: " + stringBase64)
}
