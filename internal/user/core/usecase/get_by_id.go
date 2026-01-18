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

// GetByID retrieves a user by the provided userID.
func (s *Service) GetByID(ctx context.Context, userID uint64) (domain.User, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, SpanGetSelf)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanGetSelf),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	)

	span.AddEvent(SpanEventCheckCache)
	cachedUser, err := s.userCache.GetUserByID(ctx, userID)
	if err == nil && cachedUser.ID != 0 {
		span.AddEvent(SpanEventCacheHit)
		span.SetStatus(codes.Ok, SuccessUserRetrieved)
		s.logger.InfowCtx(ctx, InfoUserRetrievedFromCache,
			commonkeys.UserID, strconv.FormatUint(userID, 10),
		)
		return cachedUser, nil
	}

	span.AddEvent(SpanEventCacheMiss)
	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrorToGetSelf)
		s.logger.ErrorwCtx(ctx, ErrorToGetSelf,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, strconv.FormatUint(userID, 10),
		)
		return domain.User{}, fmt.Errorf("%w: %w", ErrGetSelf, err)
	}

	span.AddEvent(SpanEventSaveToCache)
	if err := s.userCache.SaveUser(ctx, user, 0); err != nil {
		s.logger.WarnwCtx(ctx, WarnFailedToSaveUserToCacheGeneric,
			commonkeys.UserID, userID,
			commonkeys.Error, err,
		)
	}

	span.SetStatus(codes.Ok, SuccessUserRetrieved)
	s.logger.InfowCtx(ctx, SuccessUserRetrieved,
		commonkeys.UserID, strconv.FormatUint(user.ID, 10),
	)
	return user, nil
}
