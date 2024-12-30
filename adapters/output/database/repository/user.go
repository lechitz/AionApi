package repository

import (
	"github.com/jinzhu/copier"
	"github.com/lechitz/AionApi/internal/core/domain"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type UserPostgresDB struct {
	DB          *gorm.DB
	LoggerSugar *zap.SugaredLogger
}

type UserDB struct {
	ID        uint64    `gorm:"primaryKey, column:id"`
	Name      string    `gorm:"column:name"`
	Username  string    `gorm:"column:username"`
	Email     string    `gorm:"column:email"`
	Password  string    `gorm:"column:password"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func NewUserPostgresDB(gormDB *gorm.DB, loggerSugar *zap.SugaredLogger) UserPostgresDB {
	return UserPostgresDB{
		DB:          gormDB,
		LoggerSugar: loggerSugar,
	}
}

func (UserDB) TableName() string {
	return "aion_api.users"
}

func (u UserDB) CopyToUserDomain() domain.UserDomain {
	return domain.UserDomain{
		ID:        u.ID,
		Name:      u.Name,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
	}
}

func (up UserPostgresDB) CreateUser(contextControl domain.ContextControl, userDomain domain.UserDomain) (domain.UserDomain, error) {

	var userDB UserDB
	copier.Copy(&userDB, &userDomain)

	if err := up.DB.WithContext(contextControl.Context).
		Create(&userDB).Error; err != nil {
		up.LoggerSugar.Errorw("error to save the user into postgres", "error", err.Error())
		return domain.UserDomain{}, err
	}

	return userDB.CopyToUserDomain(), nil
}
