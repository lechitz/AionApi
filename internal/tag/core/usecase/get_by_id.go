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

	if tagID == 0 {
		span.SetStatus(codes.Error, TagNameIsRequired)
		s.Logger.ErrorwCtx(ctx, TagNameIsRequired, commonkeys.TagName, strconv.FormatUint(tagID, 10))
		return domain.Tag{}, errors.New(TagNameIsRequired)
	}

	span.AddEvent(EventRepositoryGet)
	tag, err := s.TagRepository.GetByID(ctx, tagID, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, FailedToGetTagByName)
		s.Logger.ErrorwCtx(ctx, FailedToGetTagByName, commonkeys.TagName, strconv.FormatUint(tagID, 10), commonkeys.Error, err)
		return domain.Tag{}, err
	}

	span.AddEvent(EventSuccess)
	span.SetStatus(codes.Ok, StatusRetrievedByName)
	return tag, nil
}
