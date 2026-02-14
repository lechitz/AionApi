package controller

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// SoftDelete deletes a single record (soft) owned by the user.
func (h *controller) SoftDelete(ctx context.Context, recordID, userID uint64) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanSoftDelete)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.RecordID, strconv.FormatUint(recordID, 10)),
	)

	if err := h.RecordService.Delete(ctx, recordID, userID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, MsgSoftDeleteError)
		return err
	}
	span.SetStatus(codes.Ok, StatusDeleted)
	return nil
}

// SoftDeleteAll soft deletes all records for the user.
func (h *controller) SoftDeleteAll(ctx context.Context, userID uint64) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanSoftDeleteAll)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	)

	if err := h.RecordService.DeleteAll(ctx, userID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, MsgSoftDeleteAllError)
		return err
	}
	span.SetStatus(codes.Ok, StatusDeletedAll)
	return nil
}
