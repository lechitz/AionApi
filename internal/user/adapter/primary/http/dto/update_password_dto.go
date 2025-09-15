// Package dto (user) contains data transfer objects.
package dto

// UpdatePasswordUserRequest represents a request to update a user's password.
// Both fields are required. Minimum length recommendation: 8 characters.
type UpdatePasswordUserRequest struct {
	// Password is the current password for the user.
	// Example: "P@ssw0rd123"
	Password string `json:"password" example:"P@ssw0rd123"`

	// NewPassword is the new password to be set.
	// Example: "N3wP@ssw0rd456"
	NewPassword string `json:"new_password" example:"N3wP@ssw0rd456"`
}
