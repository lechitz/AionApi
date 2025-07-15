// Package repository provides methods for interacting with the user database.
package repository

import (
	"context"
	"errors"
	"github.com/lechitz/AionApi/internal/adapters/secondary/db/postgres/mapper"
	"strconv"
	"strings"
	"time"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"

	"github.com/lechitz/AionApi/internal/shared/sharederrors"

	"github.com/lechitz/AionApi/internal/adapters/secondary/db/model"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

// UserRepository handles interactions with the user database, providing methods for CRUD operations and user retrieval.
type UserRepository struct {
	db     *gorm.DB
	logger output.ContextLogger
}

// NewUser initializes a new UserRepository with the provided database connection and contextlogger.
func NewUser(db *gorm.DB, logger output.ContextLogger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

// isUniqueViolation detects a Postgres unique constraint violation and returns the affected field.
func isUniqueViolation(err error) (string, bool) {
	if err == nil {
		return "", false
	}
	msg := err.Error()
	// Postgres unique constraint violation
	switch {
	case strings.Contains(msg, "users_username_key"):
		return commonkeys.Username, true
	case strings.Contains(msg, "users_email_key"):
		return commonkeys.Email, true
	}
	return "", false
}

// CreateUser inserts a new user. Returns a ValidationError if username/email is already in use.
func (up UserRepository) CreateUser(ctx context.Context, userDomain domain.User) (domain.User, error) {
	tr := otel.Tracer("UserRepository")
	ctx, span := tr.Start(ctx, "CreateUser", trace.WithAttributes(
		attribute.String(commonkeys.Username, userDomain.Username),
		attribute.String(commonkeys.Email, userDomain.Email),
		attribute.String(commonkeys.Operation, "create"),
	))
	defer span.End()

	userDB := mapper.UserToDB(userDomain)

	if err := up.db.WithContext(ctx).Create(&userDB).Error; err != nil {
		if field, ok := isUniqueViolation(err); ok {
			span.SetStatus(codes.Error, "validation: "+field+" is already in use")
			span.SetAttributes(attribute.String("http.error_reason", field+"_already_exists"))
			span.RecordError(err)
			up.logger.ErrorwCtx(ctx, "unique constraint violation on create user", "field", field, commonkeys.Error, err.Error())
			return domain.User{}, sharederrors.NewValidationError(field, field+" is already in use")
		}
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		up.logger.ErrorwCtx(ctx, "failed to create user", commonkeys.Error, err.Error())
		return domain.User{}, err
	}

	span.SetStatus(codes.Ok, "user created successfully")
	up.logger.InfowCtx(ctx, "user created successfully", commonkeys.UserID, userDB.ID)
	return mapper.UserFromDB(userDB), nil
}

// GetAllUsers returns all users (excluding soft-deleted). Returns error if DB operation fails.
func (up UserRepository) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	tr := otel.Tracer("UserRepository")
	ctx, span := tr.Start(ctx, "GetAllUsers", trace.WithAttributes(
		attribute.String(commonkeys.Operation, "get_all"),
	))
	defer span.End()

	var usersDB []model.UserDB
	var usersDomain []domain.User

	if err := up.db.WithContext(ctx).
		Model(&model.UserDB{}).
		Select("user_id, name, username, email, created_at").
		Find(&usersDB).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		up.logger.ErrorwCtx(ctx, "failed to get all users", commonkeys.Error, err.Error())
		return nil, err
	}

	for _, userDB := range usersDB {
		usersDomain = append(usersDomain, mapper.UserFromDB(userDB))
	}

	span.SetStatus(codes.Ok, "all users retrieved successfully")
	up.logger.InfowCtx(ctx, "all users retrieved successfully", commonkeys.UsersCount, len(usersDomain))
	return usersDomain, nil
}

