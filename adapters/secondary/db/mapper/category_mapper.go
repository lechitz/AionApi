package mapper

import (
	"time"

	"github.com/lechitz/AionApi/adapters/secondary/db/model"
	"github.com/lechitz/AionApi/internal/core/domain"
	"gorm.io/gorm"
)

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
