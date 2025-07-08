// Package output token defines interfaces for token-related operations.
package output

import (
	"context"
	"net/http"

	"github.com/lechitz/AionApi/internal/core/domain"
)

// TokenSaver defines a method for persisting a token in the system with its associated user information.
type TokenSaver interface {
	Save(ctx context.Context, token domain.TokenDomain) error
}

// TokenChecker defines a method to retrieve a user identifier associated with a given token from the system.
type TokenChecker interface {
	Get(ctx context.Context, token domain.TokenDomain) (string, error)
}

// TokenDeleter defines a method for deleting a token from the system associated with a given user.
type TokenDeleter interface {
	Delete(ctx context.Context, token domain.TokenDomain) error
}

// TokenClaimsExtractor defines an interface for extracting claimskeys from a token.
type TokenClaimsExtractor interface {
	ExtractFromRequest(r *http.Request) (map[string]interface{}, error)
	ExtractFromContext(ctx context.Context) (map[string]interface{}, error)
}

// TokenStore is an interface that combines token-related operations such as checking, saving, updating, and deleting tokens.
type TokenStore interface {
	TokenSaver
	TokenChecker
	TokenDeleter
}