// GetUserByID retrieves a user from the database by their unique user ID and returns the user in domain object format or an error.
func (up UserRepository) GetUserByID(ctx context.Context, userID uint64) (domain.User, error) {
	tr := otel.Tracer("UserRepository")
	ctx, span := tr.Start(ctx, "GetUserByID", trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, "get_by_id"),
	))
	defer span.End()

	var userDB model.UserDB

	if err := up.db.WithContext(ctx).
		Model(&model.UserDB{}).
		Where("user_id = ?", userID).
		First(&userDB).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		up.logger.ErrorwCtx(ctx, "failed to get user by id", commonkeys.UserID, userID, commonkeys.Error, err.Error())
		return domain.User{}, err
	}

	span.SetStatus(codes.Ok, "user retrieved by id successfully")
	up.logger.InfowCtx(ctx, "user retrieved by id successfully", commonkeys.UserID, userDB.ID)
	return mapper.UserFromDB(userDB), nil
}

// GetUserByUsername retrieves a user from the database using their unique username. Returns a domain.User or an error if the user is not found.
//
//nolint:dupl // TODO: Refactor duplication with GetUserByEmail / GetUserByUsername when business logic diverges or for greater DRY. Prioritizing explicitness and speed for now.
func (up UserRepository) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
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
			up.logger.InfowCtx(ctx, "user not found by username", commonkeys.Username, username)
			return domain.User{}, nil
		}
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		up.logger.ErrorwCtx(ctx, "failed to get user by username", commonkeys.Username, username, commonkeys.Error, err.Error())
		return domain.User{}, err
	}

	span.SetStatus(codes.Ok, "user retrieved by username successfully")
	up.logger.InfowCtx(ctx, "user retrieved by username successfully", commonkeys.UserID, userDB.ID, commonkeys.Username, userDB.Username)
	return mapper.UserFromDB(userDB), nil
}

// GetUserByEmail retrieves a user by their email address from the database and returns a domain.User or nil if not found.
//
//nolint:dupl // TODO:" Refactor duplication with GetUserByEmail / GetUserByUsername when business logic diverges or for greater DRY. Prioritizing explicitness and speed for now.
func (up UserRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	tr := otel.Tracer("UserRepository")
	ctx, span := tr.Start(ctx, "GetUserByEmail", trace.WithAttributes(
		attribute.String(commonkeys.Email, email),
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
			up.logger.InfowCtx(ctx, "user not found by email", commonkeys.Email, email)
			return domain.User{}, nil
		}
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		up.logger.ErrorwCtx(ctx, "failed to get user by email", commonkeys.Email, email, commonkeys.Error, err.Error())
		return domain.User{}, err
	}

	span.SetStatus(codes.Ok, "user retrieved by email successfully")
	up.logger.InfowCtx(ctx, "user retrieved by email successfully", commonkeys.UserID, userDB.ID, commonkeys.Email, userDB.Email)
	return mapper.UserFromDB(userDB), nil
}

// UpdateUser updates specified fields for a user by their ID and returns the updated user or an error if the operation fails.
func (up UserRepository) UpdateUser(
	ctx context.Context,
	userID uint64,
	fields map[string]interface{},
) (domain.User, error) {
	tr := otel.Tracer("UserRepository")
	ctx, span := tr.Start(ctx, "UpdateUser", trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String("operation", "update"),
	))
	defer span.End()

	delete(fields, commonkeys.UserCreatedAt)

	if err := up.db.WithContext(ctx).
		Model(&model.UserDB{}).
		Where("user_id = ?", userID).
		Updates(fields).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		up.logger.ErrorwCtx(ctx, "failed to update user", commonkeys.UserID, userID, commonkeys.Error, err.Error())
		return domain.User{}, err
	}

	span.SetStatus(codes.Ok, "user updated successfully")
	up.logger.InfowCtx(ctx, "user updated successfully", commonkeys.UserID, userID)
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
		commonkeys.UserDeletedAt: time.Now().UTC(),
		commonkeys.UserUpdatedAt: time.Now().UTC(),
	}

	if err := up.db.WithContext(ctx).
		Model(&model.UserDB{}).
		Where("user_id = ?", userID).
		Updates(fields).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		up.logger.ErrorwCtx(ctx, "failed to soft delete user", commonkeys.UserID, userID, commonkeys.Error, err.Error())
		return err
	}

	span.SetStatus(codes.Ok, "user soft deleted successfully")
	up.logger.InfowCtx(ctx, "user soft deleted successfully", commonkeys.UserID, userID)
	return nil
}
