// Package domain contains core business entities used throughout the application.
package domain

import "time"

// UserDomain represents a user within the AionApi system.
// It contains identification data, credentials, and lifecycle metadata.
type UserDomain struct {
	CreatedAt time.Time  // Timestamp of when the user was created
	UpdatedAt time.Time  // Timestamp of the last update
	DeletedAt *time.Time // Soft delete marker (nil if active)
	Name      string     // Full name of the user
	Username  string     // Username used for login
	Email     string     // Email address
	Password  string     // Hashed password
	ID        uint64     // Unique identifier for the user
}
