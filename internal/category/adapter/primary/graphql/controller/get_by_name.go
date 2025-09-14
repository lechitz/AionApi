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

// GetByName fetches a category by name, adding tracing/logging and delegating to the input port.
func (h *controller) GetByName(ctx context.Context, categoryName string, userID uint64) (*model.Category, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanGetByName)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanGetByName),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.CategoryName, categoryName),
	)

	// Basic guards (controller-level preconditions).
	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound)
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound, commonkeys.UserID, userID)
		return nil, errors.New(ErrUserIDNotFound)
	}
	if categoryName == "" {
		span.SetStatus(codes.Error, ErrCategoryNotFound)
		h.Logger.ErrorwCtx(ctx, ErrCategoryNotFound, commonkeys.CategoryName, categoryName)
		return nil, errors.New(ErrCategoryNotFound)
	}

	// Delegate to the input port (use case).
	category, err := h.CategoryService.GetByName(ctx, categoryName, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrCategoryNotFound)
		h.Logger.ErrorwCtx(
			ctx,
			ErrCategoryNotFound,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, userID,
			commonkeys.CategoryName, categoryName,
		)
		return nil, err
	}

	out := toModelOut(category)
	span.SetAttributes(attribute.String(commonkeys.CategoryID, out.ID))
	span.SetStatus(codes.Ok, StatusFetched)
	return out, nil
}
