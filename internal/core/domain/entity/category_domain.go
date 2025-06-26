// Package entity contains core business entities used throughout the application.
package entity

import "time"

// Category represents a user-defined classification for organizing habits or tasks within the Aion system.
// It contains metadata such as color, icon, and optional description.
type Category struct {
	CreatedAt   time.Time  // Timestamp of creation
	UpdatedAt   time.Time  // Timestamp of the last update
	DeletedAt   *time.Time // Soft-delete marker (null if active)
	Name        string     // Name of the category (e.g., "Health", "Learning")
	Description string     // Optional description providing context
	Color       string     // Color code (hex or name) used for UI representation
	Icon        string     // Icon name or identifier for category visualization
	ID          uint64     // Unique identifier for the category
	UserID      uint64     // ID of the user who owns this category
}
