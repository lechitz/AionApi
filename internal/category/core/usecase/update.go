package usecase

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/core/category/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Update updates an existing handler in the system with provided fields and logs the operation outcome.
// Returns the updated handler or an error.
func (s *Service) Update(ctx context.Context, category domain.Category) (domain.Category, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanUpdateCategory)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanUpdateCategory),
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(category.ID, 10)),
		attribute.String(commonkeys.UserID, strconv.FormatUint(category.UserID, 10)),
	)

	fieldsToUpdate := extractUpdateFields(category)

	span.AddEvent(EventRepositoryUpdate)
	updatedCategory, err := s.CategoryRepository.UpdateCategory(ctx, category.ID, category.UserID, fieldsToUpdate)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, FailedToUpdateCategory)
		s.Logger.ErrorwCtx(ctx, FailedToUpdateCategory, commonkeys.CategoryID, strconv.FormatUint(category.ID, 10), commonkeys.Error, err)
		return domain.Category{}, err
	}

	span.AddEvent(EventSuccess)
	span.SetStatus(codes.Ok, "updated")
	s.Logger.InfowCtx(ctx, SuccessfullyUpdatedCategory, commonkeys.CategoryID, strconv.FormatUint(updatedCategory.ID, 10))

	return updatedCategory, nil
}

// extractUpdateFields constructs a map of non-empty handler fields for updating.
func extractUpdateFields(category domain.Category) map[string]interface{} {
	updateFields := make(map[string]interface{})

	if category.Name != "" {
		updateFields[commonkeys.CategoryName] = category.Name
	}
	if category.Description != "" {
		updateFields[commonkeys.CategoryDescription] = category.Description
	}
	if category.Color != "" {
		updateFields[commonkeys.CategoryColor] = category.Color
	}
	if category.Icon != "" {
		updateFields[commonkeys.CategoryIcon] = category.Icon
	}

	return updateFields
}
