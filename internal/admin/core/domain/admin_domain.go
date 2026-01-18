// Package domain is the domain for the admin context.
package domain

import "time"

// AdminUser represents a user in the admin context.
// This is a separate domain model from user.User to maintain bounded context separation.
// It contains only the fields relevant for administrative operations.
type AdminUser struct {
	ID        uint64    // User identifier
	Username  string    // Username for display/reference
	Email     string    // Email for contact
	Name      string    // Full name
	Roles     []string  // Current roles assigned to the user
	IsActive  bool      // Whether the user is active (not deleted)
	CreatedAt time.Time // When the user was created
	UpdatedAt time.Time // When the user was last updated
}

// HasRole checks if the admin user has a specific role.
func (u AdminUser) HasRole(role string) bool {
	return HasRole(u.Roles, role)
}

// GetHighestRole returns the highest privilege role of the admin user.
func (u AdminUser) GetHighestRole() string {
	return GetHighestRole(u.Roles)
}

// IsBlocked checks if the user is blocked.
func (u AdminUser) IsBlocked() bool {
	return u.HasRole(RoleBlocked)
}

// IsAdmin checks if the user is an admin.
func (u AdminUser) IsAdmin() bool {
	return u.HasRole(RoleAdmin)
}

// IsOwner checks if the user is an owner.
func (u AdminUser) IsOwner() bool {
	return u.HasRole(RoleOwner)
}

// Role represents a role in the system.
type Role struct {
	ID          uint64
	Name        string
	Description string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// UserRole represents the assignment of a role to a user.
type UserRole struct {
	ID         uint64
	UserID     uint64
	RoleID     uint64
	AssignedBy *uint64
	AssignedAt time.Time
}

// RoleUpdate represents the data for updating a user's roles.
type RoleUpdate struct {
	Roles      []string // New roles to assign to the user
	UserID     uint64   // ID of the user to update
	AssignedBy uint64   // ID of the admin making the change
}

// ValidRoles returns the list of valid roles in the system.
func ValidRoles() []string {
	return []string{RoleOwner, RoleAdmin, RoleUser, RoleBlocked}
}

// Role constants with hierarchy: Owner > Admin > User > Blocked.
const (
	RoleOwner   = "owner"   // System owner - highest privilege
	RoleAdmin   = "admin"   // Administrator
	RoleUser    = "user"    // Default user
	RoleBlocked = "blocked" // Blocked user
)

// IsValidRole checks if a role is valid.
func IsValidRole(role string) bool {
	for _, r := range ValidRoles() {
		if r == role {
			return true
		}
	}
	return false
}
