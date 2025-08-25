package input

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain"
)

//TODO: preciso ajustar isso aqui ..

// TokenCreator creates a new token for a user.
type TokenCreator interface {
	Create(ctx context.Context, userID uint64) (domain.Token, error)
}

// TokenValidator validates a raw token and returns the resolved userID and claims.
type TokenValidator interface {
	Validate(ctx context.Context, token string) (userID uint64, claims map[string]any, err error)
}

// TokenRemover revokes (deletes) the token for a given user.
type TokenRemover interface {
	Revoke(ctx context.Context, userID uint64) error
}

// TokenService is the full token application port.
type TokenService interface {
	TokenCreator
	TokenValidator
	TokenRemover
}
