// Package dto (user) contains data transfer objects for the HTTP layer.
package dto

import (
	"time"

	"github.com/lechitz/AionApi/internal/user/core/ports/input"
)

// UpdateUserRequest represents the payload for updating user information.
// Any field left nil will NOT update the corresponding user attribute.
// At least one field must be provided.
type UpdateUserRequest struct {
	// Name is the new display name for the user.
	// Example: "Alice Doe"
	Name *string `json:"name,omitempty" example:"Alice Doe"`

	// Username is the new unique handle for the user.
	// Example: "alice"
	Username *string `json:"username,omitempty" example:"alice"`

	// Email is the new email address for the user.
	// Example: "alice@example.com"
	Email *string `json:"email,omitempty" example:"alice@example.com"`
}

// ToCommand converts the request to a domain command.
func (r UpdateUserRequest) ToCommand() input.UpdateUserCommand {
	return input.UpdateUserCommand{
		Name:     r.Name,
		Username: r.Username,
		Email:    r.Email,
	}
}

// UpdateUserResponse represents the response returned after a successful user update.
type UpdateUserResponse struct {
	// UpdatedAt is the timestamp when the user was updated.
	// Example: "2025-09-14T22:01:02Z"
	UpdatedAt time.Time `json:"updated_at" format:"date-time" example:"2025-09-14T22:01:02Z"`

	// Name is the current display name after the update (if changed).
	// Example: "Alice Doe"
	Name *string `json:"name" example:"Alice Doe"`

	// Username is the current username after the update (if changed).
	// Example: "alice"
	Username *string `json:"username" example:"alice"`

	// Email is the current email after the update (if changed).
	// Example: "alice@example.com"
	Email *string `json:"email" example:"alice@example.com"`

	// ID is the user's unique identifier.
	// Example: 42
	ID uint64 `json:"user_id" example:"42"`
}
