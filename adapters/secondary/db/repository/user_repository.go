package repository

import (
	"context"
	"fmt"
	"github.com/lechitz/AionApi/adapters/secondary/db/constants"
	"github.com/lechitz/AionApi/adapters/secondary/db/mapper"
	"github.com/lechitz/AionApi/adapters/secondary/db/model"
	"time"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	"gorm.io/gorm"
)

type UserRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewUserRepository(db *gorm.DB, logger logger.Logger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

func (up UserRepository) CreateUser(ctx context.Context, userDomain domain.UserDomain) (domain.UserDomain, error) {
	userDB := mapper.ToDB(userDomain)

	if err := up.db.WithContext(ctx).
		Create(&userDB).Error; err != nil {
		wrappedErr := fmt.Errorf(constants.ErrorToCreateUser, err)
		up.logger.Errorw(constants.ErrorToCreateUser, constants.Error, wrappedErr.Error())
		return domain.UserDomain{}, wrappedErr
	}

	return mapper.FromDB(userDB), nil
}

func (up UserRepository) GetAllUsers(ctx context.Context) ([]domain.UserDomain, error) {
	var usersDB []model.UserDB
	var usersDomain []domain.UserDomain

	if err := up.db.WithContext(ctx).
		Model(&model.UserDB{}).
		Select("id, name, username, email, created_at").
		Find(&usersDB).Error; err != nil {
		up.logger.Errorw(constants.ErrorToGetAllUsers, constants.Error, err.Error())
		return nil, err
	}

	for _, userDB := range usersDB {
		usersDomain = append(usersDomain, mapper.FromDB(userDB))
	}

	return usersDomain, nil
}

func (up UserRepository) GetUserByID(ctx context.Context, userID uint64) (domain.UserDomain, error) {
	var userDB model.UserDB

	if err := up.db.WithContext(ctx).
		Model(&model.UserDB{}).
		Where("id = ?", userID).
		First(&userDB).Error; err != nil {
		up.logger.Errorw(constants.ErrorToGetUserByID, constants.Error, err.Error())
		return domain.UserDomain{}, err
	}

	return mapper.FromDB(userDB), nil
}

func (up UserRepository) GetUserByUsername(ctx context.Context, username string) (domain.UserDomain, error) {
	var userDB model.UserDB

	if err := up.db.WithContext(ctx).
		Select("id, username, email, password, created_at").
		Where("username = ?", username).
		First(&userDB).Error; err != nil {
		up.logger.Errorw(constants.ErrorToGetUserByUsername, constants.Error, err.Error())
		return domain.UserDomain{}, err
	}

	return mapper.FromDB(userDB), nil
}

func (up UserRepository) GetUserByEmail(ctx context.Context, email string) (domain.UserDomain, error) {
	var userDB model.UserDB

	if err := up.db.WithContext(ctx).
		Select("id, email, created_at").
		Where("email = ?", email).
		First(&userDB).Error; err != nil {
		up.logger.Errorw(constants.ErrorToGetUserByEmail, constants.Error, err.Error())
		return domain.UserDomain{}, err
	}

	return mapper.FromDB(userDB), nil
}

func (up UserRepository) UpdateUser(ctx context.Context, userID uint64, fields map[string]interface{}) (domain.UserDomain, error) {
	delete(fields, constants.CreatedAt)

	if err := up.db.WithContext(ctx).
		Model(&model.UserDB{}).
		Where("id = ?", userID).
		Updates(fields).Error; err != nil {
		up.logger.Errorw(constants.ErrorToUpdateUser, constants.Error, err.Error())
		return domain.UserDomain{}, err
	}

	return up.GetUserByID(ctx, userID)
}

func (up UserRepository) SoftDeleteUser(ctx context.Context, userID uint64) error {
	fields := map[string]interface{}{
		constants.DeletedAt: time.Now().UTC(),
		constants.UpdatedAt: time.Now().UTC(),
	}

	if err := up.db.WithContext(ctx).
		Model(&model.UserDB{}).
		Where("id = ?", userID).
		Updates(fields).Error; err != nil {
		up.logger.Errorw(constants.ErrorToSoftDeleteUser, constants.Error, err.Error())
		return err
	}

	return nil
}
