package graphql

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
)

// UserStats is the resolver for the userStats field.
func (r *queryResolver) UserStats(ctx context.Context) (*model.UserStats, error) {
	userID, _ := ctx.Value(ctxkeys.UserID).(uint64)

	stats, err := r.UserService.GetUserStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Safe conversion with range validation (following project pattern)
	safeInt32 := func(v int) int32 {
		if v >= -2147483648 && v <= 2147483647 {
			return int32(v) // #nosec G115 - validated range
		}
		return 0 // fallback for overflow
	}

	result := &model.UserStats{
		TotalRecords:     safeInt32(stats.TotalRecords),
		TotalCategories:  safeInt32(stats.TotalCategories),
		TotalTags:        safeInt32(stats.TotalTags),
		RecordsThisWeek:  safeInt32(stats.RecordsThisWeek),
		RecordsThisMonth: safeInt32(stats.RecordsThisMonth),
	}

	if stats.MostUsedCategory != nil {
		result.MostUsedCategory = &model.CategoryCount{
			ID:    strconv.FormatUint(stats.MostUsedCategory.ID, 10),
			Name:  stats.MostUsedCategory.Name,
			Count: safeInt32(stats.MostUsedCategory.Count),
		}
	}

	if stats.MostUsedTag != nil {
		result.MostUsedTag = &model.TagCount{
			ID:    strconv.FormatUint(stats.MostUsedTag.ID, 10),
			Name:  stats.MostUsedTag.Name,
			Count: safeInt32(stats.MostUsedTag.Count),
		}
	}

	return result, nil
}
