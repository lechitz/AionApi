package handler

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/lechitz/AionApi/internal/shared/adapters/primary/graph/model"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Create starts a new tracing span, validates context, maps input to domain,
// invokes the creation use case, handles mapping of output, logs, and returns the result.
func (h *Handler) Create(ctx context.Context, in model.CreateCategoryInput, userID uint64) (*model.Category, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, SpanCreate)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.CategoryName, in.Name),
	)
	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound)
		return nil, errors.New(ErrUserIDNotFound)
	}

	domainCat := toDomainCreate(in, userID)
	created, err := h.CategoryService.Create(ctx, domainCat)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.Logger.ErrorwCtx(ctx, "error creating handler",
			commonkeys.Error, err.Error(),
			commonkeys.UserID, strconv.FormatUint(userID, 10),
			commonkeys.Input, fmt.Sprintf("%+v", in),
		)
		return nil, err
	}

	out := toModelOut(created)
	span.SetAttributes(attribute.String(commonkeys.CategoryID, out.ID))
	span.SetStatus(codes.Ok, StatusCreated)

	ip, _ := ctx.Value(ctxkeys.RequestIP).(string)
	ua, _ := ctx.Value(ctxkeys.RequestUserAgent).(string)
	h.Logger.InfowCtx(ctx, MsgCreated,
		commonkeys.UserID, strconv.FormatUint(userID, 10),
		commonkeys.CategoryID, out.ID,
		commonkeys.CategoryName, out.Name,
		tracingkeys.RequestIPKey, ip,
		tracingkeys.RequestUserAgentKey, ua,
	)
	return out, nil
}
