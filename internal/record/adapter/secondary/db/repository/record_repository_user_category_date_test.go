package repository_test

import (
	"errors"
	"testing"

	"github.com/lechitz/aion-api/internal/platform/ports/output/db"
	"github.com/lechitz/aion-api/internal/record/adapter/secondary/db/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRecordGetByUserCategoryDateSuccess(t *testing.T) {
	repo, dbMock := newRecordRepo(t)
	rec := sampleRecord()

	dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
	dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
	dbMock.EXPECT().First(gomock.Any()).DoAndReturn(func(dest any, _ ...any) db.DB {
		row, ok := dest.(*model.Record)
		require.True(t, ok)
		*row = model.Record{ID: rec.ID, UserID: rec.UserID, TagID: rec.TagID, EventTime: rec.EventTime}
		return dbMock
	})
	dbMock.EXPECT().Error().Return(nil)

	got, err := repo.GetByUserCategoryDate(t.Context(), rec.UserID, 77, rec.EventTime)
	require.NoError(t, err)
	require.Equal(t, rec.ID, got.ID)
}

func TestRecordGetByUserCategoryDateError(t *testing.T) {
	repo, dbMock := newRecordRepo(t)
	rec := sampleRecord()

	dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
	dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
	dbMock.EXPECT().First(gomock.Any()).Return(dbMock)
	dbMock.EXPECT().Error().Return(errors.New("not found"))

	got, err := repo.GetByUserCategoryDate(t.Context(), rec.UserID, 77, rec.EventTime)
	require.Error(t, err)
	require.Equal(t, uint64(0), got.ID)
}
