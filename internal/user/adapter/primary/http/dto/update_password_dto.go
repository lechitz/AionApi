// Package dto (user) contains data transfer objects.
package dto

// UpdatePasswordUserRequest represents a request to update a user's password.
// It includes the current password and a new password for the user.
type UpdatePasswordUserRequest struct {
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
}
