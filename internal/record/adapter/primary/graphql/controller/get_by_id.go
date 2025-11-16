// Package controller contains GraphQL-facing controllers for the Tag context.
package controller

import (
	"context"
	"errors"
	"strconv"

	gmodel "github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetByID fetches a Record by its ID for the authenticated user.
// It adds tracing/logging, applies basic guards, and delegates to the input port.
func (h *controller) GetByID(ctx context.Context, recordID, userID uint64) (*gmodel.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanGetByName)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanGetByName),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.RecordID, strconv.FormatUint(recordID, 10)),
	)

	// Controller-level preconditions.
	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound)
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound, commonkeys.UserID, userID)
		return nil, errors.New(ErrUserIDNotFound)
	}

	if recordID == 0 {
		span.SetStatus(codes.Error, ErrRecordNotFound)
		h.Logger.ErrorwCtx(ctx, ErrRecordNotFound, commonkeys.RecordID, recordID)
		return nil, errors.New(ErrRecordNotFound)
	}

	record, err := h.RecordService.GetByID(ctx, recordID, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrRecordNotFound)
		h.Logger.ErrorwCtx(
			ctx,
			ErrRecordNotFound,
			"error", err.Error(),
			commonkeys.UserID, userID,
			commonkeys.RecordID, recordID,
		)
		return nil, err
	}

	out := toModelOut(record)
	span.SetAttributes(attribute.String(commonkeys.RecordID, out.ID))
	span.SetStatus(codes.Ok, StatusFetched)
	return out, nil
}
