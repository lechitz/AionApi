package testdata

import (
	"time"

	"github.com/lechitz/AionApi/internal/core/domain"
)

// PerfectCategory is a predefined Category instance representing a "Work" category with user ID 3, default blue color, and an optional description for testing purposes.
var PerfectCategory = domain.Category{
	ID:          1,
	UserID:      3,
	Name:        "Work",
	Description: "my work description",
	Color:       "blue",
	Icon:        "work",
	CreatedAt:   time.Time{},
	UpdatedAt:   time.Time{},
	DeletedAt:   nil,
}
