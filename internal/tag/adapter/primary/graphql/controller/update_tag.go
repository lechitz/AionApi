// Package controller contains GraphQL-facing controllers for the Tag context.
package controller

import (
	"context"
	"strconv"

	"github.com/lechitz/aion-api/internal/adapter/primary/graphql/model"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Update updates a tag using the provided GraphQL input, adding tracing/logging and delegating to the input port.
func (h *controller) Update(ctx context.Context, in model.UpdateTagInput, userID uint64) (*model.Tag, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanUpdate)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanUpdate),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.TagID, in.ID),
	)

	// Basic guards (controller-level preconditions).
	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound.Error())
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound.Error(), commonkeys.UserID, userID)
		return nil, ErrUserIDNotFound
	}

	cmd, err := toUpdateCommand(in, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrInvalidTagID.Error())
		h.Logger.ErrorwCtx(ctx, ErrInvalidTagID.Error(), commonkeys.Error, err.Error())
		return nil, ErrInvalidTagID
	}

	// Delegate to the input port (use case).
	tag, err := h.TagService.Update(ctx, cmd)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, MsgUpdateError)
		h.Logger.ErrorwCtx(
			ctx,
			MsgUpdateError,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, strconv.FormatUint(userID, 10),
			commonkeys.TagID, in.ID,
		)
		return nil, err
	}

	out := toModelOut(tag)
	span.SetAttributes(attribute.String(commonkeys.TagName, out.Name))
	span.SetStatus(codes.Ok, StatusUpdated)

	h.Logger.InfowCtx(
		ctx,
		MsgUpdated,
		commonkeys.UserID, strconv.FormatUint(userID, 10),
		commonkeys.TagID, out.ID,
		commonkeys.TagName, out.Name,
	)
	return out, nil
}
