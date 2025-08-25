package dto

import "time"

// UpdateUserResponse represents the structure of the response data returned after a user update operation.
// It contains the updated timestamp, name, username, email, and user ID.
type UpdateUserResponse struct {
	UpdatedAt time.Time `json:"updated_at"`
	Name      *string   `json:"name"`
	Username  *string   `json:"username"`
	Email     *string   `json:"email"`
	ID        uint64    `json:"user_id"`
}
