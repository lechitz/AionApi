package category

import (
	"context"
	"errors"
	"strconv"

	constants "github.com/lechitz/AionApi/internal/adapters/primary/graph/controllers/category/constants"
	"github.com/lechitz/AionApi/internal/adapters/primary/graph/model"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Update handles category update following the same orchestration: tracing, mapping,
// error handling, and delegating to CategoryService.Update.
func (h *Handler) Update(ctx context.Context, in model.UpdateCategoryInput, userID uint64) (*model.Category, error) {
	tracer := otel.Tracer(constants.TracerName)
	ctx, span := tracer.Start(ctx, constants.SpanUpdate)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.CategoryID, in.ID),
	)
	if userID == 0 {
		span.SetStatus(codes.Error, constants.ErrUserIDNotFound)
		return nil, errors.New(constants.ErrUserIDNotFound)
	}

	domainCat, err := toDomainUpdate(in, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, constants.ErrInvalidCategoryID)
		return nil, errors.New(constants.ErrInvalidCategoryID)
	}

	updated, err := h.CategoryService.Update(ctx, domainCat)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.Logger.ErrorwCtx(ctx, "error updating category",
			commonkeys.Error, err.Error(),
			commonkeys.UserID, strconv.FormatUint(userID, 10),
			commonkeys.CategoryID, in.ID,
		)
		return nil, err
	}

	out := toModelOut(updated)
	span.SetAttributes(attribute.String(commonkeys.CategoryName, out.Name))
	span.SetStatus(codes.Ok, constants.StatusUpdated)

	ip, _ := ctx.Value(ctxkeys.RequestIP).(string)
	ua, _ := ctx.Value(ctxkeys.RequestUserAgent).(string)
	h.Logger.InfowCtx(ctx, constants.MsgUpdated,
		commonkeys.UserID, strconv.FormatUint(userID, 10),
		commonkeys.CategoryID, out.ID,
		commonkeys.CategoryName, out.Name,
		tracingkeys.RequestIPKey, ip,
		tracingkeys.RequestUserAgentKey, ua,
	)
	return out, nil
}
