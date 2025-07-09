// Package repository provides methods for interacting with the user database.
package repository

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/shared/commonkeys"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/ports/output"

	"github.com/lechitz/AionApi/internal/adapters/secondary/db/mapper"
	"github.com/lechitz/AionApi/internal/adapters/secondary/db/model"

	"gorm.io/gorm"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// UserRepository handles interactions with the user database, providing methods for CRUD operations and user retrieval.
type UserRepository struct {
	db     *gorm.DB
	logger output.Logger
}

// NewUserRepository initializes a new UserRepository with the provided database connection and logger.
func NewUserRepository(db *gorm.DB, logger output.Logger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

// CreateUser adds a new user to the database, mapping the provided domain object and returning the created user or an error if the operation fails.
func (up UserRepository) CreateUser(ctx context.Context, userDomain domain.UserDomain) (domain.UserDomain, error) {
	tr := otel.Tracer("UserRepository")
	ctx, span := tr.Start(ctx, "CreateUser", trace.WithAttributes(
		attribute.String(commonkeys.Username, userDomain.Username),
		attribute.String(commonkeys.UserEmail, userDomain.Email),
		attribute.String("operation", "create"),
	))
	defer span.End()

	userDB := mapper.UserToDB(userDomain)

	if err := up.db.WithContext(ctx).
		Create(&userDB).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return domain.UserDomain{}, err
	}

	span.SetStatus(codes.Ok, "user created successfully")
	return mapper.UserFromDB(userDB), nil
}

// GetAllUsers retrieves all active users from the database and maps them to the domain.UserDomain format. Returns a slice of users or an error.
func (up UserRepository) GetAllUsers(ctx context.Context) ([]domain.UserDomain, error) {
	tr := otel.Tracer("UserRepository")
	ctx, span := tr.Start(ctx, "GetAllUsers", trace.WithAttributes(
		attribute.String("operation", "get_all"),
	))
	defer span.End()

	var usersDB []model.UserDB
	var usersDomain []domain.UserDomain

	if err := up.db.WithContext(ctx).
		Model(&model.UserDB{}).
		Select("user_id, name, username, email, created_at").
		Find(&usersDB).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return nil, err
	}

	for _, userDB := range usersDB {
		usersDomain = append(usersDomain, mapper.UserFromDB(userDB))
	}

	span.SetStatus(codes.Ok, "all users retrieved successfully")
	return usersDomain, nil
}

// GetUserByID retrieves a user from the database by their unique user ID and returns the user in domain object format or an error.
func (up UserRepository) GetUserByID(ctx context.Context, userID uint64) (domain.UserDomain, error) {
	tr := otel.Tracer("UserRepository")
	ctx, span := tr.Start(ctx, "GetUserByID", trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	))
	defer span.End()

	var userDB model.UserDB

	if err := up.db.WithContext(ctx).
		Model(&model.UserDB{}).
		Where("user_id = ?", userID).
		First(&userDB).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		return domain.UserDomain{}, err
	}

	span.SetStatus(codes.Ok, "user retrieved by id successfully")

	return mapper.UserFromDB(userDB), nil
}

// GetUserByUsername retrieves a user from the database using their unique username. Returns a domain.UserDomain or an error if the user is not found.
//
//nolint:dupl // TODO: Refactor duplication with GetUserByEmail / GetUserByUsername when business logic diverges or for greater DRY. Prioritizing explicitness and speed for now.
func (up UserRepository) GetUserByUsername(ctx context.Context, username string) (domain.UserDomain, error) {
	tr := otel.Tracer("UserRepository")
	ctx, span := tr.Start(ctx, "GetUserByUsername", trace.WithAttributes(
		attribute.String(commonkeys.Username, username),
		attribute.String("operation", "get_by_username"),
	))
	defer span.End()

	var userDB model.UserDB

	err := up.db.WithContext(ctx).
		Select("user_id, username, email, password, created_at").
		Where("username = ?", username).
		First(&userDB).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.SetStatus(codes.Ok, "user not found (business as usual)")
			return domain.UserDomain{}, nil
		}

		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return domain.UserDomain{}, err
	}

	span.SetStatus(codes.Ok, "user retrieved by username successfully")
	return mapper.UserFromDB(userDB), nil
}

// GetUserByEmail retrieves a user by their email address from the database and returns a domain.UserDomain or nil if not found.
//
//nolint:dupl // TODO:" Refactor duplication with GetUserByEmail / GetUserByUsername when business logic diverges or for greater DRY. Prioritizing explicitness and speed for now.
func (up UserRepository) GetUserByEmail(ctx context.Context, email string) (domain.UserDomain, error) {
	tr := otel.Tracer("UserRepository")
	ctx, span := tr.Start(ctx, "GetUserByEmail", trace.WithAttributes(
		attribute.String(commonkeys.UserEmail, email),
		attribute.String("operation", "get_by_email"),
	))
	defer span.End()

	var userDB model.UserDB

	err := up.db.WithContext(ctx).
		Select("user_id, email, created_at").
		Where("email = ?", email).
		First(&userDB).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.SetStatus(codes.Ok, "user not found (business as usual)")
			return domain.UserDomain{}, nil
		}

		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return domain.UserDomain{}, err
	}

	span.SetStatus(codes.Ok, "user retrieved by email successfully")
	return mapper.UserFromDB(userDB), nil
}

// UpdateUser updates specified fields for a user by their ID and returns the updated user or an error if the operation fails.
func (up UserRepository) UpdateUser(
	ctx context.Context,
	userID uint64,
	fields map[string]interface{},
) (domain.UserDomain, error) {
	tr := otel.Tracer("UserRepository")
	ctx, span := tr.Start(ctx, "UpdateUser", trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String("operation", "update"),
	))
	defer span.End()

	delete(fields, commonkeys.CreatedAt)

	if err := up.db.WithContext(ctx).
		Model(&model.UserDB{}).
		Where("user_id = ?", userID).
		Updates(fields).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return domain.UserDomain{}, err
	}

	span.SetStatus(codes.Ok, "user updated successfully")

	return up.GetUserByID(ctx, userID)
}

// SoftDeleteUser marks a user as deleted by updating the DeletedAt and UserUpdatedAt fields for the specified userID. Returns an error if the update fails.
func (up UserRepository) SoftDeleteUser(ctx context.Context, userID uint64) error {
	tr := otel.Tracer("UserRepository")
	ctx, span := tr.Start(ctx, "SoftDeleteUser", trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String("operation", "soft_delete"),
	))
	defer span.End()

	fields := map[string]interface{}{
		commonkeys.DeletedAt: time.Now().UTC(),
		commonkeys.UpdatedAt: time.Now().UTC(),
	}

	if err := up.db.WithContext(ctx).
		Model(&model.UserDB{}).
		Where("user_id = ?", userID).
		Updates(fields).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return err
	}

	span.SetStatus(codes.Ok, "user soft deleted successfully")

	return nil
}
