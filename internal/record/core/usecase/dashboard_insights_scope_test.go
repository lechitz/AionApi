package usecase

import (
	"testing"
	"time"

	"github.com/lechitz/aion-api/internal/record/core/domain"
	tagdomain "github.com/lechitz/aion-api/internal/tag/core/domain"
	"github.com/stretchr/testify/require"
)

func TestFilterInsightRecordsByScope(t *testing.T) {
	records := []domain.Record{
		{ID: 1, TagID: 10, EventTime: time.Now().UTC()},
		{ID: 2, TagID: 11, EventTime: time.Now().UTC()},
		{ID: 3, TagID: 20, EventTime: time.Now().UTC()},
	}
	tags := []tagdomain.Tag{
		{ID: 10, CategoryID: 100},
		{ID: 11, CategoryID: 100},
		{ID: 20, CategoryID: 200},
	}

	t.Run("returns all records when no scope is defined", func(t *testing.T) {
		got := filterInsightRecordsByScope(records, nil, nil, tags)
		require.Len(t, got, 3)
	})

	t.Run("filters by category", func(t *testing.T) {
		categoryID := uint64(100)
		got := filterInsightRecordsByScope(records, &categoryID, nil, tags)
		require.Len(t, got, 2)
		require.Equal(t, uint64(10), got[0].TagID)
		require.Equal(t, uint64(11), got[1].TagID)
	})

	t.Run("filters by tags", func(t *testing.T) {
		got := filterInsightRecordsByScope(records, nil, []uint64{20}, tags)
		require.Len(t, got, 1)
		require.Equal(t, uint64(20), got[0].TagID)
	})

	t.Run("combines category and tag filters", func(t *testing.T) {
		categoryID := uint64(100)
		got := filterInsightRecordsByScope(records, &categoryID, []uint64{11, 20}, tags)
		require.Len(t, got, 1)
		require.Equal(t, uint64(11), got[0].TagID)
	})
}

func TestBuildCategoryConcentrationInsight_SkipsWhenScoped(t *testing.T) {
	now := time.Now().UTC()
	records := []domain.Record{
		{ID: 1, TagID: 10, EventTime: now},
		{ID: 2, TagID: 10, EventTime: now},
		{ID: 3, TagID: 11, EventTime: now},
	}
	defs := []domain.MetricDefinition{
		{TagID: 10, TagIDs: []uint64{10}, DisplayName: "Saude Fisica"},
		{TagID: 11, TagIDs: []uint64{11}, DisplayName: "Mobilidade"},
	}

	categoryID := uint64(1)
	require.Nil(t, buildCategoryConcentrationInsight(records, defs, domain.InsightWindow7D, now, &categoryID, nil))
	require.Nil(t, buildCategoryConcentrationInsight(records, defs, domain.InsightWindow7D, now, nil, []uint64{10}))
}
