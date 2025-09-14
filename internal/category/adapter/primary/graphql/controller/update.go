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

// Update updates a category using the provided GraphQL input, adding tracing/logging and delegating to the input port.
func (h *controller) Update(ctx context.Context, in model.UpdateCategoryInput, userID uint64) (*model.Category, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanUpdate)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanUpdate),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.CategoryID, in.ID),
	)

	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound)
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound, commonkeys.UserID, userID)
		return nil, errors.New(ErrUserIDNotFound)
	}

	cmd, err := toUpdateCommand(in, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrInvalidCategoryID)
		h.Logger.ErrorwCtx(ctx, ErrInvalidCategoryID, commonkeys.Error, err.Error())
		return nil, errors.New(ErrInvalidCategoryID)
	}

	category, err := h.CategoryService.Update(ctx, cmd)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, MsgUpdateError)
		h.Logger.ErrorwCtx(
			ctx,
			MsgUpdateError,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, strconv.FormatUint(userID, 10),
			commonkeys.CategoryID, in.ID,
		)
		return nil, err
	}

	out := toModelOut(category)
	span.SetAttributes(attribute.String(commonkeys.CategoryName, out.Name))
	span.SetStatus(codes.Ok, StatusUpdated)

	h.Logger.InfowCtx(
		ctx,
		MsgUpdated,
		commonkeys.UserID, strconv.FormatUint(userID, 10),
		commonkeys.CategoryID, out.ID,
		commonkeys.CategoryName, out.Name,
	)
	return out, nil
}
