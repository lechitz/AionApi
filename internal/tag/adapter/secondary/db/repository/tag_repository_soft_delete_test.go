package repository_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSoftDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, dbMock := newTagRepo(t)
		tag := sampleTag()

		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Model(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Updates(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(nil)

		err := repo.SoftDelete(t.Context(), tag.ID, tag.UserID)
		require.NoError(t, err)
	})

	t.Run("db error", func(t *testing.T) {
		repo, dbMock := newTagRepo(t)
		tag := sampleTag()

		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Model(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Updates(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("soft delete fail"))

		err := repo.SoftDelete(t.Context(), tag.ID, tag.UserID)
		require.Error(t, err)
	})
}
