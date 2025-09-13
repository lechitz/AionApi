package controller

import (
	"context"
	"errors"
	"strconv"

	"github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetByID fetches a handler by ID.
func (h *controller) GetByID(ctx context.Context, categoryID uint64, userID uint64) (*model.Category, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, SpanGetByID)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(categoryID, 10)),
	)

	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound)
		return nil, errors.New(ErrUserIDNotFound)
	}

	if categoryID == 0 {
		span.SetStatus(codes.Error, ErrCategoryNotFound)
		return nil, errors.New(ErrCategoryNotFound)
	}

	category, err := h.CategoryService.GetByID(ctx, categoryID, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrCategoryNotFound)
		h.Logger.ErrorwCtx(ctx, ErrCategoryNotFound,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, userID,
			commonkeys.CategoryID, categoryID,
		)
		return nil, err
	}

	out := toModelOut(category)
	span.SetAttributes(attribute.String(commonkeys.CategoryName, out.Name))
	span.SetStatus(codes.Ok, StatusFetched)
	return out, nil
}
