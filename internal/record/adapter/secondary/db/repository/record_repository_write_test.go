package repository_test

import (
	"errors"
	"testing"

	"github.com/lechitz/aion-api/internal/platform/ports/output/db"
	"github.com/lechitz/aion-api/internal/record/adapter/secondary/db/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRecordCreateUpdateDelete(t *testing.T) {
	repo, dbMock := newRecordRepo(t)
	rec := sampleRecord()

	t.Run("create success", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Create(gomock.Any()).DoAndReturn(func(v any) db.DB {
			row, ok := v.(*model.Record)
			require.True(t, ok)
			row.ID = 99
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		got, err := repo.Create(t.Context(), rec)
		require.NoError(t, err)
		require.Equal(t, uint64(99), got.ID)
	})

	t.Run("create error", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Create(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("create fail"))
		_, err := repo.Create(t.Context(), rec)
		require.Error(t, err)
	})

	t.Run("update error", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Model(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Updates(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("update fail"))
		_, err := repo.Update(t.Context(), rec)
		require.Error(t, err)
	})

	t.Run("update success", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Model(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Updates(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(nil)

		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().First(gomock.Any()).DoAndReturn(func(dest any, _ ...any) db.DB {
			row, ok := dest.(*model.Record)
			require.True(t, ok)
			row.ID = rec.ID
			row.UserID = rec.UserID
			row.TagID = rec.TagID
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		got, err := repo.Update(t.Context(), rec)
		require.NoError(t, err)
		require.Equal(t, rec.ID, got.ID)
	})

	t.Run("delete error", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Model(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("delete fail"))
		require.Error(t, repo.Delete(t.Context(), rec.ID, rec.UserID))
	})

	t.Run("delete all by user success", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Model(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(nil)
		require.NoError(t, repo.DeleteAllByUser(t.Context(), rec.UserID))
	})
}
