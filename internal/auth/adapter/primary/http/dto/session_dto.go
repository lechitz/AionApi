// Package dto contains Data Transfer Objects used by the HTTP layer.
package dto

import "time"

// SessionResponse represents the current authenticated session snapshot.
// This endpoint exists for UIs (e.g. dashboards) that need to render navigation/permissions.
type SessionResponse struct {
	// Authenticated indicates whether the request is authenticated.
	Authenticated bool `json:"authenticated" example:"true"`

	// UserID is the unique identifier for the authenticated user.
	UserID uint64 `json:"user_id" example:"42"`

	// Username is the unique username for the authenticated user.
	Username string `json:"username" example:"joaopereira"`

	// Email is the email address of the authenticated user.
	Email string `json:"email" example:"joao@example.com"`

	// Name is a friendly display name for the authenticated user.
	Name string `json:"name" example:"João Pereira"`

	// Roles are the user's role names (authorization snapshot).
	Roles []string `json:"roles" example:"user,admin"`

	// ExpiresAt is the access token expiration time (if available).
	ExpiresAt *time.Time `json:"expires_at,omitempty" example:"2026-01-10T20:05:00Z"`
}
