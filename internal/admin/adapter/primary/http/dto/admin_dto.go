// Package dto is the data transfer object for the admin context.
package dto

import (
	"errors"
	"time"

	"github.com/lechitz/aion-api/internal/admin/core/ports/input"
)

// UpdateUserRolesRequest represents the payload for updating a user's roles.
type UpdateUserRolesRequest struct {
	// Roles is the list of roles to assign to the user.
	// Valid roles: "user", "admin", "blocked"
	// Example: ["user", "admin"]
	Roles []string `json:"roles" example:"user,admin"`
}

// Validate ensures the request contains valid data.
func (r UpdateUserRolesRequest) Validate() error {
	if len(r.Roles) == 0 {
		return ErrNoRolesProvided
	}
	for _, role := range r.Roles {
		if role == "" {
			return ErrEmptyRoleProvided
		}
	}
	return nil
}

// ToCommand converts the request to a domain command.
func (r UpdateUserRolesRequest) ToCommand(userID uint64) input.UpdateUserRolesCommand {
	return input.UpdateUserRolesCommand{
		UserID: userID,
		Roles:  r.Roles,
	}
}

// UpdateUserRolesResponse represents the response after updating a user's roles.
type UpdateUserRolesResponse struct {
	// ID is the user's unique identifier.
	// Example: 42
	ID uint64 `json:"user_id" example:"42"`

	// Username is the user's username.
	// Example: "alice"
	Username string `json:"username" example:"alice"`

	// Email is the user's email.
	// Example: "alice@example.com"
	Email string `json:"email" example:"alice@example.com"`

	// Roles is the updated list of roles.
	// Example: ["user", "admin"]
	Roles []string `json:"roles" example:"user,admin"`

	// UpdatedAt is the timestamp when the user was updated.
	// Example: "2025-09-14T22:01:02Z"
	UpdatedAt time.Time `json:"updated_at" format:"date-time" example:"2025-09-14T22:01:02Z"`
}

var (
	// ErrNoRolesProvided is returned when no roles are provided in the request.
	ErrNoRolesProvided = errors.New("no roles provided")

	// ErrEmptyRoleProvided is returned when an empty role string is provided.
	ErrEmptyRoleProvided = errors.New("empty role provided")
)
