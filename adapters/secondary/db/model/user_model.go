package model

import (
	"time"

	"github.com/lechitz/AionApi/adapters/secondary/db/constants"
	"gorm.io/gorm"
)

// UserDB represents the database model for storing user information.
// It includes fields for user details, timestamps, and soft deletion.
type UserDB struct {
	Name      string         `gorm:"column:name"`
	Username  string         `gorm:"column:username"`
	Email     string         `gorm:"column:email"`
	Password  string         `gorm:"column:password"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
	ID        uint64         `gorm:"primaryKey;column:user_id"`
}

// TableName specifies the custom database table name for the UserDB model.
func (UserDB) TableName() string {
	return constants.TableUsers
}
