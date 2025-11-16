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

// GetAll retrieves all tags for a given user.
func (s *Service) GetAll(ctx context.Context, userID uint64) ([]domain.Tag, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanGetAll)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanGetAll),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	)

	if userID == 0 {
		span.SetStatus(codes.Error, UserIDIsRequired)
		s.Logger.ErrorwCtx(ctx, UserIDIsRequired, commonkeys.UserID, strconv.FormatUint(userID, 10))
		return []domain.Tag{}, errors.New(UserIDIsRequired)
	}

	span.AddEvent(EventRepositoryListAll)
	tags, err := s.TagRepository.GetAll(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, FailedToListTags)
		s.Logger.ErrorwCtx(ctx, FailedToListTags, commonkeys.UserID, strconv.FormatUint(userID, 10), commonkeys.Error, err)
		return []domain.Tag{}, err
	}

	span.AddEvent(EventSuccess)
	span.SetStatus(codes.Ok, StatusListedAll)
	return tags, nil
}
