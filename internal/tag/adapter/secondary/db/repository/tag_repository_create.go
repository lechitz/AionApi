package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/tag/adapter/secondary/db/mapper"
	"github.com/lechitz/AionApi/internal/tag/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Create inserts a new tag and returns it, or an error on failure.
func (r TagRepository) Create(ctx context.Context, tag domain.Tag) (domain.Tag, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanCreateRepo, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(tag.UserID, 10)),
		attribute.String(commonkeys.TagName, tag.Name),
		attribute.String(commonkeys.Operation, OpCreate),
	))
	defer span.End()

	row := mapper.TagToDB(tag)
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, OpCreate)
		r.logger.ErrorwCtx(ctx, ErrCreateTagMsg, commonkeys.Error, err.Error())
		return domain.Tag{}, fmt.Errorf("insert tag: %w", err)
	}
	return mapper.TagFromDB(row), nil
}
