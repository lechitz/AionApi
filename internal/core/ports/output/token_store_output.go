package output

import (
	"context"
	"github.com/lechitz/AionApi/internal/core/domain"
)

// TokenStore defines an interface for managing tokens in the system.
type TokenStore interface {
	Save(ctx context.Context, token domain.TokenDomain) error
	Get(ctx context.Context, userID uint64) (string, error)
	Delete(ctx context.Context, token string) error
}
