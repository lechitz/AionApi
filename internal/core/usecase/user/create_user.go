// Package user contains use cases for managing users in the system.
package user

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/sharederrors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// TODO: PROVISÓRIO PRA BUILDAR.
const (
	DBErrorCheckingUsername = "a"
	ErrUsernameInUse        = "a"
	DBErrorCheckingEmail    = "a"
	ErrEmailInUse           = "a"

	ErrorToHashPassword       = "a"
	ErrorToCreateUser         = "a"
	ErrorToValidateCreateUser = "a"
	SuccessUserCreated        = "a"
)

// CreateUser creates a new user with the given data and password, ensuring validations and unique constraints are met. Returns the created user or an error.
func (s *Service) CreateUser(ctx context.Context, user domain.UserDomain, password string) (domain.UserDomain, error) {
	tracer := otel.Tracer("aionapi.user.usecase")
	ctx, span := tracer.Start(ctx, "CreateUser")
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", "CreateUser"),
		attribute.String(commonkeys.Username, user.Username),
		attribute.String(commonkeys.Email, user.Email),
	)

	s.normalizeUserData(&user)

	if err := s.validateCreateUserRequired(user, password); err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, "validation_failed"))
		s.logger.ErrorwCtx(ctx, constants.ErrorToValidateCreateUser, commonkeys.Error, err.Error())
		return domain.UserDomain{}, sharederrors.NewValidationError("user", err.Error())
	}

	existingByUsername, err := s.userStore.GetUserByUsername(ctx, user.Username)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, "db_error_checking_username"))
		s.logger.ErrorwCtx(ctx, DBErrorCheckingUsername, commonkeys.Error, err.Error())
		return domain.UserDomain{}, err
	}
	if existingByUsername.ID != 0 {
		span.SetAttributes(attribute.String(commonkeys.Status, "username_exists"))
		return domain.UserDomain{}, sharederrors.NewValidationError("username", ErrUsernameInUse)
	}

	existingByEmail, err := s.userStore.GetUserByEmail(ctx, user.Email)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, "db_error_checking_email"))
		s.logger.ErrorwCtx(ctx, DBErrorCheckingEmail, commonkeys.Error, err.Error())
		return domain.UserDomain{}, err
	}
	if existingByEmail.ID != 0 {
		span.SetAttributes(attribute.String(commonkeys.Status, "email_exists"))
		return domain.UserDomain{}, sharederrors.NewValidationError("email", ErrEmailInUse)
	}

	hashedPassword, err := s.hashStore.HashPassword(password)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, "hash_password_failed"))
		s.logger.ErrorwCtx(ctx, constants.ErrorToHashPassword, commonkeys.Error, err.Error())
		return domain.UserDomain{}, err
	}
	user.Password = hashedPassword

	userDB, err := s.userStore.CreateUser(ctx, user)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, "db_error_create_user"))
		s.logger.ErrorwCtx(ctx, constants.ErrorToCreateUser, commonkeys.Error, err.Error())
		return domain.UserDomain{}, err
	}

	span.SetAttributes(
		attribute.String(commonkeys.Status, "success"),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userDB.ID, 10)),
	)
	s.logger.InfowCtx(ctx, constants.SuccessUserCreated, commonkeys.UserID, strconv.FormatUint(userDB.ID, 10))

	return userDB, nil
}
