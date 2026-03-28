package repository_test

import (
	"errors"
	"testing"

	"github.com/lechitz/aion-api/internal/platform/ports/output/db"
	"github.com/lechitz/aion-api/internal/tag/adapter/secondary/db/model"
	repository "github.com/lechitz/aion-api/internal/tag/adapter/secondary/db/repository"
	"github.com/lechitz/aion-api/internal/tag/core/domain"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestGetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, dbMock := newTagRepo(t)
		tag := sampleTag()

		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Find(gomock.Any()).DoAndReturn(func(dest any, _ ...any) db.DB {
			rows, ok := dest.(*[]model.TagDB)
			require.True(t, ok)
			*rows = []model.TagDB{{ID: tag.ID, UserID: tag.UserID, CategoryID: tag.CategoryID, Name: tag.Name, Description: tag.Description, Icon: tag.Icon}}
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		got, err := repo.GetAll(t.Context(), tag.UserID)
		require.NoError(t, err)
		require.Len(t, got, 1)
		require.Equal(t, tag.Name, got[0].Name)
	})

	t.Run("not found", func(t *testing.T) {
		repo, dbMock := newTagRepo(t)
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Find(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(gorm.ErrRecordNotFound)

		got, err := repo.GetAll(t.Context(), 1)
		require.NoError(t, err)
		require.Empty(t, got)
	})

	t.Run("generic error", func(t *testing.T) {
		repo, dbMock := newTagRepo(t)
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Find(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("db fail"))

		got, err := repo.GetAll(t.Context(), 1)
		require.Error(t, err)
		require.Empty(t, got)
		require.Contains(t, err.Error(), "get all tags")
	})
}

func TestGetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, dbMock := newTagRepo(t)
		tag := sampleTag()

		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().First(gomock.Any()).DoAndReturn(func(dest any, _ ...any) db.DB {
			row, ok := dest.(*model.TagDB)
			require.True(t, ok)
			*row = model.TagDB{ID: tag.ID, UserID: tag.UserID, CategoryID: tag.CategoryID, Name: tag.Name}
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		got, err := repo.GetByID(t.Context(), tag.ID, tag.UserID)
		require.NoError(t, err)
		require.Equal(t, tag.ID, got.ID)
	})

	t.Run("not found returns empty", func(t *testing.T) {
		repo, dbMock := newTagRepo(t)
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().First(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(gorm.ErrRecordNotFound)

		got, err := repo.GetByID(t.Context(), 1, 2)
		require.NoError(t, err)
		require.Equal(t, domain.Tag{}, got)
	})

	t.Run("generic error", func(t *testing.T) {
		repo, dbMock := newTagRepo(t)
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().First(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("db fail"))

		got, err := repo.GetByID(t.Context(), 1, 2)
		require.Error(t, err)
		require.Empty(t, got)
		require.Contains(t, err.Error(), "get tag by id")
	})
}

func TestGetByName(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, dbMock := newTagRepo(t)
		tag := sampleTag()

		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().First(gomock.Any()).DoAndReturn(func(dest any, _ ...any) db.DB {
			row, ok := dest.(*model.TagDB)
			require.True(t, ok)
			*row = model.TagDB{ID: tag.ID, UserID: tag.UserID, Name: tag.Name}
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		got, err := repo.GetByName(t.Context(), tag.Name, tag.UserID)
		require.NoError(t, err)
		require.Equal(t, tag.Name, got.Name)
	})

	t.Run("not found returns explicit error", func(t *testing.T) {
		repo, dbMock := newTagRepo(t)
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().First(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(gorm.ErrRecordNotFound)

		got, err := repo.GetByName(t.Context(), "x", 1)
		require.Error(t, err)
		require.Empty(t, got)
		require.Equal(t, repository.ErrTagNotFoundMsg, err.Error())
	})

	t.Run("generic error", func(t *testing.T) {
		repo, dbMock := newTagRepo(t)
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().First(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("db fail"))

		got, err := repo.GetByName(t.Context(), "x", 1)
		require.Error(t, err)
		require.Empty(t, got)
		require.Contains(t, err.Error(), "get tag by name")
	})
}

func TestGetByCategoryID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, dbMock := newTagRepo(t)
		tag := sampleTag()

		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Find(gomock.Any()).DoAndReturn(func(dest any, _ ...any) db.DB {
			rows, ok := dest.(*[]model.TagDB)
			require.True(t, ok)
			*rows = []model.TagDB{{ID: tag.ID, UserID: tag.UserID, CategoryID: tag.CategoryID, Name: tag.Name}}
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		got, err := repo.GetByCategoryID(t.Context(), tag.CategoryID, tag.UserID)
		require.NoError(t, err)
		require.Len(t, got, 1)
		require.Equal(t, tag.CategoryID, got[0].CategoryID)
	})

	t.Run("error", func(t *testing.T) {
		repo, dbMock := newTagRepo(t)
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Find(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("db fail"))

		got, err := repo.GetByCategoryID(t.Context(), 1, 2)
		require.Error(t, err)
		require.Empty(t, got)
		require.Contains(t, err.Error(), "get tags by category id")
	})
}
