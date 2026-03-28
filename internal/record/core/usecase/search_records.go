package usecase

import (
	"context"
	"strconv"

	"github.com/lechitz/aion-api/internal/record/core/domain"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// SearchRecords searches for records using full-text search with optional filters.
func (s *Service) SearchRecords(ctx context.Context, userID uint64, filters domain.SearchFilters) ([]domain.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, "record.usecase.search")
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String("search_query", filters.Query),
		attribute.Int("limit", filters.Limit),
	)

	s.Logger.InfowCtx(ctx, "searching records",
		commonkeys.UserID, userID,
		"query", filters.Query,
		"limit", filters.Limit,
	)

	// Set default limit if not provided
	if filters.Limit == 0 {
		filters.Limit = 20
	}

	// Validate limit
	if filters.Limit > 100 {
		filters.Limit = 100
	}

	records, err := s.RecordRepository.SearchRecords(ctx, userID, filters)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "search failed")
		s.Logger.ErrorwCtx(ctx, "error searching records",
			commonkeys.Error, err.Error(),
			commonkeys.UserID, userID,
		)
		return nil, err
	}

	span.SetAttributes(attribute.Int("results_count", len(records)))
	span.SetStatus(codes.Ok, "search completed successfully")

	s.Logger.InfowCtx(ctx, "records searched successfully",
		commonkeys.UserID, userID,
		"results_count", len(records),
	)

	return records, nil
}
