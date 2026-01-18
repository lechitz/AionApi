package repository

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/category/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
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

	q := c.db.WithContext(ctx).
		Model(&model.CategoryDB{}).
		Where("category_id = ? AND user_id = ?", categoryID, userID).
		Updates(fields)
	if err := q.Error(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return err
	}

	// If no rows were affected, the category either doesn't exist or doesn't belong to the user.
	if q.RowsAffected() == 0 {
		err := gorm.ErrRecordNotFound
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return errors.New(ErrCategoryNotFoundMsg)
	}

	span.SetStatus(codes.Ok, StatusSoftDeleted)
	return nil
}
