// Package input is the input port for the admin context.
package input

import (
	"context"

	"github.com/lechitz/AionApi/internal/admin/core/domain"
)

// RoleUpdater defines methods for updating user roles in the system.
type RoleUpdater interface {
	UpdateUserRoles(ctx context.Context, cmd UpdateUserRolesCommand) (domain.AdminUser, error)
}

// RoleManager defines methods for managing user roles with hierarchy validation.
type RoleManager interface {
	PromoteToAdmin(ctx context.Context, cmd PromoteToAdminCommand) (domain.AdminUser, error)
	DemoteFromAdmin(ctx context.Context, cmd DemoteFromAdminCommand) (domain.AdminUser, error)
	BlockUser(ctx context.Context, cmd BlockUserCommand) (domain.AdminUser, error)
	UnblockUser(ctx context.Context, cmd UnblockUserCommand) (domain.AdminUser, error)
}

// AdminService aggregates admin operations.
type AdminService interface {
	RoleUpdater // Legacy: generic role update (will be deprecated)
	RoleManager // New: specific role operations with hierarchy validation
}
