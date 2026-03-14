package controller

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetByCategoryID fetches Tags by Category ID for the authenticated user.
func (h *controller) GetByCategoryID(ctx context.Context, categoryID, userID uint64) ([]*model.Tag, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanGetByCategory)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanGetByCategory),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(categoryID, 10)),
	)

	// Controller-level preconditions.
	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound.Error())
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound.Error(), commonkeys.UserID, userID)
		return nil, ErrUserIDNotFound
	}

	if categoryID == 0 {
		span.SetStatus(codes.Error, ErrCategoryNotFound.Error())
		h.Logger.ErrorwCtx(ctx, ErrCategoryNotFound.Error(), commonkeys.CategoryID, categoryID)
		return nil, ErrCategoryNotFound
	}

	// Delegate to use case (input port).
	tags, err := h.TagService.GetByCategoryID(ctx, categoryID, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrFailedToListTags.Error())
		h.Logger.ErrorwCtx(
			ctx,
			ErrFailedToListTags.Error(),
			commonkeys.Error, err.Error(),
			commonkeys.UserID, userID,
			commonkeys.CategoryID, categoryID,
		)
		return nil, err
	}

	var out []*model.Tag
	for _, tag := range tags {
		out = append(out, toModelOut(tag))
	}

	span.SetStatus(codes.Ok, StatusFetched)
	return out, nil
}
