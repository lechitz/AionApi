package postgres

import (
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/lechitz/AionApi/core/domain/entities"
	"github.com/lechitz/AionApi/pkg/contextkeys"
	portsdb "github.com/lechitz/AionApi/ports/output/db"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type UserStore struct {
	db          *gorm.DB
	loggerSugar *zap.SugaredLogger
}

func NewUserRepo(db *gorm.DB, loggerSugar *zap.SugaredLogger) portsdb.IUserRepository {
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

func (u UserDB) CopyToUserDomain() entities.UserDomain {
	return entities.UserDomain{
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

func (up UserStore) CreateUser(contextControl entities.ContextControl, userDomain entities.UserDomain) (entities.UserDomain, error) {

	var userDB UserDB
	copier.Copy(&userDB, &userDomain)

	if err := up.db.WithContext(contextControl.BaseContext).
		Create(&userDB).Error; err != nil {
		wrappedErr := fmt.Errorf(ErrorToCreateUser, err)
		up.loggerSugar.Errorw(ErrorToCreateUser, contextkeys.Error, wrappedErr.Error())
		return entities.UserDomain{}, wrappedErr
	}

	return userDB.CopyToUserDomain(), nil
}

func (up UserStore) GetAllUsers(contextControl entities.ContextControl) ([]entities.UserDomain, error) {
	var usersDB []UserDB
	var usersDomain []entities.UserDomain

	if err := up.db.WithContext(contextControl.BaseContext).
		Where(DeleteAtIsNull).
		Select(contextkeys.UserID, contextkeys.Name, contextkeys.Username, contextkeys.Email, contextkeys.CreatedAt).
		Find(&usersDB).Error; err != nil {
		up.loggerSugar.Errorw(ErrorToGetAllUsers, contextkeys.Error, err.Error())
		return []entities.UserDomain{}, err
	}

	for _, userDB := range usersDB {
		usersDomain = append(usersDomain, userDB.CopyToUserDomain())
	}

	return usersDomain, nil
}

func (up UserStore) GetUserByID(contextControl entities.ContextControl, userID uint64) (entities.UserDomain, error) {

	var userDB UserDB

	if err := up.db.WithContext(contextControl.BaseContext).
		Where("id = ?", userID).
		First(&userDB).Error; err != nil {
		up.loggerSugar.Errorw(ErrorToGetUserByID, contextkeys.Error, err.Error())
		return entities.UserDomain{}, err
	}

	return userDB.CopyToUserDomain(), nil
}

func (up UserStore) GetUserByUsername(contextControl entities.ContextControl, username string) (entities.UserDomain, error) {

	var userDB UserDB

	if err := up.db.WithContext(contextControl.BaseContext).
		Select("id, username, password").
		Where("username = ?", username).
		First(&userDB).Error; err != nil {
		up.loggerSugar.Errorw(ErrorToGetUserByUserName, contextkeys.Error, err.Error())
		return entities.UserDomain{}, err
	}

	return userDB.CopyToUserDomain(), nil
}

func (up UserStore) UpdateUser(contextControl entities.ContextControl, userDomain entities.UserDomain) (entities.UserDomain, error) {

	var userDB UserDB
	copier.Copy(&userDB, &userDomain)

	if err := up.db.WithContext(contextControl.BaseContext).
		Model(&UserDB{}).
		Where("id = ?", userDB.ID).
		Updates(userDB).Error; err != nil {
		up.loggerSugar.Errorw(ErrorToUpdateUser, contextkeys.Error, err.Error())
		return entities.UserDomain{}, err
	}

	return userDB.CopyToUserDomain(), nil
}

func (up UserStore) UpdatePassword(contextControl entities.ContextControl, userDomain entities.UserDomain) (entities.UserDomain, error) {

	var userDB UserDB
	copier.Copy(&userDB, &userDomain)

	if err := up.db.WithContext(contextControl.BaseContext).
		Model(&UserDB{}).
		Where("id = ?", userDB.ID).
		Update(contextkeys.Password, userDB.Password).Error; err != nil {
		up.loggerSugar.Errorw(ErrorToUpdatePassword, contextkeys.Error, err.Error())
		return entities.UserDomain{}, err
	}

	return userDB.CopyToUserDomain(), nil
}

func (up UserStore) SoftDeleteUser(contextControl entities.ContextControl, userID uint64) error {
	now := time.Now()
	if err := up.db.WithContext(contextControl.BaseContext).
		Model(&UserDB{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			contextkeys.DeletedAt: now,
		}).Error; err != nil {
		up.loggerSugar.Errorw(ErrorToSoftDeleteUser, contextkeys.UserID, userID, contextkeys.Error, err.Error())
		return err
	}

	up.loggerSugar.Infow(SuccesfullyDeletedUser, contextkeys.UserID, userID, contextkeys.DeletedAt, now)
	return nil
}
