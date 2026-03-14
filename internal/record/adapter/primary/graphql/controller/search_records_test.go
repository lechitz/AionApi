package controller_test

import (
	"context"
	"testing"
	"time"

	gmodel "github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
	"github.com/lechitz/AionApi/internal/record/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchRecords_ConvertsFilters(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	limit := int32(5)
	offset := int32(2)

	svc := &recordServiceStub{
		searchFn: func(_ context.Context, userID uint64, filters domain.SearchFilters) ([]domain.Record, error) {
			require.Equal(t, uint64(7), userID)
			require.Equal(t, "hello", filters.Query)
			assert.Equal(t, []uint64{1}, filters.CategoryIDs)
			assert.Equal(t, []uint64{2}, filters.TagIDs)
			require.NotNil(t, filters.StartDate)
			require.NotNil(t, filters.EndDate)
			assert.Equal(t, start, *filters.StartDate)
			assert.Equal(t, end, *filters.EndDate)
			assert.Equal(t, int(limit), filters.Limit)
			assert.Equal(t, int(offset), filters.Offset)
			return []domain.Record{}, nil
		},
	}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	startStr := start.Format(time.RFC3339)
	endStr := end.Format(time.RFC3339)
	filters := gmodel.SearchFilters{
		Query:       "hello",
		CategoryIds: []string{"1", "bad"},
		TagIds:      []string{"2"},
		StartDate:   &startStr,
		EndDate:     &endStr,
		Limit:       &limit,
		Offset:      &offset,
	}

	out, err := h.SearchRecords(t.Context(), filters, 7)
	require.NoError(t, err)
	assert.Empty(t, out)
}

func TestSearchRecords_DefaultLimit(t *testing.T) {
	svc := &recordServiceStub{
		searchFn: func(_ context.Context, _ uint64, filters domain.SearchFilters) ([]domain.Record, error) {
			assert.Equal(t, 20, filters.Limit)
			return []domain.Record{}, nil
		},
	}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	filters := gmodel.SearchFilters{Query: "hello"}
	_, err := h.SearchRecords(t.Context(), filters, 1)
	require.NoError(t, err)
}
