package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lechitz/AionApi/internal/record/adapter/secondary/db/mapper"
	"github.com/lechitz/AionApi/internal/record/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/record/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// SearchRecords performs full-text search on records with optional filters using PostgreSQL tsvector.
func (r RecordRepository) SearchRecords(ctx context.Context, userID uint64, filters domain.SearchFilters) ([]domain.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, "record.repository.search", trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String("search_query", filters.Query),
		attribute.Int("limit", filters.Limit),
		attribute.Int("offset", filters.Offset),
	))
	defer span.End()

	var records []model.Record

	// Build SQL query with full-text search and ts_rank for relevance ordering
	// Note: category_id is obtained via JOIN with tags table, not directly from records
	query := `
		SELECT id, user_id, tag_id, description,
		       value, duration_seconds, event_time, recorded_at,
		       source, timezone, status,
		       created_at, updated_at, deleted_at,
		       ts_rank(search_vector, plainto_tsquery('portuguese', $1)) as rank
		FROM aion_api.records
		WHERE user_id = $2 
		  AND deleted_at IS NULL
		  AND search_vector @@ plainto_tsquery('portuguese', $1)
	`

	args := []interface{}{filters.Query, userID}
	argIndex := 3

	// Add category filter (requires JOIN with tags table)
	if len(filters.CategoryIDs) > 0 {
		// Need to modify query to include JOIN
		query = `
			SELECT r.id, r.user_id, r.tag_id, r.description,
			       r.value, r.duration_seconds, r.event_time, r.recorded_at,
			       r.source, r.timezone, r.status,
			       r.created_at, r.updated_at, r.deleted_at,
			       ts_rank(r.search_vector, plainto_tsquery('portuguese', $1)) as rank
			FROM aion_api.records r
			INNER JOIN aion_api.tags t ON r.tag_id = t.tag_id
			WHERE r.user_id = $2
			  AND r.deleted_at IS NULL
			  AND r.search_vector @@ plainto_tsquery('portuguese', $1)
			  AND t.category_id = ANY($3)
		`
		args = append(args, filters.CategoryIDs)
		argIndex++
	}

	// Add tag filter
	if len(filters.TagIDs) > 0 {
		query += fmt.Sprintf(" AND tag_id = ANY($%d)", argIndex)
		args = append(args, filters.TagIDs)
		argIndex++
	}

	// Add date range filters
	if filters.StartDate != nil {
		query += fmt.Sprintf(" AND event_time >= $%d", argIndex)
		args = append(args, filters.StartDate)
		argIndex++
	}

	if filters.EndDate != nil {
		query += fmt.Sprintf(" AND event_time <= $%d", argIndex)
		args = append(args, filters.EndDate)
		argIndex++
	}

	// Order by relevance (ts_rank) first, then by date
	query += " ORDER BY rank DESC, created_at DESC"

	// Add pagination
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filters.Limit, filters.Offset)

	// Execute query with Raw() and Scan()
	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&records).Error(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "search failed")
		r.logger.ErrorwCtx(ctx, "error searching records",
			commonkeys.Error, err.Error(),
			commonkeys.UserID, strconv.FormatUint(userID, 10),
		)
		return nil, fmt.Errorf("search records: %w", err)
	}

	span.SetAttributes(attribute.Int("results_count", len(records)))
	span.SetStatus(codes.Ok, "search completed successfully")

	r.logger.InfowCtx(ctx, "records searched successfully",
		commonkeys.UserID, strconv.FormatUint(userID, 10),
		"results_count", len(records),
		"query", filters.Query,
	)

	return mapper.RecordsFromDB(records), nil
}
