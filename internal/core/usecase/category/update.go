package category

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/category/constants"
)

// Update updates an existing category in the system with provided fields and logs the operation outcome.
// Returns the updated category or an error.
func (s *Service) Update(ctx context.Context, category domain.Category) (domain.Category, error) {
	tr := otel.Tracer(constants.TracerName)
	ctx, span := tr.Start(ctx, constants.SpanUpdateCategory)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, constants.SpanUpdateCategory),
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(category.ID, 10)),
		attribute.String(commonkeys.UserID, strconv.FormatUint(category.UserID, 10)),
	)

	fieldsToUpdate := extractUpdateFields(category)

	span.AddEvent(constants.EventRepositoryUpdate)
	updatedCategory, err := s.CategoryRepository.UpdateCategory(ctx, category.ID, category.UserID, fieldsToUpdate)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, constants.FailedToUpdateCategory)
		s.Logger.ErrorwCtx(ctx, constants.FailedToUpdateCategory, commonkeys.CategoryID, strconv.FormatUint(category.ID, 10), commonkeys.Error, err)
		return domain.Category{}, err
	}

	span.AddEvent(constants.EventSuccess)
	span.SetStatus(codes.Ok, "updated")
	s.Logger.InfowCtx(ctx, constants.SuccessfullyUpdatedCategory, commonkeys.CategoryID, strconv.FormatUint(updatedCategory.ID, 10))

	return updatedCategory, nil
}

// extractUpdateFields constructs a map of non-empty category fields for updating.
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
