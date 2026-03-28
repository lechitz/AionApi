package usecase

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lechitz/aion-api/internal/record/core/domain"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ListByCategory returns records filtered by category for the authenticated user.
// Records are retrieved via JOIN (records → tags → categories).
func (s *Service) ListByCategory(ctx context.Context, categoryID uint64, userID uint64, limit int) ([]domain.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, "record.list_by_category")
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, "list_by_category"),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(categoryID, 10)),
		attribute.Int("limit", limit),
	)

	if limit <= 0 || limit > 100 {
		limit = 50 // default limit
	}

	span.AddEvent(EventRepositoryList)
	records, err := s.RecordRepository.ListByCategory(ctx, categoryID, userID, limit, nil, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, FailedToListRecords)
		s.Logger.ErrorwCtx(ctx, FailedToListRecords,
			commonkeys.CategoryID, categoryID,
			commonkeys.UserID, userID,
			commonkeys.Error, err,
		)
		return nil, fmt.Errorf("%s: %w", FailedToListRecords, err)
	}

	span.AddEvent(EventSuccess)
	span.SetStatus(codes.Ok, StatusListedAll)
	s.Logger.InfowCtx(ctx, "records listed by category successfully",
		commonkeys.CategoryID, categoryID,
		commonkeys.UserID, userID,
		"count", len(records),
	)

	return records, nil
}
