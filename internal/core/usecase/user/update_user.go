// Package user contains use cases for managing users in the system.
package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/lechitz/AionApi/internal/shared/sharederrors"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// UpdateUser updates an existing user's attributes based on the provided data. Returns the updated user or an error if the operation fails.
func (s *Service) UpdateUser(ctx context.Context, user domain.UserDomain) (domain.UserDomain, error) {
	tracer := otel.Tracer(constants.TracerName)
	ctx, span := tracer.Start(ctx, constants.SpanUpdateUser)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, constants.SpanUpdateUser),
		attribute.String(commonkeys.UserID, strconv.FormatUint(user.ID, 10)),
	)

	updateFields := buildUserUpdateFields(user)

	if len(updateFields) == 0 {
		span.SetAttributes(attribute.String(commonkeys.Status, sharederrors.NoFieldsToUpdate))
		err := errors.New(constants.ErrorNoFieldsToUpdate)
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, constants.ErrorNoFieldsToUpdate, commonkeys.Error, err.Error())
		return domain.UserDomain{}, err
	}

	updateFields[constants.UpdatedAt] = time.Now().UTC()

	updatedUser, err := s.userStore.UpdateUser(ctx, user.ID, updateFields)
	if err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, constants.ErrorToUpdateUser))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, constants.ErrorToUpdateUser, commonkeys.Error, err.Error())
		return domain.UserDomain{}, fmt.Errorf("%s: %w", constants.ErrorToUpdateUser, err)
	}

	span.SetAttributes(
		attribute.String(commonkeys.Status, constants.StatusSuccess),
		attribute.String(commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10)),
	)
	s.logger.InfowCtx(ctx, constants.SuccessUserUpdated, commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10), commonkeys.UserUpdatedFields, updateFields)

	return updatedUser, nil
}

// UpdateUserPassword updates a user's password after validating the old password and hashing the new password, then returns the updated user and a new token.
func (s *Service) UpdateUserPassword(ctx context.Context, user domain.UserDomain, oldPassword, newPassword string) (domain.UserDomain, string, error) {
	tracer := otel.Tracer(constants.TracerName)
	ctx, span := tracer.Start(ctx, constants.SpanUpdateUserPassword)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, constants.SpanUpdateUserPassword),
		attribute.String(commonkeys.UserID, strconv.FormatUint(user.ID, 10)),
	)

	userDB, err := s.userStore.GetUserByID(ctx, user.ID)
	if err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, constants.ErrorToGetUserByID))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, constants.ErrorToGetUserByID, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(user.ID, 10))
		return domain.UserDomain{}, "", fmt.Errorf("%s: %w", constants.ErrorToGetUserByID, err)
	}

	if err := s.hashStore.ValidatePassword(userDB.Password, oldPassword); err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, constants.ErrorToCompareHashAndPassword))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, constants.ErrorToCompareHashAndPassword, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(user.ID, 10))
		return domain.UserDomain{}, "", fmt.Errorf("%s: %w", constants.ErrorToCompareHashAndPassword, err)
	}

	hashedPassword, err := s.hashStore.HashPassword(newPassword)
	if err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, constants.ErrorToHashPassword))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, constants.ErrorToHashPassword, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(user.ID, 10))
		return domain.UserDomain{}, "", fmt.Errorf("%s: %w", constants.ErrorToHashPassword, err)
	}

	fields := map[string]interface{}{
		commonkeys.Password:      hashedPassword,
		commonkeys.UserUpdatedAt: time.Now().UTC(),
	}

	updatedUser, err := s.userStore.UpdateUser(ctx, user.ID, fields)
	if err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, constants.ErrorToUpdatePassword))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, constants.ErrorToUpdatePassword, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(user.ID, 10))
		return domain.UserDomain{}, "", fmt.Errorf("%s: %w", constants.ErrorToUpdatePassword, err)
	}

	tokenDomain := domain.TokenDomain{UserID: user.ID}

	token, err := s.tokenService.CreateToken(ctx, tokenDomain)
	if err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, sharederrors.ErrorToCreateToken))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, constants.ErrorToCreateToken, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(user.ID, 10))
		return domain.UserDomain{}, "", fmt.Errorf("%s: %w", constants.ErrorToCreateToken, err)
	}

	span.SetAttributes(
		attribute.String(commonkeys.Status, constants.StatusSuccess),
		attribute.String(commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10)),
	)
	s.logger.InfowCtx(ctx, constants.SuccessPasswordUpdated, commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10))

	return updatedUser, token, nil
}
