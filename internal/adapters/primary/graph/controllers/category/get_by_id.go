package category

import (
	"context"
	"errors"
	"strconv"

	constants "github.com/lechitz/AionApi/internal/adapters/primary/graph/controllers/category/constants"
	"github.com/lechitz/AionApi/internal/adapters/primary/graph/model"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetByID fetches a category by ID.
func (h *Handler) GetByID(ctx context.Context, categoryID uint64, userID uint64) (*model.Category, error) {
	tracer := otel.Tracer(constants.TracerName)
	ctx, span := tracer.Start(ctx, constants.SpanGetByID)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(categoryID, 10)),
	)

	if userID == 0 {
		span.SetStatus(codes.Error, constants.ErrUserIDNotFound)
		return nil, errors.New(constants.ErrUserIDNotFound)
	}

	if categoryID == 0 {
		span.SetStatus(codes.Error, constants.ErrCategoryNotFound)
		return nil, errors.New(constants.ErrCategoryNotFound)
	}

	category, err := h.CategoryService.GetByID(ctx, categoryID, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, constants.ErrCategoryNotFound)
		h.Logger.ErrorwCtx(ctx, constants.ErrCategoryNotFound,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, userID,
			commonkeys.CategoryID, categoryID,
		)
		return nil, err
	}

	out := toModelOut(category)
	span.SetAttributes(attribute.String(commonkeys.CategoryName, out.Name))
	span.SetStatus(codes.Ok, constants.StatusFetched)
	return out, nil
}
