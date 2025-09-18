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

// Create is the resolver for the createTag field.
func (h *controller) Create(ctx context.Context, in model.CreateTagInput, userID uint64) (*model.Tag, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanCreate)
	defer span.End()

	categoryID, err := strconv.ParseUint(in.CategoryID, 10, 64)
	if err != nil {
		span.SetStatus(codes.Error, ErrInvalidCategoryID)
		h.Logger.ErrorwCtx(ctx, ErrInvalidCategoryID, commonkeys.CategoryID, in.CategoryID, commonkeys.Error, err.Error())
		return nil, errors.New(ErrInvalidCategoryID)
	}

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanCreate),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(categoryID, 10)),
		attribute.String(commonkeys.TagName, in.Name),
	)

	// Basic guards (controller-level preconditions).
	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound)
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound, commonkeys.UserID, userID)
		return nil, errors.New(ErrUserIDNotFound)
	}

	cmd := toCreateCommand(in, userID, categoryID)

	// Delegate to the input port (use case).
	domainOut, err := h.TagService.Create(ctx, cmd)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, MsgCreateError)
		h.Logger.ErrorwCtx(
			ctx,
			MsgCreateError,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, strconv.FormatUint(userID, 10),
			commonkeys.CategoryName, in.Name,
		)
		return nil, err
	}

	out := toModelOut(domainOut)
	span.SetAttributes(
		attribute.String(commonkeys.CategoryID, out.ID),
		attribute.String(commonkeys.CategoryName, out.Name),
	)
	span.SetStatus(codes.Ok, StatusCreated)

	h.Logger.InfowCtx(
		ctx,
		MsgCreated,
		commonkeys.UserID, strconv.FormatUint(userID, 10),
		commonkeys.TagName, out.Name,
		commonkeys.CategoryID, out.ID,
		commonkeys.CategoryName, out.Name,
	)
	return out, nil
}
