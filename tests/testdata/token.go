// Package testdata contains test data used for testing purposes.
package testdata

import (
	"github.com/lechitz/AionApi/internal/core/domain"
)

// TestPerfectToken is a predefined TokenDomain instance used for testing purposes, representing a valid token associated with a user.
var TestPerfectToken = domain.TokenDomain{
	UserID: 1,
	Token:  "token_abc123",
	//CreatedAt: time.Now(),
	//ExpiresAt: time.Now().Add(24 * time.Hour),
}
