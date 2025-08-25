package dto

import (
	"github.com/lechitz/AionApi/internal/core/ports/input"
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
