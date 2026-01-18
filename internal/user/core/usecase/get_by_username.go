package usecase

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/user/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetUserByUsername retrieves a user by their username.
func (s *Service) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, SpanGetUserByUsername)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanGetUserByUsername),
		attribute.String(commonkeys.Username, username),
	)

	span.AddEvent(SpanEventCheckCache)
	cachedUser, err := s.userCache.GetUserByUsername(ctx, username)
	if err == nil && cachedUser.ID != 0 {
		span.AddEvent(SpanEventCacheHit)
		span.SetStatus(codes.Ok, SuccessUserRetrieved)
		s.logger.InfowCtx(ctx, InfoUserRetrievedFromCache,
			commonkeys.Username, username,
			commonkeys.UserID, strconv.FormatUint(cachedUser.ID, 10),
		)
		return cachedUser, nil
	}

	span.AddEvent(SpanEventCacheMiss)
	user, err := s.userRepository.GetByUsername(ctx, username)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrorToGetUserByUsername)
		s.logger.ErrorwCtx(ctx, ErrorToGetUserByUsername,
			commonkeys.Error, err.Error(),
			commonkeys.Username, username,
		)
		return domain.User{}, fmt.Errorf("%w: %w", ErrGetUserByUsername, err)
	}

	span.AddEvent(SpanEventSaveToCache)
	if err := s.userCache.SaveUser(ctx, user, 0); err != nil {
		s.logger.WarnwCtx(ctx, WarnFailedToSaveUserToCacheGeneric,
			commonkeys.UserID, user.ID,
			commonkeys.Username, username,
			commonkeys.Error, err,
		)
	}

	span.SetStatus(codes.Ok, SuccessUserRetrieved)
	s.logger.InfowCtx(ctx, SuccessUserRetrieved,
		commonkeys.UserID, strconv.FormatUint(user.ID, 10),
		commonkeys.Username, username,
	)
	return user, nil
}
