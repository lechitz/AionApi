package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"github.com/lechitz/aion-api/internal/tag/adapter/secondary/db/mapper"
	"github.com/lechitz/aion-api/internal/tag/adapter/secondary/db/model"
	"github.com/lechitz/aion-api/internal/tag/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// UpdateTag updates a tag in the database based on its ID and user ID, updating only fields specified in the updateFields map.
func (r TagRepository) UpdateTag(ctx context.Context, tagID uint64, userID uint64, updateFields map[string]interface{}) (domain.Tag, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanUpdateRepo, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.TagID, strconv.FormatUint(tagID, 10)),
		attribute.String(commonkeys.Operation, OpUpdate),
	))
	defer span.End()

	// Prevent overwriting created_at
	delete(updateFields, commonkeys.TagCreatedAt)

	var tagDB model.TagDB
	if err := r.db.WithContext(ctx).
		Model(&tagDB).
		Where("tag_id = ? AND user_id = ?", tagID, userID).
		Updates(updateFields).Error(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		r.logger.ErrorwCtx(ctx, "error updating tag", commonkeys.Error, err.Error())
		return domain.Tag{}, fmt.Errorf("update tag: %w", err)
	}

	// Fetch the updated tag
	if err := r.db.WithContext(ctx).
		Where("tag_id = ? AND user_id = ?", tagID, userID).
		First(&tagDB).Error(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		r.logger.ErrorwCtx(ctx, "error fetching updated tag", commonkeys.Error, err.Error())
		return domain.Tag{}, fmt.Errorf("fetch updated tag: %w", err)
	}

	span.SetStatus(codes.Ok, StatusUpdated)
	return mapper.TagFromDB(tagDB), nil
}
