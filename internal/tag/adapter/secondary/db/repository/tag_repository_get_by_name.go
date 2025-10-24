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

// GetByName retrieves a handler by its name and user ID from the database and returns it as a domain.Tag or an error if not found.
func (r TagRepository) GetByName(ctx context.Context, tagName string, userID uint64) (domain.Tag, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanGetByNameRepo, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.TagName, tagName),
		attribute.String(commonkeys.Operation, OpGetByName),
	))
	defer span.End()

	var tagDB model.TagDB
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND name = ?", userID, tagName).
		First(&tagDB).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.SetStatus(codes.Ok, ErrTagNotFoundMsg)
			return domain.Tag{}, nil
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, OpGetByName)
		r.logger.ErrorwCtx(ctx, ErrGetTagByIDMsg,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, strconv.FormatUint(userID, 10),
			commonkeys.TagName, tagName,
		)
		return domain.Tag{}, fmt.Errorf("get tag by name: %w", err)
	}

	span.SetStatus(codes.Ok, StatusRetrievedByName)
	return mapper.TagFromDB(tagDB), nil
}
