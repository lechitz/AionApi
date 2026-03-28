// Package output is the output port for the admin context.
package output

import (
	"context"

	"github.com/lechitz/aion-api/internal/admin/core/domain"
)

// UserRoleUpdater defines an interface for updating user roles in the system.
type UserRoleUpdater interface {
	UpdateRoles(ctx context.Context, userID uint64, roles []string) (domain.AdminUser, error)
}

// AdminUserFinder defines an interface for finding users in admin context.
type AdminUserFinder interface {
	GetByID(ctx context.Context, userID uint64) (domain.AdminUser, error)
}

// RolesReader defines an interface for reading user roles.
// This is intentionally read-only and is used by other bounded contexts (e.g., /auth).
type RolesReader interface {
	GetRolesByUserID(ctx context.Context, userID uint64) ([]string, error)
}

// AdminRepository aggregates interfaces for admin operations.
type AdminRepository interface {
	UserRoleUpdater
	AdminUserFinder
	RoleAssigner // Allows other contexts (e.g., /user) to delegate default role assignment
	RolesReader  // Allows other contexts (e.g., /auth) to read roles
}
