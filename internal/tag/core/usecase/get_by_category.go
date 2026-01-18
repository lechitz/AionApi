package usecase

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/tag/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// GetByCategoryID retrieves all tags associated with a specific category for a given user.
func (s *Service) GetByCategoryID(ctx context.Context, categoryID uint64, userID uint64) ([]domain.Tag, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanGetByCategory)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanGetByCategory),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(categoryID, 10)),
	)

	span.AddEvent("CheckCache")
	cachedTags, err := s.TagCache.GetTagsByCategory(ctx, categoryID, userID)
	if err == nil && cachedTags != nil {
		span.AddEvent("CacheHit")
		return cachedTags, nil
	}

	span.AddEvent(EventRepositoryGet)
	tags, err := s.TagRepository.GetByCategoryID(ctx, categoryID, userID)
	if err != nil {
		span.RecordError(err)
		s.Logger.ErrorwCtx(ctx, ErrFailedToListTags, commonkeys.Error, err.Error(), commonkeys.UserID, userID, commonkeys.CategoryID, categoryID)
		return []domain.Tag{}, err
	}

	span.AddEvent("SaveToCache")
	err = s.TagCache.SaveTagsByCategory(ctx, categoryID, userID, tags, 0) // use default TTL
	if err != nil {
		s.Logger.WarnwCtx(ctx, "failed to save tags to cache", commonkeys.Error, err.Error(), commonkeys.UserID, userID, commonkeys.CategoryID, categoryID)
	}

	span.AddEvent(EventSuccess)
	return tags, nil
}
