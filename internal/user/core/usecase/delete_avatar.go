package usecase

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/user/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// RemoveAvatar clears the current avatar URL for the authenticated user.
func (s *Service) RemoveAvatar(ctx context.Context, userID uint64) (domain.User, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, SpanDeleteAvatar)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanDeleteAvatar),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	)

	fields := map[string]interface{}{
		commonkeys.AvatarURL:     nil,
		commonkeys.UserUpdatedAt: time.Now().UTC(),
	}

	updatedUser, err := s.userRepository.Update(ctx, userID, fields)
	if err != nil {
		span.SetStatus(codes.Error, ErrorToDeleteAvatar)
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, ErrorToDeleteAvatar, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(userID, 10))
		return domain.User{}, fmt.Errorf("%w: %w", ErrDeleteAvatar, err)
	}

	span.AddEvent(SpanEventInvalidateCache)
	if err := s.userCache.DeleteUser(ctx, updatedUser.ID, updatedUser.Username, updatedUser.Email); err != nil {
		s.logger.WarnwCtx(ctx, WarnFailedToInvalidateUserCacheAfterAvatarDelete,
			commonkeys.UserID, updatedUser.ID,
			commonkeys.Username, updatedUser.Username,
			commonkeys.Email, updatedUser.Email,
			commonkeys.Error, err,
		)
	}

	span.SetAttributes(
		attribute.String(commonkeys.Status, commonkeys.StatusSuccess),
		attribute.String(commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10)),
	)
	s.logger.InfowCtx(ctx, SuccessUserAvatarDeleted, commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10))

	return updatedUser, nil
}
