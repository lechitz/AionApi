// Package output defines auth output ports.
package output

import "context"

// RolesReader abstracts how the auth context reads authorization data.
//
// Important: roles belong to the /admin bounded context. The auth context must not
// query the user domain/entity for roles.
// This port allows the auth use cases to fetch roles without importing admin packages.
type RolesReader interface {
	// GetRolesByUserID returns the role names for the given user.
	// Implementations should return at least the default role (e.g., ["user"]) when applicable.
	GetRolesByUserID(ctx context.Context, userID uint64) ([]string, error)
}
