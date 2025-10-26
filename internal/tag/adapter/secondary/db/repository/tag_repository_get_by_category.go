package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/tag/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/tag/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// GetByCategoryID retrieves all tags associated with a specific category for a given user.
func (r TagRepository) GetByCategoryID(ctx context.Context, categoryID uint64, userID uint64) ([]domain.Tag, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanGetByCategoryRepo, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(categoryID, 10)),
		attribute.String(commonkeys.Operation, OpGetByCategory),
	))
	defer span.End()

	var tagsDB []model.TagDB
	err := r.db.WithContext(ctx).
		Joins("JOIN category_tags ON category_tags.tag_id = tags.id").
		Where("category_tags.category_id = ? AND tags.user_id = ?", categoryID, userID).
		Find(&tagsDB).Error
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, OpGetByCategory)
		r.logger.ErrorwCtx(ctx, ErrGetTagsByCategoryMsg,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, strconv.FormatUint(userID, 10),
			commonkeys.CategoryID, strconv.FormatUint(categoryID, 10),
		)
		return []domain.Tag{}, fmt.Errorf("get tags by category id: %w", err)
	}

	var tags []domain.Tag
	for _, tagDB := range tagsDB {
		tags = append(tags, domain.Tag{
			ID:     tagDB.ID,
			UserID: tagDB.UserID,
			Name:   tagDB.Name,
		})
	}

	span.SetStatus(codes.Ok, StatusRetrievedByCategory)
	return tags, nil
}
