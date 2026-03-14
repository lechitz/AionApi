package controller

import (
	"context"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
	"github.com/lechitz/AionApi/internal/record/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// SearchRecords performs full-text search on records with optional filters.
func (c *controller) SearchRecords(ctx context.Context, filters model.SearchFilters, userID uint64) ([]*model.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanSearch)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, SpanSearch),
	)

	// Convert GraphQL filters to domain filters
	domainFilters := convertSearchFilters(filters)

	// Call service
	records, err := c.RecordService.SearchRecords(ctx, userID, domainFilters)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, MsgSearchError)
		c.Logger.ErrorwCtx(ctx, MsgSearchError,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, strconv.FormatUint(userID, 10),
		)
		return nil, err
	}

	// Convert to GraphQL models
	result := toModelOutSlice(records)

	span.SetAttributes(attribute.Int(AttrResultsCount, len(result)))
	span.SetStatus(codes.Ok, StatusSearchCompleted)

	c.Logger.InfowCtx(ctx, MsgSearched,
		commonkeys.UserID, strconv.FormatUint(userID, 10),
		AttrResultsCount, len(result),
	)

	return result, nil
}

// convertSearchFilters converts GraphQL SearchFilters to domain SearchFilters.
func convertSearchFilters(filters model.SearchFilters) domain.SearchFilters {
	domainFilters := domain.SearchFilters{
		Query:  filters.Query,
		Limit:  20, // default
		Offset: 0,
	}

	// Convert category IDs
	domainFilters.CategoryIDs = convertIDSlice(filters.CategoryIds)

	// Convert tag IDs
	domainFilters.TagIDs = convertIDSlice(filters.TagIds)

	// Parse dates
	if filters.StartDate != nil {
		if t, err := time.Parse(time.RFC3339, *filters.StartDate); err == nil {
			domainFilters.StartDate = &t
		}
	}

	if filters.EndDate != nil {
		if t, err := time.Parse(time.RFC3339, *filters.EndDate); err == nil {
			domainFilters.EndDate = &t
		}
	}

	// Set limit and offset
	if filters.Limit != nil {
		domainFilters.Limit = int(*filters.Limit)
	}

	if filters.Offset != nil {
		domainFilters.Offset = int(*filters.Offset)
	}

	return domainFilters
}

// convertIDSlice converts a slice of string IDs to uint64 IDs.
func convertIDSlice(stringIDs []string) []uint64 {
	if stringIDs == nil {
		return nil
	}

	ids := make([]uint64, 0, len(stringIDs))
	for _, id := range stringIDs {
		if parsed, err := strconv.ParseUint(id, 10, 64); err == nil {
			ids = append(ids, parsed)
		}
	}
	return ids
}
