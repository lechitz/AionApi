// Package dto contains Data Transfer Objects used by the HTTP layer.//TODO: ajustar magic string.
package dto

import (
	"time"
)

// TODO: Analisar o impacto ao retirar a lib time do domain, aqui vai precisar ser ajustado.

// CreateUserRequest represents the data required to create a new user in the system.
type CreateUserRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateUserResponse represents the structure of the response returned after a user is successfully created.
type CreateUserResponse struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	ID       uint64 `json:"user_id"`
}

// GetUserResponse represents the response structure returned when fetching user details.
type GetUserResponse struct {
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	ID        uint64    `json:"user_id"`
}

// UpdateUserRequest represents the payload for updating user information.
// Fields that are nil will not update the corresponding values on the user.
// ID is required to identify the user to be updated.
type UpdateUserRequest struct {
	Name     *string `json:"name,omitempty"`
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
	ID       uint64  `json:"user_id"`
}

// UpdateUserResponse represents the structure of the response data returned after a user update operation.
// It contains the updated timestamp, name, username, email, and user ID.
type UpdateUserResponse struct {
	UpdatedAt time.Time `json:"updated_at"`
	Name      *string   `json:"name"`
	Username  *string   `json:"username"`
	Email     *string   `json:"email"`
	ID        uint64    `json:"user_id"`
}

// UpdatePasswordUserRequest represents a request to update a user's password.
// It includes the current password and a new password for the user.
type UpdatePasswordUserRequest struct {
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
}
