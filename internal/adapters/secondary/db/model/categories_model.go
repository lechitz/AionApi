// Package model The package model contains database models for the application.
package model

import (
	"time"

	"github.com/lechitz/AionApi/internal/adapters/secondary/db/constants"

	"gorm.io/gorm"
)

// CategoryDB represents the database model for a category entity with metadata and user association.
type CategoryDB struct {
	CreatedAt   time.Time      `gorm:"column:created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at"`
	Name        string         `gorm:"column:name"`
	Description string         `gorm:"column:description"`
	Color       string         `gorm:"column:color_hex"`
	Icon        string         `gorm:"column:icon"`
	ID          uint64         `gorm:"primaryKey;column:category_id"`
	UserID      uint64         `gorm:"column:user_id"`
}

// TableName specifies the database table name for the CategoryDB model.
func (CategoryDB) TableName() string {
	return constants.CategoryTable
}
