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

// GetByName fetches a Tag by its name for the authenticated user.
// It adds tracing/logging, applies basic guards, and delegates to the input port.
func (h *controller) GetByName(ctx context.Context, tagName string, userID uint64) (*model.Tag, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanGetByName)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanGetByName),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.TagName, tagName),
	)

	// Controller-level preconditions.
	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound)
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound, commonkeys.UserID, userID)
		return nil, errors.New(ErrUserIDNotFound)
	}

	if tagName == "" {
		span.SetStatus(codes.Error, ErrTagNotFound)
		h.Logger.ErrorwCtx(ctx, ErrTagNotFound, commonkeys.CategoryName, tagName)
		return nil, errors.New(ErrTagNotFound)
	}

	// Delegate to use case (input port).
	tag, err := h.TagService.GetByName(ctx, tagName, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrTagNotFound)
		h.Logger.ErrorwCtx(
			ctx,
			ErrTagNotFound,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, userID,
			commonkeys.TagName, tagName,
		)
		return nil, err
	}

	out := toModelOut(tag)
	span.SetAttributes(attribute.String(commonkeys.TagName, out.Name))
	span.SetStatus(codes.Ok, StatusFetched)
	return out, nil
}
