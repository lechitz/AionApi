// Package output defines output ports for the admin context.
package output

import "context"

// SessionRevoker invalidates active user sessions (access/refresh tokens).
type SessionRevoker interface {
	RevokeUserSessions(ctx context.Context, userID uint64) error
}
