// Package input defines commands for admin operations.
package input

// UpdateUserRolesCommand defines the command for updating a user's roles.
type UpdateUserRolesCommand struct {
	Roles  []string // New roles to assign
	UserID uint64   // User ID to update
}

// HasUpdates returns true if the command contains valid role updates.
func (c UpdateUserRolesCommand) HasUpdates() bool {
	return len(c.Roles) > 0
}
