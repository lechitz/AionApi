package output

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain"
)

// TokenProvider defines an interface for generating and verifying tokens.
type TokenProvider interface {
	Generate(ctx context.Context, userID uint64) (string, error)
	Verify(ctx context.Context, tokenValue string) (map[string]any, error)
}

// TokenStore abstracts persistence for the current valid token per user.
type TokenStore interface {
	Save(ctx context.Context, token domain.Token) error
	Get(ctx context.Context, tokenKey uint64) (domain.Token, error)
	Delete(ctx context.Context, tokenKey uint64) error
}
