package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/user/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// GetUserStats retrieves aggregated statistics for a user using SQL queries with Raw().
func (r UserRepository) GetUserStats(ctx context.Context, userID uint64) (domain.UserStats, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, "user.repository.get_stats", trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	))
	defer span.End()

	stats := domain.UserStats{}
	now := time.Now()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// Total records
	r.db.WithContext(ctx).Raw("SELECT COUNT(*) FROM aion_api.records WHERE user_id = ? AND deleted_at IS NULL", userID).Scan(&stats.TotalRecords)

	// Total categories
	r.db.WithContext(ctx).Raw("SELECT COUNT(*) FROM aion_api.categories WHERE user_id = ? AND deleted_at IS NULL", userID).Scan(&stats.TotalCategories)

	// Total tags
	r.db.WithContext(ctx).Raw("SELECT COUNT(*) FROM aion_api.tags WHERE user_id = ? AND deleted_at IS NULL", userID).Scan(&stats.TotalTags)

	// Records this week
	r.db.WithContext(ctx).
		Raw("SELECT COUNT(*) FROM aion_api.records WHERE user_id = ? AND created_at >= ? AND deleted_at IS NULL", userID, startOfWeek).
		Scan(&stats.RecordsThisWeek)

	// Records this month
	r.db.WithContext(ctx).
		Raw("SELECT COUNT(*) FROM aion_api.records WHERE user_id = ? AND created_at >= ? AND deleted_at IS NULL", userID, startOfMonth).
		Scan(&stats.RecordsThisMonth)

	// Most used category (with JOIN and GROUP BY)
	var categoryResult struct {
		CategoryID uint64
		Name       string
		Count      int
	}
	if err := r.db.WithContext(ctx).Raw(`
		SELECT r.category_id, c.name, COUNT(*) as count
		FROM aion_api.records r
		JOIN aion_api.categories c ON r.category_id = c.category_id
		WHERE r.user_id = ? AND r.deleted_at IS NULL AND c.deleted_at IS NULL
		GROUP BY r.category_id, c.name
		ORDER BY count DESC
		LIMIT 1
	`, userID).Scan(&categoryResult).Error(); err == nil && categoryResult.CategoryID > 0 {
		stats.MostUsedCategory = &domain.UsageCount{
			ID:    categoryResult.CategoryID,
			Name:  categoryResult.Name,
			Count: categoryResult.Count,
		}
	}

	// Most used tag (with JOIN and GROUP BY)
	var tagResult struct {
		TagID uint64
		Name  string
		Count int
	}
	if err := r.db.WithContext(ctx).Raw(`
		SELECT r.tag_id, t.name, COUNT(*) as count
		FROM aion_api.records r
		JOIN aion_api.tags t ON r.tag_id = t.tag_id
		WHERE r.user_id = ? AND r.deleted_at IS NULL AND t.deleted_at IS NULL
		GROUP BY r.tag_id, t.name
		ORDER BY count DESC
		LIMIT 1
	`, userID).Scan(&tagResult).Error(); err == nil && tagResult.TagID > 0 {
		stats.MostUsedTag = &domain.UsageCount{
			ID:    tagResult.TagID,
			Name:  tagResult.Name,
			Count: tagResult.Count,
		}
	}

	span.SetAttributes(
		attribute.Int("total_records", stats.TotalRecords),
		attribute.Int("total_categories", stats.TotalCategories),
		attribute.Int("total_tags", stats.TotalTags),
		attribute.Int("records_this_week", stats.RecordsThisWeek),
		attribute.Int("records_this_month", stats.RecordsThisMonth),
	)
	span.SetStatus(codes.Ok, "user stats retrieved successfully")

	r.logger.InfowCtx(ctx, "user stats retrieved successfully",
		commonkeys.UserID, strconv.FormatUint(userID, 10),
		"total_records", stats.TotalRecords,
		"records_this_month", stats.RecordsThisMonth,
	)

	return stats, nil
}
