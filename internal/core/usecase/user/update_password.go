package user

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/sharederrors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// UpdatePassword updates a user's password after validating the old password and hashing the new password, then returns the updated user and a new token.
func (s *Service) UpdatePassword(ctx context.Context, userID uint64, oldPassword, newPassword string) (string, error) {
	tracer := otel.Tracer(constants.TracerName)
	ctx, span := tracer.Start(ctx, constants.SpanUpdateUserPassword)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, constants.SpanUpdateUserPassword),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	)

	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, constants.ErrorToGetSelf))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, constants.ErrorToGetSelf, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(user.ID, 10))
		return "", fmt.Errorf("%s: %w", constants.ErrorToGetSelf, err)
	}

	if err := s.hasher.Compare(user.Password, oldPassword); err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, constants.ErrorToCompareHashAndPassword))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, constants.ErrorToCompareHashAndPassword, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(user.ID, 10))
		return "", fmt.Errorf("%s: %w", constants.ErrorToCompareHashAndPassword, err)
	}

	hashedPassword, err := s.hasher.Hash(newPassword)
	if err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, constants.ErrorToHashPassword))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, constants.ErrorToHashPassword, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(user.ID, 10))
		return "", fmt.Errorf("%s: %w", constants.ErrorToHashPassword, err)
	}

	fields := map[string]interface{}{
		commonkeys.Password:      hashedPassword,
		commonkeys.UserUpdatedAt: time.Now().UTC(),
	}

	updatedUser, err := s.userRepository.Update(ctx, user.ID, fields)
	if err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, constants.ErrorToUpdatePassword))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, constants.ErrorToUpdatePassword, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(user.ID, 10))
		return "", fmt.Errorf("%s: %w", constants.ErrorToUpdatePassword, err)
	}

	tokenValue, err := s.tokenProvider.Generate(ctx, updatedUser.ID)
	if err != nil || tokenValue == "" {
		span.SetAttributes(attribute.String(commonkeys.Status, sharederrors.ErrMsgCreateToken))
		if err != nil {
			span.RecordError(err)
			s.logger.ErrorwCtx(ctx, constants.ErrorToCreateToken, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10))
			return "", fmt.Errorf("%s: %w", constants.ErrorToCreateToken, err)
		}
		s.logger.ErrorwCtx(ctx, constants.ErrorToCreateToken, commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10))
		return "", fmt.Errorf("%s: empty token", constants.ErrorToCreateToken)
	}

	if err := s.tokenStore.Save(ctx, domain.Token{Key: updatedUser.ID, Value: tokenValue}); err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, sharederrors.ErrMsgCreateToken))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, constants.ErrorToCreateToken, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10))
		return "", fmt.Errorf("%s: %w", constants.ErrorToCreateToken, err)
	}

	span.SetAttributes(
		attribute.String(commonkeys.Status, commonkeys.StatusSuccess),
		attribute.String(commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10)),
	)
	s.logger.InfowCtx(ctx, constants.SuccessPasswordUpdated, commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10))

	return tokenValue, nil
}
