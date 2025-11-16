// Package model contains database models for the record bounded context.
package model

import (
	"time"
)

// Record represents the database model for a record in aion_api.records table.
type Record struct {
	ID           uint64     `gorm:"column:id;primaryKey;autoIncrement"`
	UserID       uint64     `gorm:"column:user_id;not null;index:idx_records_user_event,priority:1"`
	Title        string     `gorm:"column:title;type:varchar(255);not null"`
	Description  *string    `gorm:"column:description;type:text"`
	CategoryID   uint64     `gorm:"column:category_id;not null;index:idx_records_category_user"`
	TagID        uint64     `gorm:"column:tag_id;not null;index:idx_records_tag_user"`
	EventTime    time.Time  `gorm:"column:event_time;not null;index:idx_records_user_event,priority:2"`
	RecordedAt   *time.Time `gorm:"column:recorded_at"`
	DurationSecs *int       `gorm:"column:duration_seconds"`
	Value        *float64   `gorm:"column:value"`
	Source       *string    `gorm:"column:source;type:varchar(100)"`
	Timezone     *string    `gorm:"column:timezone;type:varchar(100)"`
	Status       *string    `gorm:"column:status;type:varchar(50)"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt    *time.Time `gorm:"column:deleted_at;index"`
}

// TableName specifies the table name for GORM.
func (Record) TableName() string {
	return "aion_api.records"
}
