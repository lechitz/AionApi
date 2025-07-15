package model

import (
	"github.com/lechitz/AionApi/internal/adapters/secondary/db/postgres/constants"
	"time"

	"gorm.io/gorm"
)

// TODO: analisar o impacto ao retirar a lib time do domain.

// UserDB represents the database model for storing user information.
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
