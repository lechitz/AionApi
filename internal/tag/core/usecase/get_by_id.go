package usecase

import (
	"context"
	"errors"
	"strconv"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/tag/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetByID retrieves a handler by ID.
func (s *Service) GetByID(ctx context.Context, tagID uint64, userID uint64) (domain.Tag, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanGetTagByName)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanGetTagByName),
		attribute.String(commonkeys.TagName, strconv.FormatUint(tagID, 10)),
	)

	if userID == 0 {
		span.SetStatus(codes.Error, UserIDIsRequired)
		s.Logger.ErrorwCtx(ctx, UserIDIsRequired, commonkeys.UserID, strconv.FormatUint(userID, 10))
		return domain.Tag{}, errors.New(UserIDIsRequired)
	}

	if tagID == 0 {
		span.SetStatus(codes.Error, UserIDIsRequired)
		s.Logger.ErrorwCtx(ctx, UserIDIsRequired, commonkeys.TagName, strconv.FormatUint(tagID, 10))
		return domain.Tag{}, errors.New(UserIDIsRequired)
	}

	span.AddEvent("CheckCache")
	cachedTag, err := s.TagCache.GetTag(ctx, tagID, userID)
	if err == nil && cachedTag.ID != 0 {
		span.AddEvent("CacheHit")
		span.SetStatus(codes.Ok, StatusRetrievedByName)
		return cachedTag, nil
	}

	span.AddEvent(EventRepositoryGet)
	tag, err := s.TagRepository.GetByID(ctx, tagID, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, FailedToGetTagByName)
		s.Logger.ErrorwCtx(ctx, FailedToGetTagByName, commonkeys.TagName, strconv.FormatUint(tagID, 10), commonkeys.Error, err)
		return domain.Tag{}, err
	}

	span.AddEvent("SaveToCache")
	err = s.TagCache.SaveTag(ctx, tag, 0) // use default TTL
	if err != nil {
		s.Logger.WarnwCtx(ctx, "failed to save tag to cache", commonkeys.TagID, tagID, commonkeys.Error, err.Error())
	}
	span.AddEvent(EventSuccess)
	span.SetStatus(codes.Ok, StatusRetrievedByName)
	return tag, nil
}
