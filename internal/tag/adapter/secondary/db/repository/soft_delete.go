package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"github.com/lechitz/aion-api/internal/tag/adapter/secondary/db/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// SoftDelete marks a tag as deleted for a given user.
func (r TagRepository) SoftDelete(ctx context.Context, tagID, userID uint64) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanSoftDeleteRepo, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.TagID, strconv.FormatUint(tagID, 10)),
		attribute.String(commonkeys.Operation, OpSoftDelete),
	))
	defer span.End()

	fields := map[string]interface{}{
		commonkeys.TagDeletedAt: time.Now().UTC(),
		commonkeys.TagUpdatedAt: time.Now().UTC(),
	}

	if err := r.db.WithContext(ctx).
		Model(&model.TagDB{}).
		Where("tag_id = ? AND user_id = ?", tagID, userID).
		Updates(fields).Error(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return err
	}

	span.SetStatus(codes.Ok, StatusSoftDeleted)
	return nil
}
