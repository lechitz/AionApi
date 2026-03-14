package usecase

import (
	"context"
	"strconv"
	"strings"

	"github.com/lechitz/AionApi/internal/category/core/domain"
	"github.com/lechitz/AionApi/internal/category/core/ports/input"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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

	normalizedIcon, hasIcon, err := normalizeUpdateIcon(span, cmd.Icon)
	if err != nil {
		return domain.Category{}, err
	}
	if hasIcon {
		cmd.Icon = &normalizedIcon
	} else {
		cmd.Icon = nil
	}

	newName := extractNewName(cmd)
	if err := s.ensureUpdateNameUnique(ctx, span, cmd, newName); err != nil {
		return domain.Category{}, err
	}

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

	span.AddEvent(EventInvalidateCache)
	s.invalidateCategoryCaches(ctx, updatedCategory)

	span.AddEvent(EventSuccess)
	span.SetStatus(codes.Ok, StatusUpdated)
	s.Logger.InfowCtx(
		ctx,
		SuccessfullyUpdatedCategory,
		commonkeys.CategoryID, strconv.FormatUint(updatedCategory.ID, 10),
	)

	return updatedCategory, nil
}

func normalizeUpdateIcon(span trace.Span, icon *string) (string, bool, error) {
	if icon == nil {
		return "", false, nil
	}

	normalized := normalizeIconKey(icon)
	if normalized == "" {
		return "", false, nil
	}
	if !isValidIconKey(normalized) {
		span.SetStatus(codes.Error, ErrToValidateCategory)
		return "", false, ErrCategoryIconInvalid
	}

	return normalized, true, nil
}

func extractNewName(cmd input.UpdateCategoryCommand) string {
	if cmd.Name == nil {
		return ""
	}
	return strings.TrimSpace(*cmd.Name)
}

func (s *Service) ensureUpdateNameUnique(
	ctx context.Context,
	span trace.Span,
	cmd input.UpdateCategoryCommand,
	newName string,
) error {
	if newName == "" {
		return nil
	}

	span.AddEvent(EventCheckUniqueness)
	existing, err := s.CategoryRepository.GetByName(ctx, newName, cmd.UserID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, FailedToGetCategoryByName)
		s.Logger.ErrorwCtx(ctx, FailedToGetCategoryByName,
			commonkeys.CategoryName, newName,
			commonkeys.UserID, strconv.FormatUint(cmd.UserID, 10),
			commonkeys.Error, err,
		)
		return err
	}

	if existing.Name != "" && existing.ID != cmd.ID {
		span.SetStatus(codes.Error, CategoryAlreadyExists)
		s.Logger.ErrorwCtx(ctx, CategoryAlreadyExists,
			commonkeys.CategoryName, newName,
			commonkeys.UserID, strconv.FormatUint(cmd.UserID, 10),
		)
		return ErrCategoryAlreadyExists
	}

	return nil
}

func (s *Service) invalidateCategoryCaches(ctx context.Context, updatedCategory domain.Category) {
	err := s.CategoryCache.DeleteCategory(ctx, updatedCategory.ID, updatedCategory.UserID)
	if err != nil {
		s.Logger.WarnwCtx(ctx, WarnFailedToDeleteCategoryCache,
			commonkeys.CategoryID, updatedCategory.ID,
			commonkeys.UserID, updatedCategory.UserID,
			commonkeys.Error, err,
		)
	}

	err = s.CategoryCache.DeleteCategoryByName(ctx, updatedCategory.Name, updatedCategory.UserID)
	if err != nil {
		s.Logger.WarnwCtx(ctx, WarnFailedToDeleteCategoryByNameCache,
			commonkeys.CategoryName, updatedCategory.Name,
			commonkeys.UserID, updatedCategory.UserID,
			commonkeys.Error, err,
		)
	}

	err = s.CategoryCache.DeleteCategoryList(ctx, updatedCategory.UserID)
	if err != nil {
		s.Logger.WarnwCtx(ctx, WarnFailedToInvalidateCategoryListCache,
			commonkeys.UserID, updatedCategory.UserID,
			commonkeys.Error, err,
		)
	}
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
