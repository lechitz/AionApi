package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/tag/adapter/secondary/db/mapper"
	"github.com/lechitz/AionApi/internal/tag/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/tag/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

// GetAll retrieves all tags for a given user from the database.
func (r TagRepository) GetAll(ctx context.Context, userID uint64) ([]domain.Tag, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanGetAllRepo, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, OpGetAll),
	))
	defer span.End()

	var tagsDB []model.TagDB
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&tagsDB).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.SetStatus(codes.Ok, ErrNoTagsFoundMsg)
			return []domain.Tag{}, nil
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, OpGetAll)
		r.logger.ErrorwCtx(ctx, ErrGetAllTagsMsg,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, strconv.FormatUint(userID, 10),
		)
		return []domain.Tag{}, fmt.Errorf("get all tags: %w", err)
	}

	tags := mapper.TagsFromDB(tagsDB)

	span.SetStatus(codes.Ok, StatusRetrievedAll)
	return tags, nil
}
