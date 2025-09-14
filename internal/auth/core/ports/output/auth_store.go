package output

import (
	"context"

	"github.com/lechitz/AionApi/internal/auth/core/domain"
)

// AuthStore abstracts persistence for the current valid token per user.
type AuthStore interface {
	Save(ctx context.Context, token domain.Auth) error
	Get(ctx context.Context, tokenKey uint64) (domain.Auth, error)
	Delete(ctx context.Context, tokenKey uint64) error
}
