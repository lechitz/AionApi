package repository

import (
	"context"
	"errors"
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

		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return domain.Tag{}, err
	}

	span.SetStatus(codes.Ok, StatusRetrievedByName)
	return mapper.TagFromDB(tagDB), nil
}
