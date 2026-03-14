package repository_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/platform/ports/output/db"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/tag/adapter/secondary/db/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUpdateTag(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, dbMock := newTagRepo(t)
		tag := sampleTag()
		fields := map[string]interface{}{
			"name":                  "renamed",
			commonkeys.TagCreatedAt: tag.CreatedAt,
		}

		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Model(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Updates(gomock.Any()).DoAndReturn(func(updateFields any) db.DB {
			typed, ok := updateFields.(map[string]interface{})
			require.True(t, ok)
			_, hasCreatedAt := typed[commonkeys.TagCreatedAt]
			require.False(t, hasCreatedAt)
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().First(gomock.Any()).DoAndReturn(func(dest any, _ ...any) db.DB {
			row, ok := dest.(*model.TagDB)
			require.True(t, ok)
			*row = model.TagDB{
				ID:          tag.ID,
				UserID:      tag.UserID,
				CategoryID:  tag.CategoryID,
				Name:        "renamed",
				Description: tag.Description,
				Icon:        tag.Icon,
			}
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		got, err := repo.UpdateTag(t.Context(), tag.ID, tag.UserID, fields)
		require.NoError(t, err)
		require.Equal(t, "renamed", got.Name)
	})

	t.Run("update error", func(t *testing.T) {
		repo, dbMock := newTagRepo(t)
		tag := sampleTag()

		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Model(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Updates(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("update fail"))

		got, err := repo.UpdateTag(t.Context(), tag.ID, tag.UserID, map[string]interface{}{"name": "x"})
		require.Error(t, err)
		require.Empty(t, got)
		require.ErrorContains(t, err, "update tag")
	})

	t.Run("fetch updated error", func(t *testing.T) {
		repo, dbMock := newTagRepo(t)
		tag := sampleTag()

		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Model(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Updates(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(nil)

		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().First(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("fetch fail"))

		got, err := repo.UpdateTag(t.Context(), tag.ID, tag.UserID, map[string]interface{}{"name": "x"})
		require.Error(t, err)
		require.Empty(t, got)
		require.ErrorContains(t, err, "fetch updated tag")
	})
}
