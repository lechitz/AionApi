package usecase

import (
	"context"
	"errors"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/tag/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetByName retrieves a handler by its name from the database and returns it.
func (s *Service) GetByName(ctx context.Context, tagName string, userID uint64) (domain.Tag, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanGetTagByName)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanGetTagByName),
		attribute.String(commonkeys.TagName, tagName),
	)

	if tagName == "" {
		span.SetStatus(codes.Error, TagNameIsRequired)
		s.Logger.ErrorwCtx(ctx, TagNameIsRequired, commonkeys.TagName, tagName)
		return domain.Tag{}, errors.New(TagNameIsRequired)
	}

	span.AddEvent(EventRepositoryGet)
	tag, err := s.TagRepository.GetByName(ctx, tagName, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, FailedToGetTagByName)
		s.Logger.ErrorwCtx(ctx, FailedToGetTagByName, commonkeys.TagName, tagName, commonkeys.Error, err)
		return domain.Tag{}, err
	}

	span.AddEvent(EventSuccess)
	span.SetStatus(codes.Ok, StatusRetrievedByName)
	return tag, nil
}
