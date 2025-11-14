// Package controller contains GraphQL-facing controllers for the Tag context.
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

// GetByID fetches a Tag by its ID for the authenticated user.
// It adds tracing/logging, applies basic guards, and delegates to the input port.
func (h *controller) GetByID(ctx context.Context, tagID, userID uint64) (*model.Tag, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanGetByName)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanGetByName),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.TagID, strconv.FormatUint(tagID, 10)),
	)

	// Controller-level preconditions.
	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound)
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound, commonkeys.UserID, userID)
		return nil, errors.New(ErrUserIDNotFound)
	}

	if tagID == 0 {
		span.SetStatus(codes.Error, ErrTagNotFound)
		h.Logger.ErrorwCtx(ctx, ErrTagNotFound, commonkeys.TagID, tagID)
		return nil, errors.New(ErrTagNotFound)
	}

	// Delegate to use case (input port).
	tag, err := h.TagService.GetByID(ctx, tagID, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrTagNotFound)
		h.Logger.ErrorwCtx(
			ctx,
			ErrTagNotFound,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, userID,
			commonkeys.TagID, tagID,
		)
		return nil, err
	}

	out := toModelOut(tag)
	span.SetAttributes(attribute.String(commonkeys.TagID, out.ID))
	span.SetStatus(codes.Ok, StatusFetched)
	return out, nil
}
