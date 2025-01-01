package repository

import (
	"github.com/jinzhu/copier"
	"github.com/lechitz/AionApi/internal/core/domain"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

const (
	ErrorToCreateUser  = "error to create user into postgres"
	ErrorToGetAllUser  = "error to get all users from postgres"
	ErrorToGetUserByID = "error to get user by ID from postgres"
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
		UpdatedAt: u.UpdatedAt,
	}
}

func (up UserPostgresDB) CreateUser(contextControl domain.ContextControl, userDomain domain.UserDomain) (domain.UserDomain, error) {

	var userDB UserDB
	copier.Copy(&userDB, &userDomain)

	if err := up.DB.WithContext(contextControl.Context).
		Create(&userDB).Error; err != nil {
		up.LoggerSugar.Errorw(ErrorToCreateUser, "error", err.Error())
		return domain.UserDomain{}, err
	}

	return userDB.CopyToUserDomain(), nil
}

func (up UserPostgresDB) GetAllUsers(contextControl domain.ContextControl) ([]domain.UserDomain, error) {
	var usersDB []UserDB
	var usersDomain []domain.UserDomain

	if err := up.DB.WithContext(contextControl.Context).
		Select("id", "name", "username", "email").
		Find(&usersDB).Error; err != nil {
		up.LoggerSugar.Errorw(ErrorToGetAllUser, "error", err.Error())
		return []domain.UserDomain{}, err
	}

	for _, userDB := range usersDB {
		usersDomain = append(usersDomain, userDB.CopyToUserDomain())
	}

	return usersDomain, nil
}

func (up UserPostgresDB) GetUserByID(contextControl domain.ContextControl, userID uint64) (domain.UserDomain, error) {
	var userDB UserDB

	if err := up.DB.WithContext(contextControl.Context).
		Where("id = ?", userID).
		First(&userDB).Error; err != nil {
		up.LoggerSugar.Errorw(ErrorToGetUserByID, "error", err.Error())
		return domain.UserDomain{}, err
	}

	return userDB.CopyToUserDomain(), nil
}
