// Package controller contains GraphQL-facing controllers for the Category context.
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

// ListAll retrieves all categories for a given user, adding tracing/logging and delegating to the input port.
func (h *controller) ListAll(ctx context.Context, userID uint64) ([]*model.Category, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanListAll)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanListAll),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	)

	// Basic guards (controller-level preconditions).
	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound)
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound, commonkeys.UserID, userID)
		return nil, errors.New(ErrUserIDNotFound)
	}

	// Delegate to the input port (use case).
	all, err := h.CategoryService.ListAll(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrCategoriesNotFound)
		h.Logger.ErrorwCtx(
			ctx,
			ErrCategoriesNotFound,
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
