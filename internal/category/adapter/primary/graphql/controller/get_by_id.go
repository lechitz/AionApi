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

// GetByID fetches a category by ID, adding tracing/logging and delegating to the input port.
func (h *controller) GetByID(ctx context.Context, categoryID, userID uint64) (*model.Category, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanGetByID)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanGetByID),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(categoryID, 10)),
	)

	// Basic guards (controller-level preconditions).
	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound)
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound, commonkeys.UserID, userID)
		return nil, errors.New(ErrUserIDNotFound)
	}
	if categoryID == 0 {
		span.SetStatus(codes.Error, ErrCategoryNotFound)
		h.Logger.ErrorwCtx(ctx, ErrCategoryNotFound, commonkeys.CategoryID, categoryID)
		return nil, errors.New(ErrCategoryNotFound)
	}

	// Delegate to the input port (use case).
	category, err := h.CategoryService.GetByID(ctx, categoryID, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrCategoryNotFound)
		h.Logger.ErrorwCtx(
			ctx,
			ErrCategoryNotFound,
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
