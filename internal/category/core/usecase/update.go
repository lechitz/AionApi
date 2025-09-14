package usecase

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/category/core/domain"
	"github.com/lechitz/AionApi/internal/category/core/ports/input"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Update updates an existing Category in the system using the provided UpdateCategoryCommand.
// It builds the update fields, calls the repository, and logs/traces the operation.
func (s *Service) Update(ctx context.Context, cmd input.UpdateCategoryCommand) (domain.Category, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanUpdateCategory)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanUpdateCategory),
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(cmd.ID, 10)),
		attribute.String(commonkeys.UserID, strconv.FormatUint(cmd.UserID, 10)),
	)

	fieldsToUpdate := extractUpdateFields(cmd)

	span.AddEvent(EventRepositoryUpdate)
	updatedCategory, err := s.CategoryRepository.UpdateCategory(ctx, cmd.ID, cmd.UserID, fieldsToUpdate)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, FailedToUpdateCategory)
		s.Logger.ErrorwCtx(
			ctx,
			FailedToUpdateCategory,
			commonkeys.CategoryID, strconv.FormatUint(cmd.ID, 10),
			commonkeys.Error, err,
		)
		return domain.Category{}, err
	}

	span.AddEvent(EventSuccess)
	span.SetStatus(codes.Ok, "updated")
	s.Logger.InfowCtx(
		ctx,
		SuccessfullyUpdatedCategory,
		commonkeys.CategoryID, strconv.FormatUint(updatedCategory.ID, 10),
	)

	return updatedCategory, nil
}

// extractUpdateFields builds a map with only the non-nil/non-empty fields from UpdateCategoryCommand.
func extractUpdateFields(cmd input.UpdateCategoryCommand) map[string]interface{} {
	updateFields := make(map[string]interface{})

	if cmd.Name != nil && *cmd.Name != "" {
		updateFields[commonkeys.CategoryName] = *cmd.Name
	}
	if cmd.Description != nil && *cmd.Description != "" {
		updateFields[commonkeys.CategoryDescription] = *cmd.Description
	}
	if cmd.ColorHex != nil && *cmd.ColorHex != "" {
		updateFields[commonkeys.CategoryColor] = *cmd.ColorHex
	}
	if cmd.Icon != nil && *cmd.Icon != "" {
		updateFields[commonkeys.CategoryIcon] = *cmd.Icon
	}

	return updateFields
}
