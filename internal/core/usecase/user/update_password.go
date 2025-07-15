package user

import (
	"context"
	"fmt"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/sharederrors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"strconv"
	"time"
)

//TODO: AVALIAR SE POSSO RETORNAR SÓ O userID ao invés do domínio todo.

// UpdateUserPassword updates a user's password after validating the old password and hashing the new password, then returns the updated user and a new token.
func (s *Service) UpdateUserPassword(ctx context.Context, user domain.User, oldPassword, newPassword string) (domain.User, string, error) {
	tracer := otel.Tracer(constants.TracerName)
	ctx, span := tracer.Start(ctx, constants.SpanUpdateUserPassword)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, constants.SpanUpdateUserPassword),
		attribute.String(commonkeys.UserID, strconv.FormatUint(user.ID, 10)),
	)

	userDB, err := s.userRepository.GetUserByID(ctx, user.ID)
	if err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, constants.ErrorToGetUserByID))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, constants.ErrorToGetUserByID, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(user.ID, 10))
		return domain.User{}, "", fmt.Errorf("%s: %w", constants.ErrorToGetUserByID, err)
	}

	if err := s.hashStore.Compare(userDB.Password, oldPassword); err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, constants.ErrorToCompareHashAndPassword))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, constants.ErrorToCompareHashAndPassword, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(user.ID, 10))
		return domain.User{}, "", fmt.Errorf("%s: %w", constants.ErrorToCompareHashAndPassword, err)
	}

	hashedPassword, err := s.hashStore.Hash(newPassword)
	if err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, constants.ErrorToHashPassword))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, constants.ErrorToHashPassword, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(user.ID, 10))
		return domain.User{}, "", fmt.Errorf("%s: %w", constants.ErrorToHashPassword, err)
	}

	fields := map[string]interface{}{
		commonkeys.Password:      hashedPassword,
		commonkeys.UserUpdatedAt: time.Now().UTC(),
	}

	updatedUser, err := s.userRepository.UpdateUser(ctx, user.ID, fields)
	if err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, constants.ErrorToUpdatePassword))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, constants.ErrorToUpdatePassword, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(user.ID, 10))
		return domain.User{}, "", fmt.Errorf("%s: %w", constants.ErrorToUpdatePassword, err)
	}

	token, err := s.tokenService.CreateToken(ctx, updatedUser.ID)
	if err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, sharederrors.ErrMsgCreateToken))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, constants.ErrorToCreateToken, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(user.ID, 10))
		return domain.User{}, "", fmt.Errorf("%s: %w", constants.ErrorToCreateToken, err)
	}

	span.SetAttributes(
		attribute.String(commonkeys.Status, commonkeys.StatusSuccess),
		attribute.String(commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10)),
	)
	s.logger.InfowCtx(ctx, constants.SuccessPasswordUpdated, commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10))

	return updatedUser, token.Token, nil
}
