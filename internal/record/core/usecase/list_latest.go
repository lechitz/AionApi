package usecase

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lechitz/AionApi/internal/record/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ListLatest returns the N most recent records for the authenticated user, ordered by event_time DESC.
func (s *Service) ListLatest(ctx context.Context, userID uint64, limit int) ([]domain.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanListLatest)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanListLatest),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.Int("limit", limit),
	)

	if limit <= 0 || limit > 100 {
		limit = 10 // default limit for latest
	}

	span.AddEvent(EventRepositoryList)
	records, err := s.RecordRepository.ListLatest(ctx, userID, limit)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, FailedToListRecords)
		s.Logger.ErrorwCtx(ctx, FailedToListRecords,
			commonkeys.UserID, userID,
			"limit", limit,
			commonkeys.Error, err,
		)
		return nil, fmt.Errorf("%s: %w", FailedToListRecords, err)
	}

	span.AddEvent(EventSuccess)
	span.SetStatus(codes.Ok, StatusListedAll)
	s.Logger.InfowCtx(ctx, "latest records listed successfully",
		commonkeys.UserID, userID,
		"limit", limit,
		"count", len(records),
	)

	return records, nil
}
