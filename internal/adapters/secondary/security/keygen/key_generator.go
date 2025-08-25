package keygen

import (
	"crypto/rand"
	"encoding/base64"
)

// KeyLength defines the length of the key.
const KeyLength = 64

// KeyGenerator implements output.JWTKeyGenerator.
type KeyGenerator struct{}

// New creates a new KeyGenerator instance.
func New() *KeyGenerator {
	return &KeyGenerator{}
}

// Generate generates a random 64-byte JWT key as base64.
func (g *KeyGenerator) Generate() (string, error) {
	key := make([]byte, KeyLength)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}
