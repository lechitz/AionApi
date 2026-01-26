// Package output defines output ports for auth.
package output

import (
	"context"
	"time"
)

// RoleCache stores role snapshots for users to avoid repeated DB lookups.
type RoleCache interface {
	SaveRoles(ctx context.Context, userID uint64, roles []string, ttl time.Duration) error
	GetRoles(ctx context.Context, userID uint64) ([]string, error)
	InvalidateRoles(ctx context.Context, userID uint64) error
}
