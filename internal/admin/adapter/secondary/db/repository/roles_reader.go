// Package repository provides admin database operations.
package repository

import (
	"context"
)

// GetRolesByUserID returns the role names for a user.
//
// This method is consumed by the auth bounded context via an output port
// (auth/core/ports/output.RolesReader). Admin remains the source of truth for roles.
func (r *AdminRepository) GetRolesByUserID(ctx context.Context, userID uint64) ([]string, error) {
	return r.getUserRoles(ctx, userID)
}
