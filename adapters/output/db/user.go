package db

import (
	"github.com/jinzhu/copier"
	"github.com/lechitz/AionApi/internal/core/constants"
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
	ID        uint64         `gorm:"primaryKey, column:id"`
	Name      string         `gorm:"column:name"`
	Username  string         `gorm:"column:username"`
	Email     string         `gorm:"column:email"`
	Password  string         `gorm:"column:password"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
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
		DeletedAt: u.DeletedAt,
	}
}

func (up UserPostgresDB) CreateUser(contextControl domain.ContextControl, userDomain domain.UserDomain) (domain.UserDomain, error) {

	var userDB UserDB
	copier.Copy(&userDB, &userDomain)

	if err := up.DB.WithContext(contextControl.Context).
		Create(&userDB).Error; err != nil {
		up.LoggerSugar.Errorw(constants.ErrorToCreateUser, "error", err.Error())
		return domain.UserDomain{}, err
	}

	return userDB.CopyToUserDomain(), nil
}

func (up UserPostgresDB) GetAllUsers(contextControl domain.ContextControl) ([]domain.UserDomain, error) {
	var usersDB []UserDB
	var usersDomain []domain.UserDomain

	if err := up.DB.WithContext(contextControl.Context).
		Where("deleted_at IS NULL").
		Select("id", "name", "username", "email", "created_at").
		Find(&usersDB).Error; err != nil {
		up.LoggerSugar.Errorw(constants.ErrorToGetAllUsers, "error", err.Error())
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
		up.LoggerSugar.Errorw(constants.ErrorToGetUserByID, "error", err.Error())
		return domain.UserDomain{}, err
	}

	return userDB.CopyToUserDomain(), nil
}

func (up UserPostgresDB) GetUserByUsername(contextControl domain.ContextControl, userDomain domain.UserDomain) (domain.UserDomain, error) {

	var userDB UserDB
	copier.Copy(&userDB, &userDomain)

	if err := up.DB.WithContext(contextControl.Context).
		Select("id, username, password").
		Where("username = ?", userDB.Username).
		First(&userDB).Error; err != nil {
		up.LoggerSugar.Errorw(constants.ErrorToGetUserByUserName, "error", err.Error())
		return domain.UserDomain{}, err
	}

	return userDB.CopyToUserDomain(), nil
}

func (up UserPostgresDB) UpdateUser(contextControl domain.ContextControl, userDomain domain.UserDomain) (domain.UserDomain, error) {

	var userDB UserDB
	copier.Copy(&userDB, &userDomain)

	if err := up.DB.WithContext(contextControl.Context).
		Model(&UserDB{}).
		Where("id = ?", userDB.ID).
		Updates(userDB).Error; err != nil {
		up.LoggerSugar.Errorw(constants.ErrorToUpdateUser, "error", err.Error())
		return domain.UserDomain{}, err
	}

	return userDB.CopyToUserDomain(), nil
}

func (up UserPostgresDB) UpdatePassword(contextControl domain.ContextControl, userDomain domain.UserDomain) (domain.UserDomain, error) {

	var userDB UserDB
	copier.Copy(&userDB, &userDomain)

	if err := up.DB.WithContext(contextControl.Context).
		Model(&UserDB{}).
		Where("id = ?", userDB.ID).
		Update("password", userDB.Password).Error; err != nil {
		up.LoggerSugar.Errorw(constants.ErrorToUpdatePassword, "error", err.Error())
		return domain.UserDomain{}, err
	}

	return userDB.CopyToUserDomain(), nil
}

func (up UserPostgresDB) SoftDeleteUser(contextControl domain.ContextControl, userID uint64) error {

	if err := up.DB.WithContext(contextControl.Context).
		Model(&UserDB{}).
		Where("id = ?", userID).
		Update("deleted_at", time.Now()).Error; err != nil {
		up.LoggerSugar.Errorw(constants.ErrorToSoftDeleteUser, "error", err.Error())
		return err
	}

	return nil
}
