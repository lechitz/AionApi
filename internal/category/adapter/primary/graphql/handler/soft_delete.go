package handler

import (
	"context"
	"errors"
	"strconv"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// SoftDelete orchestrates soft deletion using similar patterns: span, validation,
// delegation, and logging.
func (h *Handler) SoftDelete(ctx context.Context, categoryID, userID uint64) error {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, SpanSoftDelete)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(categoryID, 10)),
	)
	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound)
		return errors.New(ErrUserIDNotFound)
	}

	if err := h.CategoryService.SoftDelete(ctx, categoryID, userID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.Logger.ErrorwCtx(ctx, "error soft deleting handler",
			commonkeys.Error, err.Error(),
			commonkeys.UserID, userID,
			commonkeys.CategoryID, categoryID,
		)
		return err
	}

	span.SetStatus(codes.Ok, StatusSoftDeleted)

	ip, _ := ctx.Value(ctxkeys.RequestIP).(string)
	ua, _ := ctx.Value(ctxkeys.RequestUserAgent).(string)
	h.Logger.InfowCtx(ctx, MsgSoftDeleted,
		commonkeys.UserID, strconv.FormatUint(userID, 10),
		commonkeys.CategoryID, categoryID,
		tracingkeys.RequestIPKey, ip,
		tracingkeys.RequestUserAgentKey, ua,
	)
	return nil
}
