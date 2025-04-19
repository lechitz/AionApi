package model

import (
	"github.com/lechitz/AionApi/adapters/secondary/db/constants"
	"gorm.io/gorm"
	"time"
)

type UserDB struct {
	ID        uint64         `gorm:"primaryKey;column:id"`
	Name      string         `gorm:"column:name"`
	Username  string         `gorm:"column:username"`
	Email     string         `gorm:"column:email"`
	Password  string         `gorm:"column:password"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (UserDB) TableName() string {
	return constants.TableUsers
}
