package controller

import (
	"context"
	"errors"
	"strconv"
	"time"

	gmodel "github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ListByTag fetches all Records for a specific tag for the authenticated user.
//
//nolint:dupl // Similar to ListByCategory but with different entity - duplication improves clarity
func (h *controller) ListByTag(ctx context.Context, tagID, userID uint64, limit int) ([]*gmodel.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanListByTag)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanListByTag),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String("tag_id", strconv.FormatUint(tagID, 10)),
		attribute.Int("limit", limit),
	)

	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound.Error())
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound.Error(), commonkeys.UserID, userID)
		return nil, ErrUserIDNotFound
	}

	if tagID == 0 {
		span.SetStatus(codes.Error, "tag id cannot be zero")
		h.Logger.ErrorwCtx(ctx, "tag id cannot be zero", "tag_id", tagID)
		return nil, errors.New("tag id cannot be zero")
	}

	records, err := h.RecordService.ListByTag(ctx, tagID, userID, limit)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error listing records by tag")
		h.Logger.ErrorwCtx(ctx, "error listing records by tag", "error", err.Error(), "tag_id", tagID, commonkeys.UserID, userID)
		return nil, err
	}

	out := make([]*gmodel.Record, len(records))
	for i, rec := range records {
		out[i] = toModelOut(rec)
	}

	span.SetAttributes(attribute.Int("count", len(out)))
	span.SetStatus(codes.Ok, StatusFetched)
	return out, nil
}

// ListByCategory fetches all Records for a specific category for the authenticated user.
// Records are retrieved via JOIN (records → tags → categories).
//
//nolint:dupl // Similar to ListByTag but with different entity - duplication improves clarity
func (h *controller) ListByCategory(ctx context.Context, categoryID, userID uint64, limit int) ([]*gmodel.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, "record.list_by_category")
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, "list_by_category"),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String("category_id", strconv.FormatUint(categoryID, 10)),
		attribute.Int("limit", limit),
	)

	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound.Error())
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound.Error(), commonkeys.UserID, userID)
		return nil, ErrUserIDNotFound
	}

	if categoryID == 0 {
		span.SetStatus(codes.Error, "category id cannot be zero")
		h.Logger.ErrorwCtx(ctx, "category id cannot be zero", "category_id", categoryID)
		return nil, errors.New("category id cannot be zero")
	}

	records, err := h.RecordService.ListByCategory(ctx, categoryID, userID, limit)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error listing records by category")
		h.Logger.ErrorwCtx(ctx, "error listing records by category", "error", err.Error(), "category_id", categoryID, commonkeys.UserID, userID)
		return nil, err
	}

	out := make([]*gmodel.Record, len(records))
	for i, rec := range records {
		out[i] = toModelOut(rec)
	}

	span.SetAttributes(attribute.Int("count", len(out)))
	span.SetStatus(codes.Ok, StatusFetched)
	return out, nil
}

// ListByDay fetches all Records for a specific day for the authenticated user.
func (h *controller) ListByDay(ctx context.Context, userID uint64, dateStr string) ([]*gmodel.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanListByDay)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanListByDay),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String("date", dateStr),
	)

	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound.Error())
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound.Error(), commonkeys.UserID, userID)
		return nil, ErrUserIDNotFound
	}

	// Parse date string (expected format: YYYY-MM-DD or RFC3339)
	var date time.Time
	var err error
	if date, err = time.Parse("2006-01-02", dateStr); err != nil {
		// Try RFC3339
		if date, err = time.Parse(time.RFC3339, dateStr); err != nil {
			span.SetStatus(codes.Error, "invalid date format")
			h.Logger.ErrorwCtx(ctx, "invalid date format", "date", dateStr, "error", err.Error())
			return nil, errors.New("invalid date format, expected YYYY-MM-DD or RFC3339")
		}
	}

	records, err := h.RecordService.ListByDay(ctx, userID, date)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error listing records by day")
		h.Logger.ErrorwCtx(ctx, "error listing records by day", "error", err.Error(), "date", dateStr, commonkeys.UserID, userID)
		return nil, err
	}

	out := make([]*gmodel.Record, len(records))
	for i, rec := range records {
		out[i] = toModelOut(rec)
	}

	span.SetAttributes(attribute.Int("count", len(out)))
	span.SetStatus(codes.Ok, StatusFetched)
	return out, nil
}

