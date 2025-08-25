package category

import (
	"context"
	"errors"
	"strconv"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"github.com/lechitz/AionApi/internal/core/usecase/category/constants"
)

// SoftDelete performs a soft delete operation on a category, marking it as inactive instead of permanently removing it from the database.
func (s *Service) SoftDelete(ctx context.Context, category domain.Category) error {
	tr := otel.Tracer(constants.TracerName)
	ctx, span := tr.Start(ctx, constants.SpanSoftDeleteCategory)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, constants.SpanSoftDeleteCategory),
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(category.ID, 10)),
	)

	span.AddEvent(constants.EventRepositoryDelete)
	if err := s.CategoryRepository.SoftDelete(ctx, category); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, constants.FailedToSoftDeleteCategory)
		s.Logger.ErrorwCtx(ctx, constants.FailedToSoftDeleteCategory, commonkeys.CategoryID, strconv.FormatUint(category.ID, 10), commonkeys.Error, err)
		return errors.New(constants.FailedToSoftDeleteCategory)
	}

	span.AddEvent(constants.EventSuccess)
	span.SetStatus(codes.Ok, "deleted")
	s.Logger.InfowCtx(ctx, constants.SuccessfullySoftDeletedCategory, commonkeys.CategoryID, strconv.FormatUint(category.ID, 10))

	return nil
}
