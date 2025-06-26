// Package cache defines interfaces for token-related operations.
package cache

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain/entity"
)

// Creator defines an interface for creating tokens within a given context.
type Creator interface {
	CreateToken(ctx context.Context, token entity.TokenDomain) (string, error)
}

// TokenChecker defines a method to retrieve a user identifier associated with a given token from the system.
type TokenChecker interface {
	Get(ctx context.Context, token entity.TokenDomain) (string, error)
}

// TokenSaver defines a method for persisting a token in the system with its associated user information.
type TokenSaver interface {
	Save(ctx context.Context, token entity.TokenDomain) error
}

// TokenUpdater defines a method to update an existing token's details in the system.
type TokenUpdater interface {
	Update(ctx context.Context, token entity.TokenDomain) error
}

// TokenDeleter defines a method for deleting a token from the system associated with a given user.
type TokenDeleter interface {
	Delete(ctx context.Context, token entity.TokenDomain) error
}

// TokenVerify defines an interface for verifying tokens in a specific context.
// VerifyToken checks the validity of a token and returns user ID, role, or an error.
type TokenVerify interface {
	VerifyToken(ctx context.Context, token string) (uint64, string, error)
}

// TokenRepositoryPort is an interface that combines token-related operations such as checking, saving, updating, and deleting tokens.
type TokenRepositoryPort interface {
	TokenChecker
	TokenSaver
	TokenUpdater
	TokenDeleter
}
