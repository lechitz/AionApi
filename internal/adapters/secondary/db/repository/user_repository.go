package repository

import (
	"context"
	"github.com/lechitz/AionApi/internal/adapters/secondary/db/constants"
	"github.com/lechitz/AionApi/internal/adapters/secondary/db/mapper"
	"github.com/lechitz/AionApi/internal/adapters/secondary/db/model"
	"time"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	"gorm.io/gorm"
)

// UserRepository handles interactions with the user database, providing methods for CRUD operations and user retrieval.
type UserRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewUserRepository initializes a new UserRepository with the provided database connection and logger.
func NewUserRepository(db *gorm.DB, logger logger.Logger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

// CreateUser adds a new user to the database, mapping the provided domain object and returning the created user or an error if the operation fails.
func (up UserRepository) CreateUser(ctx context.Context, userDomain domain.UserDomain) (domain.UserDomain, error) {
	userDB := mapper.UserToDB(userDomain)

	if err := up.db.WithContext(ctx).
		Create(&userDB).Error; err != nil {
		return domain.UserDomain{}, err
	}

	return mapper.UserFromDB(userDB), nil
}

// GetAllUsers retrieves all active users from the database and maps them to the domain.UserDomain format. Returns a slice of users or an error.
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

// GetUserByID retrieves a user from the database by their unique user ID and returns the user in domain object format or an error.
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

// GetUserByUsername retrieves a user from the database using their unique username. Returns a domain.UserDomain or an error if the user is not found.
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

// GetUserByEmail retrieves a user by their email address from the database and returns a domain.UserDomain or an error if not found.
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

// UpdateUser updates specified fields for a user by their ID and returns the updated user or an error if the operation fails.
func (up UserRepository) UpdateUser(ctx context.Context, userID uint64, fields map[string]interface{}) (domain.UserDomain, error) {
	delete(fields, constants.CreatedAt)

	if err := up.db.WithContext(ctx).
		Model(&model.UserDB{}).
		Where("user_id = ?", userID).
		Updates(fields).Error; err != nil {
		return domain.UserDomain{}, err
	}

	return up.GetUserByID(ctx, userID)
}

// SoftDeleteUser marks a user as deleted by updating the DeletedAt and UpdatedAt fields for the specified userID. Returns an error if the update fails.
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
