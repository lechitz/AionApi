package output

import (
	"context"
	"time"

	"github.com/lechitz/AionApi/internal/auth/core/domain"
)

// AuthStore abstracts persistence for the current valid token per user (access or refresh).
// Save now accepts an expiration duration that will be used as TTL in the backing cache.
type AuthStore interface {
	Save(ctx context.Context, token domain.Auth, expiration time.Duration) error
	Get(ctx context.Context, tokenKey uint64, tokenType string) (domain.Auth, error)
	Delete(ctx context.Context, tokenKey uint64, tokenType string) error

	// SaveWithKey persists a token using a custom cache key with TTL.
	// Used for grace period tokens to allow custom key formats.
	SaveWithKey(ctx context.Context, key string, token domain.Auth, expiration time.Duration) error

	// GetByKey retrieves a token using a custom cache key.
	// Used for validating tokens during grace period.
	GetByKey(ctx context.Context, key string) (domain.Auth, error)
}
