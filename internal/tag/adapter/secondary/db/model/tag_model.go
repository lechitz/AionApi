// Package model contains database models for the Tag context.
package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	// TagTable is the fully qualified table name for tags.
	TagTable = "aion_api.tags"
)

// TagDB represents the database row for a tag.
type TagDB struct {
	ID          uint64         `gorm:"primaryKey;column:tag_id"`
	UserID      uint64         `gorm:"column:user_id;not null;index"`
	CategoryID  uint64         `gorm:"column:category_id;not null;index"`
	Name        string         `gorm:"column:name;type:text;not null;index"`
	Description string         `gorm:"column:description;type:text"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// TableName implements GORM's tabler interface and returns the fully qualified
// database table name for TagDB.
func (TagDB) TableName() string { return TagTable }
