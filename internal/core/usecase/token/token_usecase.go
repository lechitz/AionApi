package token

import (
	"context"
	"github.com/lechitz/AionApi/internal/core/domain"
)

// Usecase defines operations for handling token creation, verification, deletion, and persistence.
// CreateToken generates a new token for the given domain.TokenDomain.
// VerifyToken verifies a given token and extracts user ID and associated data.
// Delete removes an existing token from the system.
// Update modifies an existing token's details in the system.
// Save persists a new token into the system.
type Usecase interface {
	CreateToken(ctx context.Context, token domain.TokenDomain) (string, error)
	VerifyToken(ctx context.Context, token string) (uint64, string, error)
	Delete(ctx context.Context, token domain.TokenDomain) error
	Update(ctx context.Context, token domain.TokenDomain) error
	Save(ctx context.Context, token domain.TokenDomain) error
}
