package testdata

import (
	"github.com/lechitz/AionApi/internal/core/domain"
	"time"
)

var TestPerfectCategory = domain.Category{
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
