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

// ListAll retrieves all categories for the provided user ID.
func (h *controller) ListAll(ctx context.Context, userID uint64) ([]*model.Category, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, SpanListAll)
	defer span.End()

	span.SetAttributes(attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)))
	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound)
		return nil, errors.New(ErrUserIDNotFound)
	}

	all, err := h.CategoryService.ListAll(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "categories not found")
		h.Logger.ErrorwCtx(ctx, "categories not found",
			commonkeys.Error, err.Error(),
			commonkeys.UserID, strconv.FormatUint(userID, 10),
		)
		return nil, err
	}

	out := make([]*model.Category, len(all))
	for i, c := range all {
		out[i] = toModelOut(c)
	}

	span.SetAttributes(attribute.Int(commonkeys.CategoriesCount, len(out)))
	span.SetStatus(codes.Ok, StatusFetchedAll)
	return out, nil
}
