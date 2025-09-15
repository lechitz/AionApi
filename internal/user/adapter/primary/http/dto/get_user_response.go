// Package dto provides Data Transfer Objects for the user HTTP layer.
package dto

import "time"

// GetUserResponse represents the response structure returned when fetching user details.
type GetUserResponse struct {
	// CreatedAt is the timestamp when the user was created.
	// Format: date-time. Example: "2024-01-02T15:04:05Z".
	CreatedAt time.Time `json:"created_at" example:"2024-01-02T15:04:05Z"`

	// Name is the user's display name.
	// Example: "Felipe Lechitz".
	Name string `json:"name" example:"Felipe Lechitz"`

	// Username is the user's unique handle.
	// Example: "lechitz".
	Username string `json:"username" example:"lechitz"`

	// Email is the user's email address.
	// Example: "dev@aion.local".
	Email string `json:"email" example:"dev@aion.local"`

	// ID is the user's identifier.
	// Example: 42.
	ID uint64 `json:"user_id" example:"42"`
}
