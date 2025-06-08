package model

import (
	"github.com/lechitz/AionApi/adapters/secondary/db/constants"
	"gorm.io/gorm"
	"time"
)

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

func (CategoryDB) TableName() string {
	return constants.CategoryTable
}
