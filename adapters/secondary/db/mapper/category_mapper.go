// Package mapper provides utility functions for mapping between domain and database objects.
package mapper

import (
	"time"

	"github.com/lechitz/AionApi/adapters/secondary/db/model"
	"github.com/lechitz/AionApi/internal/core/domain"
	"gorm.io/gorm"
)

// CategoryFromDB maps a model.CategoryDB object to a domain.Category object for further use in the application.
func CategoryFromDB(category model.CategoryDB) domain.Category {
	var deletedAt *time.Time
	if category.DeletedAt.Valid {
		deletedAt = &category.DeletedAt.Time
	}

	return domain.Category{
		ID:          category.ID,
		UserID:      category.UserID,
		Name:        category.Name,
		Description: category.Description,
		Color:       category.Color,
		Icon:        category.Icon,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
		DeletedAt:   deletedAt,
	}
}

// CategoryToDB maps a domain.Category object to a model.CategoryDB object for database operations.
func CategoryToDB(category domain.Category) model.CategoryDB {
	var deleted gorm.DeletedAt
	if category.DeletedAt != nil {
		deleted.Time = *category.DeletedAt
		deleted.Valid = true
	}

	return model.CategoryDB{
		ID:          category.ID,
		UserID:      category.UserID,
		Name:        category.Name,
		Description: category.Description,
		Color:       category.Color,
		Icon:        category.Icon,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
		DeletedAt:   deleted,
	}
}
