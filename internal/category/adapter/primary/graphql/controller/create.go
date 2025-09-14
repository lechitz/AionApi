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

// Create is the GraphQL-facing entrypoint for creating a new category.
func (h *controller) Create(ctx context.Context, in model.CreateCategoryInput, userID uint64) (*model.Category, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanCreate)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanCreate),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	)

	// Basic guards (controller-level preconditions).
	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound)
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound, commonkeys.UserID, userID)
		return nil, errors.New(ErrUserIDNotFound)
	}

	cmd := toCreateCommand(in, userID)

	// Delegate to the input port (use case).
	domainOut, err := h.CategoryService.Create(ctx, cmd)
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
		commonkeys.CategoryID, out.ID,
		commonkeys.CategoryName, out.Name,
	)
	return out, nil
}
