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

// GetByName retrieves a category by its name.
func (h *Handler) GetByName(ctx context.Context, categoryName string, userID uint64) (*model.Category, error) {
	tracer := otel.Tracer(constants.TracerName)
	ctx, span := tracer.Start(ctx, constants.SpanGetByName)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.CategoryName, categoryName),
	)
	if userID == 0 {
		span.SetStatus(codes.Error, constants.ErrUserIDNotFound)
		return nil, errors.New(constants.ErrUserIDNotFound)
	}

	category, err := h.CategoryService.GetByName(ctx, categoryName, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "category by name not found")
		h.Logger.ErrorwCtx(ctx, "category by name not found",
			commonkeys.Error, err.Error(),
			commonkeys.UserID, userID,
			commonkeys.CategoryName, categoryName,
		)
		return nil, err
	}
	if category.Name == "" {
		span.SetStatus(codes.Ok, "category not found")
		return nil, errors.New("category not found")
	}

	out := toModelOut(category)
	span.SetAttributes(attribute.String(commonkeys.CategoryID, out.ID))
	span.SetStatus(codes.Ok, constants.StatusFetched)
	return out, nil
}
