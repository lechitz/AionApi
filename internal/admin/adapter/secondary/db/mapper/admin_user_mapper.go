// Package mapper converts between user domain and admin domain models.
package mapper

import (
	"github.com/lechitz/aion-api/internal/admin/core/domain"
	userdomain "github.com/lechitz/aion-api/internal/user/core/domain"
)

// AdminUserFromUser converts a user.User to admin.AdminUser.
// This is an anti-corruption layer that maintains bounded context separation.
func AdminUserFromUser(user userdomain.User) domain.AdminUser {
	return domain.AdminUser{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Name:      user.Name,
		IsActive:  user.DeletedAt == nil,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// UserFromAdminUser converts an admin.AdminUser back to user.User.
// Note: This only includes fields available in AdminUser.
// Other user fields (password, locale, etc) are not populated.
func UserFromAdminUser(adminUser domain.AdminUser) userdomain.User {
	return userdomain.User{
		ID:        adminUser.ID,
		Username:  adminUser.Username,
		Email:     adminUser.Email,
		Name:      adminUser.Name,
		CreatedAt: adminUser.CreatedAt,
		UpdatedAt: adminUser.UpdatedAt,
		// DeletedAt, Password, and other fields are NOT set
	}
}
