package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/adapters/secondary/db/category/model"
	"github.com/lechitz/AionApi/internal/core/category/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// SoftDelete updates the DeletedAt and UserUpdatedAt fields to mark a handler as soft-deleted based on handler ID and user ID.
func (c CategoryRepository) SoftDelete(ctx context.Context, category domain.Category) error {
	tr := otel.Tracer("CategoryRepository")
	ctx, span := tr.Start(ctx, "SoftDelete", trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(category.UserID, 10)),
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(category.ID, 10)),
		attribute.String("operation", "soft_delete"),
	))
	defer span.End()

	fields := map[string]interface{}{
		commonkeys.CategoryDeletedAt: time.Now().UTC(),
		commonkeys.CategoryUpdatedAt: time.Now().UTC(),
	}

	if err := c.db.WithContext(ctx).
		Model(&model.CategoryDB{}).
		Where("category_id = ? AND user_id = ?", category.ID, category.UserID).
		Updates(fields).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return err
	}

	span.SetStatus(codes.Ok, "handler soft deleted successfully")
	return nil
}
