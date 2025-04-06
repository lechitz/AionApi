package db

import (
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/lechitz/AionApi/internal/adapters/secondary/db/constants"
	"github.com/lechitz/AionApi/internal/core/domain"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type UserRepository struct {
	db          *gorm.DB
	loggerSugar *zap.SugaredLogger
}

func NewUserRepository(db *gorm.DB, loggerSugar *zap.SugaredLogger) *UserRepository {
	return &UserRepository{
		db:          db,
		loggerSugar: loggerSugar,
	}
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

func (UserDB) TableName() string {
	return constants.TableUsers
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

func (up UserRepository) CreateUser(contextControl domain.ContextControl, userDomain domain.UserDomain) (domain.UserDomain, error) {

	var userDB UserDB
	copier.Copy(&userDB, &userDomain)

	if err := up.db.WithContext(contextControl.BaseContext).
		Create(&userDB).Error; err != nil {
		wrappedErr := fmt.Errorf(constants.ErrorToCreateUser, err)
		up.loggerSugar.Errorw(constants.ErrorToCreateUser, constants.Error, wrappedErr.Error())
		return domain.UserDomain{}, wrappedErr
	}

	return userDB.CopyToUserDomain(), nil
}

func (up UserRepository) GetAllUsers(contextControl domain.ContextControl) ([]domain.UserDomain, error) {
	var usersDB []UserDB
	var usersDomain []domain.UserDomain

	if err := up.db.WithContext(contextControl.BaseContext).
		Model(&UserDB{}).
		Select("id, name, username, email, created_at").
		Find(&usersDB).Error; err != nil {
		up.loggerSugar.Errorw(constants.ErrorToGetAllUsers, constants.Error, err.Error())
		return []domain.UserDomain{}, err
	}

	for _, userDB := range usersDB {
		usersDomain = append(usersDomain, userDB.CopyToUserDomain())
	}

	return usersDomain, nil
}

func (up UserRepository) GetUserByID(contextControl domain.ContextControl, userID uint64) (domain.UserDomain, error) {

	var userDB UserDB

	if err := up.db.WithContext(contextControl.BaseContext).
		Model(&UserDB{}).
		Where("id = ?", userID).
		First(&userDB).Error; err != nil {
		up.loggerSugar.Errorw(constants.ErrorToGetUserByID, constants.Error, err.Error())
		return domain.UserDomain{}, err
	}

	return userDB.CopyToUserDomain(), nil
}

func (up UserRepository) GetUserByUsername(contextControl domain.ContextControl, username string) (domain.UserDomain, error) {
	var userDB UserDB

	if err := up.db.WithContext(contextControl.BaseContext).
		Select("id, username, email, password, created_at").
		Where("username = ?", username).
		First(&userDB).Error; err != nil {

		up.loggerSugar.Errorw(constants.ErrorToGetUserByUsername, constants.Error, err.Error())
		return domain.UserDomain{}, err
	}

	return userDB.CopyToUserDomain(), nil
}

func (up UserRepository) GetUserByEmail(ctx domain.ContextControl, email string) (domain.UserDomain, error) {
	var userDB UserDB

	if err := up.db.WithContext(ctx.BaseContext).
		Select("id, email, created_at").
		Where("email = ?", email).
		First(&userDB).Error; err != nil {
		up.loggerSugar.Errorw(constants.ErrorToGetUserByEmail, constants.Error, err.Error())
		return domain.UserDomain{}, err
	}

	return userDB.CopyToUserDomain(), nil
}

func (up UserRepository) UpdateUser(ctx domain.ContextControl, userID uint64, fields map[string]interface{}) (domain.UserDomain, error) {
	delete(fields, constants.CreatedAt)

	if err := up.db.WithContext(ctx.BaseContext).
		Model(&UserDB{}).
		Where("id = ?", userID).
		Updates(fields).Error; err != nil {
		up.loggerSugar.Errorw(constants.ErrorToUpdateUser, constants.Error, err.Error())
		return domain.UserDomain{}, err
	}

	return up.GetUserByID(ctx, userID)
}

func (up UserRepository) SoftDeleteUser(ctx domain.ContextControl, userID uint64) error {
	fields := map[string]interface{}{
		constants.DeletedAt: time.Now().UTC(),
		constants.UpdatedAt: time.Now().UTC(),
	}

	if err := up.db.WithContext(ctx.BaseContext).
		Model(&UserDB{}).
		Where("id = ?", userID).
		Updates(fields).Error; err != nil {
		up.loggerSugar.Errorw(constants.ErrorToSoftDeleteUser, constants.Error, err.Error())
		return err
	}

	return nil
}
