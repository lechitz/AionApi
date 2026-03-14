package usecase

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/tag/core/domain"
	"github.com/lechitz/AionApi/internal/tag/core/ports/input"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Update updates an existing Tag in the system using the provided UpdateTagCommand.
// It builds the update fields, calls the repository, and logs/traces the operation.
func (s *Service) Update(ctx context.Context, cmd input.UpdateTagCommand) (domain.Tag, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanUpdateTag)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanUpdateTag),
		attribute.String(commonkeys.TagID, strconv.FormatUint(cmd.ID, 10)),
		attribute.String(commonkeys.UserID, strconv.FormatUint(cmd.UserID, 10)),
	)

	if cmd.Icon != nil {
		icon := normalizeTagIcon(cmd.Icon)
		if !isSingleEmoji(icon) {
			span.SetStatus(codes.Error, ErrToValidateTag)
			return domain.Tag{}, ErrTagIconInvalid
		}
		cmd.Icon = &icon
	}

	fieldsToUpdate := extractUpdateFields(cmd)

	span.AddEvent(EventRepositoryUpdate)
	updatedTag, err := s.TagRepository.UpdateTag(ctx, cmd.ID, cmd.UserID, fieldsToUpdate)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, FailedToUpdateTag)
		s.Logger.ErrorwCtx(
			ctx,
			FailedToUpdateTag,
			commonkeys.TagID, strconv.FormatUint(cmd.ID, 10),
			commonkeys.Error, err,
		)
		return domain.Tag{}, err
	}

	span.AddEvent(EventInvalidateCache)
	err = s.TagCache.DeleteTag(ctx, updatedTag.ID, updatedTag.UserID)
	if err != nil {
		s.Logger.WarnwCtx(ctx, "failed to delete tag cache after update",
			commonkeys.TagID, updatedTag.ID,
			commonkeys.UserID, updatedTag.UserID,
			commonkeys.Error, err,
		)
	}

	err = s.TagCache.DeleteTagByName(ctx, updatedTag.Name, updatedTag.UserID)
	if err != nil {
		s.Logger.WarnwCtx(ctx, "failed to delete tag-by-name cache after update",
			commonkeys.TagName, updatedTag.Name,
			commonkeys.UserID, updatedTag.UserID,
			commonkeys.Error, err,
		)
	}

	err = s.TagCache.DeleteTagList(ctx, updatedTag.UserID)
	if err != nil {
		s.Logger.WarnwCtx(ctx, "failed to invalidate tag list cache after updating tag",
			commonkeys.UserID, updatedTag.UserID,
			commonkeys.Error, err,
		)
	}

	err = s.TagCache.DeleteTagsByCategory(ctx, updatedTag.CategoryID, updatedTag.UserID)
	if err != nil {
		s.Logger.WarnwCtx(ctx, "failed to invalidate tag list by category cache after updating tag",
			commonkeys.CategoryID, updatedTag.CategoryID,
			commonkeys.UserID, updatedTag.UserID,
			commonkeys.Error, err,
		)
	}

	span.AddEvent(EventSuccess)
	span.SetStatus(codes.Ok, StatusUpdated)
	s.Logger.InfowCtx(
		ctx,
		SuccessfullyUpdatedTag,
		commonkeys.TagID, strconv.FormatUint(updatedTag.ID, 10),
	)

	return updatedTag, nil
}

// extractUpdateFields builds a map with only the non-nil/non-empty fields from UpdateTagCommand.
func extractUpdateFields(cmd input.UpdateTagCommand) map[string]interface{} {
	updateFields := make(map[string]interface{})

	if cmd.Name != nil && *cmd.Name != "" {
		updateFields[commonkeys.TagName] = *cmd.Name
	}
	if cmd.Description != nil && *cmd.Description != "" {
		updateFields[commonkeys.TagDescription] = *cmd.Description
	}
	if cmd.CategoryID != nil {
		updateFields[commonkeys.CategoryID] = *cmd.CategoryID
	}
	if cmd.Icon != nil && *cmd.Icon != "" {
		updateFields[commonkeys.TagIcon] = *cmd.Icon
	}

	return updateFields
}
