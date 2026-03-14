package controller

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetProjectedByID returns one derived record projection owned by the canonical API.
func (c *controller) GetProjectedByID(ctx context.Context, recordID, userID uint64) (*model.RecordProjection, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanGetProjectedByID)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanGetProjectedByID),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.RecordID, strconv.FormatUint(recordID, 10)),
	)

	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound.Error())
		c.Logger.ErrorwCtx(ctx, ErrUserIDNotFound.Error(), commonkeys.UserID, userID)
		return nil, ErrUserIDNotFound
	}
	if recordID == 0 {
		span.SetStatus(codes.Error, ErrInvalidRecordID.Error())
		c.Logger.ErrorwCtx(ctx, ErrInvalidRecordID.Error(), commonkeys.RecordID, recordID)
		return nil, ErrInvalidRecordID
	}

	projection, err := c.RecordService.GetProjectedByID(ctx, recordID, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, MsgGetProjectedByIDError)
		c.Logger.ErrorwCtx(ctx, MsgGetProjectedByIDError, commonkeys.Error, err.Error(), commonkeys.RecordID, recordID, commonkeys.UserID, userID)
		return nil, err
	}

	span.SetStatus(codes.Ok, StatusFetched)
	return toProjectedModelOut(projection), nil
}

// ListProjectedLatest returns the latest derived record projections for the user.
func (c *controller) ListProjectedLatest(ctx context.Context, userID uint64, limit int) ([]*model.RecordProjection, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanListProjectedLatest)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanListProjectedLatest),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.Int(AttrLimit, limit),
	)

	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound.Error())
		c.Logger.ErrorwCtx(ctx, ErrUserIDNotFound.Error(), commonkeys.UserID, userID)
		return nil, ErrUserIDNotFound
	}

	items, err := c.RecordService.ListProjectedLatest(ctx, userID, limit)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, MsgListProjectedLatestError)
		c.Logger.ErrorwCtx(ctx, MsgListProjectedLatestError, commonkeys.Error, err.Error(), commonkeys.UserID, userID)
		return nil, err
	}

	out := toProjectedModelOutSlice(items)
	span.SetAttributes(attribute.Int(AttrCount, len(out)))
	span.SetStatus(codes.Ok, StatusFetched)
	return out, nil
}

// ListProjectedPage returns cursor-based derived record projections for the user.
func (c *controller) ListProjectedPage(ctx context.Context, userID uint64, limit int, afterEventTime *string, afterID *int64) ([]*model.RecordProjection, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanListProjectedPage)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanListProjectedPage),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.Int(AttrLimit, limit),
	)

	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound.Error())
		c.Logger.ErrorwCtx(ctx, ErrUserIDNotFound.Error(), commonkeys.UserID, userID)
		return nil, ErrUserIDNotFound
	}

	items, err := c.RecordService.ListProjectedPage(ctx, userID, limit, afterEventTime, afterID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, MsgListProjectedPageError)
		c.Logger.ErrorwCtx(ctx, MsgListProjectedPageError, commonkeys.Error, err.Error(), commonkeys.UserID, userID)
		return nil, err
	}

	out := toProjectedModelOutSlice(items)
	span.SetAttributes(attribute.Int(AttrCount, len(out)))
	span.SetStatus(codes.Ok, StatusFetched)
	return out, nil
}
