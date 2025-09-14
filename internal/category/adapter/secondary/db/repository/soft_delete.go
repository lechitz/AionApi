package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/category/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// SoftDelete updates the DeletedAt and UserUpdatedAt fields to mark a handler as soft-deleted based on handler ID and user ID.
func (c CategoryRepository) SoftDelete(ctx context.Context, categoryID uint64, userID uint64) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanSoftDeleteRepo, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(categoryID, 10)),
		attribute.String(commonkeys.Operation, OpSoftDelete),
	))
	defer span.End()

	fields := map[string]interface{}{
		commonkeys.CategoryDeletedAt: time.Now().UTC(),
		commonkeys.CategoryUpdatedAt: time.Now().UTC(),
	}

	if err := c.db.WithContext(ctx).
		Model(&model.CategoryDB{}).
		Where("category_id = ? AND user_id = ?", categoryID, userID).
		Updates(fields).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return err
	}

	span.SetStatus(codes.Ok, StatusSoftDeleted)
	return nil
}
