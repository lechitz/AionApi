// Package testdata contains test data used for testing purposes.
package testdata

import (
	"github.com/lechitz/AionApi/internal/core/domain"
)

// SecretKey is the key for secrets used in token operations.
const SecretKey = "secret"

// TODO: Ajustar o uso dos testdata.

// TestPerfectToken is a predefined Value instance used for testing purposes, representing a valid token associated with a user.
var TestPerfectToken = domain.Token{
	Key:   1,
	Value: "token_abc123",
	// CreatedAt: time.Now(),
	// ExpiresAt: time.Now().Add(24 * time.Hour),
}
