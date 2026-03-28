// Package controller contains GraphQL-facing controllers for the Category context.
package controller

import (
	"context"
	"strconv"

	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// SoftDelete performs a soft deletion for a category, adding tracing/logging and delegating to the input port.
func (h *controller) SoftDelete(ctx context.Context, categoryID, userID uint64) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanSoftDelete)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(categoryID, 10)),
	)

	// Basic guards (controller-level preconditions).
	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound.Error())
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound.Error(), commonkeys.UserID, userID)
		return ErrUserIDNotFound
	}
	if categoryID == 0 {
		span.SetStatus(codes.Error, ErrCategoryIDNotFound.Error())
		h.Logger.ErrorwCtx(ctx, ErrCategoryIDNotFound.Error(), commonkeys.CategoryID, categoryID)
		return ErrCategoryIDNotFound
	}

	if err := h.CategoryService.SoftDelete(ctx, categoryID, userID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, MsgSoftDeleteError)
		h.Logger.ErrorwCtx(
			ctx,
			MsgSoftDeleteError,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, userID,
			commonkeys.CategoryID, categoryID,
		)
		return err
	}

	span.SetStatus(codes.Ok, StatusSoftDeleted)

	h.Logger.InfowCtx(
		ctx,
		MsgSoftDeleted,
		commonkeys.UserID, strconv.FormatUint(userID, 10),
		commonkeys.CategoryID, categoryID,
	)
	return nil
}
