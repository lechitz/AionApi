package dto

import (
	"time"

	"github.com/lechitz/AionApi/internal/user/core/ports/input"
)

// UpdateUserRequest represents the payload for updating user information.
// Fields that are nil will not update the corresponding values on the user.
type UpdateUserRequest struct {
	Name     *string `json:"name,omitempty"`
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
}

func (r UpdateUserRequest) ToCommand() input.UpdateUserCommand {
	return input.UpdateUserCommand{
		Name:     r.Name,
		Username: r.Username,
		Email:    r.Email,
	}
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
