package dto

import "time"

// GetUserResponse represents the response structure returned when fetching user details.
type GetUserResponse struct {
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	ID        uint64    `json:"user_id"`
}
