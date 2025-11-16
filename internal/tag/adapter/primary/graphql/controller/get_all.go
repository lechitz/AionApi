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

// GetAll is the resolver for the tag field.
func (h *controller) GetAll(ctx context.Context, userID uint64) ([]*model.Tag, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanListAll)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanListAll),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	)

	// Controller-level preconditions.
	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound)
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound, commonkeys.UserID, userID)
		return nil, errors.New(ErrUserIDNotFound)
	}

	// Delegate to use case (input port).
	tags, err := h.TagService.GetAll(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrFailedToListTags)
		h.Logger.ErrorwCtx(
			ctx,
			ErrFailedToListTags,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, userID,
		)
		return nil, err
	}

	var out []*model.Tag
	for _, tag := range tags {
		out = append(out, toModelOut(tag))
	}

	span.SetStatus(codes.Ok, StatusFetched)
	return out, nil
}
