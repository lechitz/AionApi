package repository_test

import (
	"errors"
	"testing"
	"time"

	"github.com/lechitz/aion-api/internal/platform/ports/output/db"
	"github.com/lechitz/aion-api/internal/record/adapter/secondary/db/model"
	"github.com/lechitz/aion-api/internal/record/core/domain"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRecordSearchRecordsSuccessWithFullTextAndFilters(t *testing.T) {
	repo, dbMock, loggerMock := newRecordRepoWithLogger(t)
	rec := sampleRecord()
	start := rec.EventTime.Add(-2 * time.Hour)
	end := rec.EventTime.Add(2 * time.Hour)
	filters := domain.SearchFilters{
		Query:       "focus",
		CategoryIDs: []uint64{11},
		TagIDs:      []uint64{22, 23},
		StartDate:   &start,
		EndDate:     &end,
		Limit:       10,
		Offset:      5,
	}

	dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
	dbMock.EXPECT().Raw(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).DoAndReturn(
		func(sql string, values ...any) db.DB {
			require.Contains(t, sql, "ts_rank")
			require.Contains(t, sql, "t.category_id = ANY($3)")
			require.Contains(t, sql, "tag_id = ANY($4)")
			require.Contains(t, sql, "event_time >= $5")
			require.Contains(t, sql, "event_time <= $6")
			require.Contains(t, sql, "ORDER BY rank DESC, created_at DESC")
			require.Contains(t, sql, "LIMIT $7 OFFSET $8")
			require.Len(t, values, 8)
			return dbMock
		},
	)
	dbMock.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest any) db.DB {
		rows, ok := dest.(*[]model.Record)
		require.True(t, ok)
		*rows = []model.Record{{ID: rec.ID, UserID: rec.UserID, TagID: rec.TagID}}
		return dbMock
	})
	dbMock.EXPECT().Error().Return(nil)
	loggerMock.EXPECT().InfowCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

	got, err := repo.SearchRecords(t.Context(), rec.UserID, filters)
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, rec.ID, got[0].ID)
}

func TestRecordSearchRecordsSuccessWithoutQuery(t *testing.T) {
	repo, dbMock, loggerMock := newRecordRepoWithLogger(t)
	rec := sampleRecord()
	filters := domain.SearchFilters{
		CategoryIDs: []uint64{11},
		Limit:       20,
		Offset:      0,
	}

	dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
	dbMock.EXPECT().Raw(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(sql string, values ...any) db.DB {
			require.Contains(t, sql, "0 as rank")
			require.Contains(t, sql, "t.category_id = ANY($2)")
			require.Contains(t, sql, "ORDER BY created_at DESC")
			require.Contains(t, sql, "LIMIT $3 OFFSET $4")
			require.Len(t, values, 4)
			return dbMock
		},
	)
	dbMock.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest any) db.DB {
		rows, ok := dest.(*[]model.Record)
		require.True(t, ok)
		*rows = []model.Record{{ID: rec.ID, UserID: rec.UserID, TagID: rec.TagID}}
		return dbMock
	})
	dbMock.EXPECT().Error().Return(nil)
	loggerMock.EXPECT().InfowCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

	got, err := repo.SearchRecords(t.Context(), rec.UserID, filters)
	require.NoError(t, err)
	require.Len(t, got, 1)
}

func TestRecordSearchRecordsError(t *testing.T) {
	repo, dbMock, loggerMock := newRecordRepoWithLogger(t)
	rec := sampleRecord()
	expectedErr := errors.New("search failed")

	dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
	dbMock.EXPECT().Raw(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
	dbMock.EXPECT().Scan(gomock.Any()).Return(dbMock)
	dbMock.EXPECT().Error().Return(expectedErr)
	loggerMock.EXPECT().ErrorwCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

	got, err := repo.SearchRecords(t.Context(), rec.UserID, domain.SearchFilters{Query: "x", Limit: 2, Offset: 1})
	require.Error(t, err)
	require.Nil(t, got)
	require.ErrorContains(t, err, "search records")
}
