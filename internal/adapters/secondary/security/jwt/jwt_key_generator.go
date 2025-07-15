package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/lechitz/AionApi/internal/adapters/secondary/security/jwt/constants"
)

// KeyGenerator implements output.JWTKeyGenerator.
type KeyGenerator struct{}

// NewKeyGenerator creates a new KeyGenerator instance.
func NewKeyGenerator() *KeyGenerator {
	return &KeyGenerator{}
}

// GenerateKey generates a random 64-byte JWT key as base64.
func (g *KeyGenerator) GenerateKey() (string, error) {
	key := make([]byte, constants.JWTKeyLength)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}
