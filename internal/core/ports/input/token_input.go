package input

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain"
)

// TokenCreator defines a method for creating a new token in the system.
type TokenCreator interface {
	CreateToken(ctx context.Context, token domain.TokenDomain) (string, error)
}

// TokenVerifier defines a method for verifying the validity of a given token.
type TokenVerifier interface {
	GetToken(ctx context.Context, token string) (uint64, string, error)
}

// TokenRemover defines a method for deleting a token from the system.
type TokenRemover interface {
	Delete(ctx context.Context, token domain.TokenDomain) error
}

// TokenService combines the functionalities of token creation, verification, deletion, and update operations.
type TokenService interface {
	TokenCreator
	TokenVerifier
	TokenRemover
}
