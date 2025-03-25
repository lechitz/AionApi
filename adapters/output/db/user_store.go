package db

import (
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/lechitz/AionApi/core/domain"
	"github.com/lechitz/AionApi/pkg/contextkeys"
	"github.com/lechitz/AionApi/ports/output/db"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type UserStore struct {
	db          *gorm.DB
	loggerSugar *zap.SugaredLogger
}

func NewUserRepo(db *gorm.DB, loggerSugar *zap.SugaredLogger) db.IUserRepository {
	return &UserStore{
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
	return TableUsers
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

func (up UserStore) CreateUser(contextControl domain.ContextControl, userDomain domain.UserDomain) (domain.UserDomain, error) {

	var userDB UserDB
	copier.Copy(&userDB, &userDomain)

	if err := up.db.WithContext(contextControl.BaseContext).
		Create(&userDB).Error; err != nil {
		wrappedErr := fmt.Errorf(ErrorToCreateUser, err)
		up.loggerSugar.Errorw(ErrorToCreateUser, contextkeys.Error, wrappedErr.Error())
		return domain.UserDomain{}, wrappedErr
	}

	return userDB.CopyToUserDomain(), nil
}

func (up UserStore) GetAllUsers(contextControl domain.ContextControl) ([]domain.UserDomain, error) {
	var usersDB []UserDB
	var usersDomain []domain.UserDomain

	if err := up.db.WithContext(contextControl.BaseContext).
		Where(DeleteAtIsNull).
		Select(contextkeys.UserID, contextkeys.Name, contextkeys.Username, contextkeys.Email, contextkeys.CreatedAt).
		Find(&usersDB).Error; err != nil {
		up.loggerSugar.Errorw(ErrorToGetAllUsers, contextkeys.Error, err.Error())
		return []domain.UserDomain{}, err
	}

	for _, userDB := range usersDB {
		usersDomain = append(usersDomain, userDB.CopyToUserDomain())
	}

	return usersDomain, nil
}

func (up UserStore) GetUserByID(contextControl domain.ContextControl, userID uint64) (domain.UserDomain, error) {

	var userDB UserDB

	if err := up.db.WithContext(contextControl.BaseContext).
		Where("id = ?", userID).
		First(&userDB).Error; err != nil {
		up.loggerSugar.Errorw(ErrorToGetUserByID, contextkeys.Error, err.Error())
		return domain.UserDomain{}, err
	}

	return userDB.CopyToUserDomain(), nil
}

func (up UserStore) GetUserByUsername(contextControl domain.ContextControl, username string) (domain.UserDomain, error) {

	var userDB UserDB

	if err := up.db.WithContext(contextControl.BaseContext).
		Select("id, username, password").
		Where("username = ?", username).
		First(&userDB).Error; err != nil {
		up.loggerSugar.Errorw(ErrorToGetUserByUserName, contextkeys.Error, err.Error())
		return domain.UserDomain{}, err
	}

	return userDB.CopyToUserDomain(), nil
}

func (up UserStore) GetUserByEmail(ctx domain.ContextControl, email string) (domain.UserDomain, error) {
	var userDB UserDB

	if err := up.db.WithContext(ctx.BaseContext).
		Select("id, email").
		Where("email = ?", email).
		First(&userDB).Error; err != nil {
		up.loggerSugar.Errorw("error to get user by email", contextkeys.Error, err.Error())
		return domain.UserDomain{}, err
	}

	return userDB.CopyToUserDomain(), nil
}

func (up UserStore) UpdateUser(ctx domain.ContextControl, userID uint64, fields map[string]interface{}) (domain.UserDomain, error) {
	if err := up.db.WithContext(ctx.BaseContext).
		Model(&UserDB{}).
		Where("id = ?", userID).
		Updates(fields).Error; err != nil {
		up.loggerSugar.Errorw(ErrorToUpdateUser, contextkeys.Error, err.Error())
		return domain.UserDomain{}, err
	}

	return up.GetUserByID(ctx, userID)
}

func (up UserStore) SoftDeleteUser(ctx domain.ContextControl, userID uint64) error {
	fields := map[string]interface{}{
		contextkeys.DeletedAt: time.Now().UTC(),
		contextkeys.UpdatedAt: time.Now().UTC(),
	}

	if err := up.db.WithContext(ctx.BaseContext).
		Model(&UserDB{}).
		Where("id = ?", userID).
		Updates(fields).Error; err != nil {
		up.loggerSugar.Errorw(ErrorToSoftDeleteUser, contextkeys.Error, err.Error())
		return err
	}

	return nil
}
