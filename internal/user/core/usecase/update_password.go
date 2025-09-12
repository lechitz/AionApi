package usecase

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/auth/core/domain"
	"github.com/lechitz/AionApi/internal/platform/server/http/helpers/sharederrors"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// UpdatePassword updates a user's password after validating the old password and hashing the new password, then returns the updated user and a new token.
func (s *Service) UpdatePassword(ctx context.Context, userID uint64, oldPassword, newPassword string) (string, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, SpanUpdateUserPassword)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanUpdateUserPassword),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	)

	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, ErrorToGetSelf))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, ErrorToGetSelf, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(user.ID, 10))
		return "", fmt.Errorf("%s: %w", ErrorToGetSelf, err)
	}

	if err := s.hasher.Compare(user.Password, oldPassword); err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, ErrorToCompareHashAndPassword))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, ErrorToCompareHashAndPassword, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(user.ID, 10))
		return "", fmt.Errorf("%s: %w", ErrorToCompareHashAndPassword, err)
	}

	hashedPassword, err := s.hasher.Hash(newPassword)
	if err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, ErrorToHashPassword))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, ErrorToHashPassword, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(user.ID, 10))
		return "", fmt.Errorf("%s: %w", ErrorToHashPassword, err)
	}

	fields := map[string]interface{}{
		commonkeys.Password:      hashedPassword,
		commonkeys.UserUpdatedAt: time.Now().UTC(),
	}

	updatedUser, err := s.userRepository.Update(ctx, user.ID, fields)
	if err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, ErrorToUpdatePassword))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, ErrorToUpdatePassword, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(user.ID, 10))
		return "", fmt.Errorf("%s: %w", ErrorToUpdatePassword, err)
	}

	tokenValue, err := s.tokenProvider.Generate(updatedUser.ID)
	if err != nil || tokenValue == "" {
		span.SetAttributes(attribute.String(commonkeys.Status, sharederrors.ErrMsgCreateToken))
		if err != nil {
			span.RecordError(err)
			s.logger.ErrorwCtx(ctx, ErrorToCreateToken, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10))
			return "", fmt.Errorf("%s: %w", ErrorToCreateToken, err)
		}
		s.logger.ErrorwCtx(ctx, ErrorToCreateToken, commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10))
		return "", fmt.Errorf("%s: empty token", ErrorToCreateToken)
	}

	if err := s.authStore.Save(ctx, domain.Auth{Key: updatedUser.ID, Token: tokenValue}); err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, sharederrors.ErrMsgCreateToken))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, ErrorToCreateToken, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10))
		return "", fmt.Errorf("%s: %w", ErrorToCreateToken, err)
	}

	span.SetAttributes(
		attribute.String(commonkeys.Status, commonkeys.StatusSuccess),
		attribute.String(commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10)),
	)
	s.logger.InfowCtx(ctx, SuccessPasswordUpdated, commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10))

	return tokenValue, nil
}
