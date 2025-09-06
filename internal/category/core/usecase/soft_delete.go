package usecase

import (
	"context"
	"errors"
	"strconv"

	"github.com/lechitz/AionApi/internal/core/category/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// SoftDelete performs a soft delete operation on a handler, marking it as inactive instead of permanently removing it from the database.
func (s *Service) SoftDelete(ctx context.Context, category domain.Category) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanSoftDeleteCategory)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanSoftDeleteCategory),
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(category.ID, 10)),
	)

	span.AddEvent(EventRepositoryDelete)
	if err := s.CategoryRepository.SoftDelete(ctx, category); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, FailedToSoftDeleteCategory)
		s.Logger.ErrorwCtx(ctx, FailedToSoftDeleteCategory, commonkeys.CategoryID, strconv.FormatUint(category.ID, 10), commonkeys.Error, err)
		return errors.New(FailedToSoftDeleteCategory)
	}

	span.AddEvent(EventSuccess)
	span.SetStatus(codes.Ok, "deleted")
	s.Logger.InfowCtx(ctx, SuccessfullySoftDeletedCategory, commonkeys.CategoryID, strconv.FormatUint(category.ID, 10))

	return nil
}
