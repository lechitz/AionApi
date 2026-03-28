// Package output defines output ports for user bounded context.
package output

import (
	"context"
	"time"

	"github.com/lechitz/aion-api/internal/user/core/domain"
)

// UserCache defines operations for caching user profile data.
// SECURITY: This cache NEVER stores password hashes - those always come from database.
type UserCache interface {
	SaveUser(ctx context.Context, user domain.User, expiration time.Duration) error
	GetUserByID(ctx context.Context, userID uint64) (domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	DeleteUser(ctx context.Context, userID uint64, username, email string) error
}
