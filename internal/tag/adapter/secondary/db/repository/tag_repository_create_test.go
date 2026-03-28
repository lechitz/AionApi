package repository_test

import (
	"errors"
	"testing"

	"github.com/lechitz/aion-api/internal/platform/ports/output/db"
	repository "github.com/lechitz/aion-api/internal/tag/adapter/secondary/db/repository"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, dbMock := newTagRepo(t)
		tag := sampleTag()

		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Create(gomock.Any()).DoAndReturn(func(any) db.DB {
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		got, err := repo.Create(t.Context(), tag)
		require.NoError(t, err)
		require.Equal(t, tag.Name, got.Name)
		require.Equal(t, tag.UserID, got.UserID)
	})

	t.Run("db error", func(t *testing.T) {
		repo, dbMock := newTagRepo(t)

		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Create(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("insert fail"))

		got, err := repo.Create(t.Context(), sampleTag())
		require.Error(t, err)
		require.Empty(t, got)
		require.Contains(t, err.Error(), "insert tag")
	})
}

func TestNew(t *testing.T) {
	repo, _ := newTagRepo(t)
	require.NotNil(t, repo)
	_ = repository.TracerName
}
