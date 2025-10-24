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
}
