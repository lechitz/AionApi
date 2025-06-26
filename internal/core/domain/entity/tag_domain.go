// Package entity contains core business entities used throughout the application.
package entity

import "time"

// Tag represents a user-defined label associated with a specific category.
// It is used to organize and classify tasks or habits in a personalized way.
type Tag struct {
	CreatedAt   time.Time  // Timestamp of creation
	UpdatedAt   time.Time  // Timestamp of the last update
	DeletedAt   *time.Time // Soft-delete marker (nil if active)
	Description *string    // Optional description for additional context
	Name        string     // Name of the tag (e.g., "Morning", "High Priority")
	ID          uint64     // Unique identifier for the tag
	UserID      uint64     // ID of the user who owns the tag
	CategoryID  uint64     // ID of the category this tag belongs to
}
