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

// GetByID retrieves a handler by its ID and userID from the database and returns it as a domain.Tag or an error if not found.
func (r TagRepository) GetByID(ctx context.Context, tagID, userID uint64) (domain.Tag, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanGetByNameRepo, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.TagID, strconv.FormatUint(tagID, 10)),
		attribute.String(commonkeys.Operation, OpGetByID),
	))
	defer span.End()

	var tagDB model.TagDB
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND tag_id = ?", userID, tagID).
		First(&tagDB).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.SetStatus(codes.Ok, ErrTagNotFoundMsg)
			return domain.Tag{}, nil
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, OpGetByID)
		r.logger.ErrorwCtx(ctx, ErrGetTagByIDMsg,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, strconv.FormatUint(userID, 10),
			commonkeys.TagID, strconv.FormatUint(tagID, 10),
		)
		return domain.Tag{}, fmt.Errorf("get tag by id: %w", err)
	}

	span.SetStatus(codes.Ok, StatusRetrievedByID)
	return mapper.TagFromDB(tagDB), nil
}
