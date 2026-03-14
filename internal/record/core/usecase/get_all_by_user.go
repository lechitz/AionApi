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

// ListByUser returns records for a user with optional cursor parameters.
func (s *Service) ListByUser(ctx context.Context, userID uint64, limit int, afterEventTime *string, afterID *int64) ([]domain.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanListAll)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanListAll),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.Int("limit", limit),
	)

	// TODO: buscar primeiro no cache

	span.AddEvent(EventRepositoryList)
	records, err := s.RecordRepository.ListByUser(ctx, userID, limit, afterEventTime, afterID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, FailedToListRecords)
		s.Logger.ErrorwCtx(ctx, FailedToListRecords,
			commonkeys.UserID, userID,
			commonkeys.Error, err,
		)
		return nil, fmt.Errorf("%s: %w", FailedToListRecords, err)
	}

	span.AddEvent(EventSuccess)
	span.SetStatus(codes.Ok, StatusListedAll)
	s.Logger.InfowCtx(ctx, "records listed by user successfully",
		commonkeys.UserID, userID,
		"count", len(records),
	)

	return records, nil
}
