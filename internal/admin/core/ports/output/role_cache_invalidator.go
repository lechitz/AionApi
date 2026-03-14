// Package output defines output ports for the admin context.
package output

import "context"

// RoleCacheInvalidator removes cached role snapshots for a user.
type RoleCacheInvalidator interface {
	InvalidateRoles(ctx context.Context, userID uint64) error
}
