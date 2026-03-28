package mapper_test

import (
	"testing"
	"time"

	"github.com/lechitz/aion-api/internal/category/adapter/secondary/db/mapper"
	"github.com/lechitz/aion-api/internal/category/adapter/secondary/db/model"
	"github.com/lechitz/aion-api/internal/category/core/domain"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCategoryFromDBAndToDB(t *testing.T) {
	now := time.Now().UTC()
	deleted := now.Add(-2 * time.Hour)
	lastUsed := now.Add(-time.Minute)

	db := model.CategoryDB{
		ID:          5,
		UserID:      7,
		Name:        "Work",
		Description: "work tasks",
		Color:       "#000",
		Icon:        "briefcase",
		CreatedAt:   now,
		UpdatedAt:   now,
		DeletedAt:   gorm.DeletedAt{Time: deleted, Valid: true},
		UsageCount:  12,
		LastUsedAt:  &lastUsed,
	}

	d := mapper.CategoryFromDB(db)
	require.NotNil(t, d.DeletedAt)
	require.Equal(t, deleted, *d.DeletedAt)
	require.Equal(t, db.UsageCount, d.UsageCount)

	back := mapper.CategoryToDB(domain.Category{
		ID:          d.ID,
		UserID:      d.UserID,
		Name:        d.Name,
		Description: d.Description,
		Color:       d.Color,
		Icon:        d.Icon,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
		DeletedAt:   d.DeletedAt,
		UsageCount:  d.UsageCount,
		LastUsedAt:  d.LastUsedAt,
	})
	require.True(t, back.DeletedAt.Valid)
	require.Equal(t, deleted, back.DeletedAt.Time)

	back = mapper.CategoryToDB(domain.Category{})
	require.False(t, back.DeletedAt.Valid)
}