// ListAllUntil fetches Records with event_time up to (and including) the given timestamp.
func (h *controller) ListAllUntil(ctx context.Context, userID uint64, untilStr string, limit int) ([]*gmodel.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanListAllUntil)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanListAllUntil),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String("until", untilStr),
		attribute.Int("limit", limit),
	)

	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound.Error())
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound.Error(), commonkeys.UserID, userID)
		return nil, ErrUserIDNotFound
	}

	until, err := time.Parse(time.RFC3339, untilStr)
	if err != nil {
		span.SetStatus(codes.Error, "invalid until timestamp")
		h.Logger.ErrorwCtx(ctx, "invalid until timestamp", "until", untilStr, "error", err.Error())
		return nil, errors.New("invalid until timestamp, expected RFC3339")
	}

	records, err := h.RecordService.ListAllUntil(ctx, userID, until, limit)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error listing records until")
		h.Logger.ErrorwCtx(ctx, "error listing records until", "error", err.Error(), "until", untilStr, commonkeys.UserID, userID)
		return nil, err
	}

	out := make([]*gmodel.Record, len(records))
	for i, rec := range records {
		out[i] = toModelOut(rec)
	}

	span.SetAttributes(attribute.Int("count", len(out)))
	span.SetStatus(codes.Ok, StatusFetched)
	return out, nil
}

// ListAllBetween fetches Records with event_time within the specified date range.
func (h *controller) ListAllBetween(ctx context.Context, userID uint64, startDateStr, endDateStr string, limit int) ([]*gmodel.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanListAllBetween)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanListAllBetween),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String("start_date", startDateStr),
		attribute.String("end_date", endDateStr),
		attribute.Int("limit", limit),
	)

	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound.Error())
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound.Error(), commonkeys.UserID, userID)
		return nil, ErrUserIDNotFound
	}

	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		span.SetStatus(codes.Error, "invalid start date")
		h.Logger.ErrorwCtx(ctx, "invalid start date", "start_date", startDateStr, "error", err.Error())
		return nil, errors.New("invalid start date, expected RFC3339")
	}

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		span.SetStatus(codes.Error, "invalid end date")
		h.Logger.ErrorwCtx(ctx, "invalid end date", "end_date", endDateStr, "error", err.Error())
		return nil, errors.New("invalid end date, expected RFC3339")
	}

	records, err := h.RecordService.ListAllBetween(ctx, userID, startDate, endDate, limit)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error listing records between dates")
		h.Logger.ErrorwCtx(
			ctx,
			"error listing records between dates",
			"error",
			err.Error(),
			"start_date",
			startDateStr,
			"end_date",
			endDateStr,
			commonkeys.UserID,
			userID,
		)
		return nil, err
	}

	out := make([]*gmodel.Record, len(records))
	for i, rec := range records {
		out[i] = toModelOut(rec)
	}

	span.SetAttributes(attribute.Int("count", len(out)))
	span.SetStatus(codes.Ok, StatusFetched)
	return out, nil
}

// ListByUser fetches records for the authenticated user with optional cursors.
func (h *controller) ListByUser(ctx context.Context, userID uint64, limit int, afterEventTime *string, afterID *int64) ([]*gmodel.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanListAll)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanListAll),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.Int("limit", limit),
	)

	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound.Error())
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound.Error(), commonkeys.UserID, userID)
		return nil, ErrUserIDNotFound
	}

	records, err := h.RecordService.ListByUser(ctx, userID, limit, afterEventTime, afterID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrFailedToListRecords.Error())
		h.Logger.ErrorwCtx(ctx, ErrFailedToListRecords.Error(), commonkeys.Error, err.Error())
		return nil, err
	}

	out := make([]*gmodel.Record, len(records))
	for i, rec := range records {
		out[i] = toModelOut(rec)
	}

	span.SetAttributes(attribute.Int("count", len(out)))
	span.SetStatus(codes.Ok, StatusFetched)
	return out, nil
}

// ListLatest fetches the N most recent records for the authenticated user.
func (h *controller) ListLatest(ctx context.Context, userID uint64, limit int) ([]*gmodel.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanListLatest)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanListLatest),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.Int("limit", limit),
	)

	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound.Error())
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound.Error(), commonkeys.UserID, userID)
		return nil, ErrUserIDNotFound
	}

	if limit <= 0 || limit > 100 {
		limit = 10 // default limit for latest
	}

	records, err := h.RecordService.ListLatest(ctx, userID, limit)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error listing latest records")
		h.Logger.ErrorwCtx(ctx, "error listing latest records", "error", err.Error(), commonkeys.UserID, userID)
		return nil, err
	}

	out := make([]*gmodel.Record, len(records))
	for i, rec := range records {
		out[i] = toModelOut(rec)
	}

	span.SetAttributes(attribute.Int("count", len(out)))
	span.SetStatus(codes.Ok, StatusFetched)
	return out, nil
}
