package controller

import (
	"context"
	"strconv"

	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// SoftDelete performs a soft deletion for a tag, delegating to the input port with tracing/logging.
func (h *controller) SoftDelete(ctx context.Context, tagID, userID uint64) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanSoftDelete, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.TagID, strconv.FormatUint(tagID, 10)),
		attribute.String(commonkeys.Operation, SpanSoftDelete),
	))
	defer span.End()

	if err := h.TagService.SoftDelete(ctx, tagID, userID); err != nil {
		span.SetStatus(codes.Error, MsgSoftDeleteError)
		h.Logger.ErrorwCtx(ctx, MsgSoftDeleteError, commonkeys.TagID, strconv.FormatUint(tagID, 10), commonkeys.Error, err.Error())
		return err
	}

	span.SetStatus(codes.Ok, StatusSoftDeleted)
	h.Logger.InfowCtx(ctx, MsgSoftDeleted, commonkeys.TagID, strconv.FormatUint(tagID, 10))
	return nil
}
