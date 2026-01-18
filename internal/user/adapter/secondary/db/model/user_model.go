// Package model contains the database model for users.
package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	// TableUsers is the name of the database table for users.
	TableUsers = "aion_api.users"
)

// UserDB represents the database model for storing user information.
type UserDB struct {
	Name      string         `gorm:"column:name"`
	Username  string         `gorm:"column:username"`
	Email     string         `gorm:"column:email"`
	Password  string         `gorm:"column:password"`
	Locale    *string        `gorm:"column:locale"`
	Timezone  *string        `gorm:"column:timezone"`
	Location  *string        `gorm:"column:location"`
	Bio       *string        `gorm:"column:bio"`
	AvatarURL *string        `gorm:"column:avatar_url"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
	ID        uint64         `gorm:"primaryKey;column:user_id"`
}

// TableName specifies the custom database table name for the UserDB model.
func (UserDB) TableName() string {
	return TableUsers
}
