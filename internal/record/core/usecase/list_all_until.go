package usecase

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/lechitz/aion-api/internal/record/core/domain"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ListAllUntil returns records with event_time up to (and including) the given timestamp.
func (s *Service) ListAllUntil(ctx context.Context, userID uint64, until time.Time, limit int) ([]domain.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanListAllUntil)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanListAllUntil),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String("until", until.Format("2006-01-02")),
		attribute.Int("limit", limit),
	)

	if limit <= 0 || limit > 100 {
		limit = 50 // default limit
	}

	span.AddEvent(EventRepositoryList)
	records, err := s.RecordRepository.ListAllUntil(ctx, userID, until, limit)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, FailedToListRecords)
		s.Logger.ErrorwCtx(ctx, FailedToListRecords,
			commonkeys.UserID, userID,
			"until", until.Format(time.RFC3339),
			commonkeys.Error, err,
		)
		return nil, fmt.Errorf("%s: %w", FailedToListRecords, err)
	}

	span.AddEvent(EventSuccess)
	span.SetStatus(codes.Ok, StatusListedAll)
	s.Logger.InfowCtx(ctx, "records listed until timestamp successfully",
		commonkeys.UserID, userID,
		"until", until.Format(time.RFC3339),
		"count", len(records),
	)

	return records, nil
}
