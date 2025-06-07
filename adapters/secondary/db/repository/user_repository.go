package repository

import (
	"context"
	"time"

	"github.com/lechitz/AionApi/adapters/secondary/db/constants"
	"github.com/lechitz/AionApi/adapters/secondary/db/mapper"
	"github.com/lechitz/AionApi/adapters/secondary/db/model"

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
	userDB := mapper.UserToDB(userDomain)

	if err := up.db.WithContext(ctx).
		Create(&userDB).Error; err != nil {
		return domain.UserDomain{}, err
	}

	return mapper.UserFromDB(userDB), nil
}

func (up UserRepository) GetAllUsers(ctx context.Context) ([]domain.UserDomain, error) {
	var usersDB []model.UserDB
	var usersDomain []domain.UserDomain

	if err := up.db.WithContext(ctx).
		Model(&model.UserDB{}).
		Select("user_id, name, username, email, created_at").
		Find(&usersDB).Error; err != nil {
		return nil, err
	}

	for _, userDB := range usersDB {
		usersDomain = append(usersDomain, mapper.UserFromDB(userDB))
	}

	return usersDomain, nil
}

func (up UserRepository) GetUserByID(ctx context.Context, userID uint64) (domain.UserDomain, error) {
	var userDB model.UserDB

	if err := up.db.WithContext(ctx).
		Model(&model.UserDB{}).
		Where("user_id = ?", userID).
		First(&userDB).Error; err != nil {
		return domain.UserDomain{}, err
	}

	return mapper.UserFromDB(userDB), nil
}

func (up UserRepository) GetUserByUsername(ctx context.Context, username string) (domain.UserDomain, error) {
	var userDB model.UserDB

	if err := up.db.WithContext(ctx).
		Select("user_id, username, email, password, created_at").
		Where("username = ?", username).
		First(&userDB).Error; err != nil {

		return domain.UserDomain{}, err
	}

	return mapper.UserFromDB(userDB), nil
}

func (up UserRepository) GetUserByEmail(ctx context.Context, email string) (domain.UserDomain, error) {
	var userDB model.UserDB

	if err := up.db.WithContext(ctx).
		Select("user_id, email, created_at").
		Where("email = ?", email).
		First(&userDB).Error; err != nil {
		return domain.UserDomain{}, err
	}

	return mapper.UserFromDB(userDB), nil
}

func (up UserRepository) UpdateUser(
	ctx context.Context,
	userID uint64,
	fields map[string]interface{},
) (domain.UserDomain, error) {
	delete(fields, constants.CreatedAt)

	if err := up.db.WithContext(ctx).
		Model(&model.UserDB{}).
		Where("user_id = ?", userID).
		Updates(fields).Error; err != nil {
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
		Where("user_id = ?", userID).
		Updates(fields).Error; err != nil {
		return err
	}

	return nil
}
