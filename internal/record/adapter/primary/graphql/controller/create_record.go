package controller

import (
	"context"
	"strconv"

	gmodel "github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Create is the resolver for creating a Record via GraphQL.
func (h *controller) Create(ctx context.Context, in gmodel.CreateRecordInput, userID uint64) (*gmodel.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanCreate)
	defer span.End()

	tagID, err := strconv.ParseUint(in.TagID, 10, 64)
	if err != nil {
		span.SetStatus(codes.Error, ErrInvalidRecordID.Error())
		h.Logger.ErrorwCtx(ctx, ErrInvalidRecordID.Error(), "tag_id", in.TagID, "error", err.Error())
		return nil, ErrInvalidRecordID
	}

	span.SetAttributes(
		attribute.String("operation", SpanCreate),
		attribute.String("user_id", strconv.FormatUint(userID, 10)),
		attribute.String("tag_id", strconv.FormatUint(tagID, 10)),
	)

	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound.Error())
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound.Error(), "user_id", userID)
		return nil, ErrUserIDNotFound
	}

	cmd := toCreateCommand(in, userID)

	// Delegate to the input port (use case).
	domainOut, err := h.RecordService.Create(ctx, cmd)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, MsgCreateError)
		h.Logger.ErrorwCtx(
			ctx,
			MsgCreateError,
			"error", err.Error(),
			"user_id", strconv.FormatUint(userID, 10),
			"tag_id", in.TagID,
		)
		return nil, err
	}

	out := toModelOut(domainOut)
	span.SetAttributes(
		attribute.String("record_id", out.ID),
		attribute.String("tag_id", out.TagID),
	)
	span.SetStatus(codes.Ok, StatusCreated)

	h.Logger.InfowCtx(
		ctx,
		MsgCreated,
		"user_id", strconv.FormatUint(userID, 10),
		"record_id", out.ID,
		"tag_id", out.TagID,
	)
	return out, nil
}
