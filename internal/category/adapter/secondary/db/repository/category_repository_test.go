package repository_test

import (
	"errors"
	"testing"

	"github.com/lechitz/aion-api/internal/category/adapter/secondary/db/model"
	repository "github.com/lechitz/aion-api/internal/category/adapter/secondary/db/repository"
	"github.com/lechitz/aion-api/internal/platform/ports/output/db"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestCategoryCreate(t *testing.T) {
	repo, dbMock := newCategoryRepo(t)

	t.Run("success", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Create(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(nil)

		got, err := repo.Create(t.Context(), sampleCategory())
		require.NoError(t, err)
		require.Equal(t, "health", got.Name)
	})

	t.Run("error", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Create(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("db fail"))

		got, err := repo.Create(t.Context(), sampleCategory())
		require.Error(t, err)
		require.Empty(t, got)
		require.Contains(t, err.Error(), repository.ErrCreateCategoryMsg)
	})
}

func TestCategoryGetByID(t *testing.T) {
	repo, dbMock := newCategoryRepo(t)

	t.Run("success", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Select(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().First(gomock.Any()).DoAndReturn(func(dest any, _ ...any) db.DB {
			row, ok := dest.(*model.CategoryDB)
			require.True(t, ok)
			*row = model.CategoryDB{ID: 1, UserID: 2, Name: "health"}
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		got, err := repo.GetByID(t.Context(), 1, 2)
		require.NoError(t, err)
		require.Equal(t, uint64(1), got.ID)
	})

	t.Run("not found", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Select(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().First(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(gorm.ErrRecordNotFound)

		_, err := repo.GetByID(t.Context(), 1, 2)
		require.EqualError(t, err, repository.ErrCategoryNotFoundMsg)
	})

	t.Run("generic", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Select(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().First(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("x"))

		_, err := repo.GetByID(t.Context(), 1, 2)
		require.EqualError(t, err, repository.ErrGetCategoryMsg)
	})
}

func TestCategoryGetByNameAndListAll(t *testing.T) {
	repo, dbMock := newCategoryRepo(t)

	t.Run("get by name success", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().First(gomock.Any()).DoAndReturn(func(dest any, _ ...any) db.DB {
			row, ok := dest.(*model.CategoryDB)
			require.True(t, ok)
			*row = model.CategoryDB{ID: 9, Name: "health"}
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		got, err := repo.GetByName(t.Context(), "health", 2)
		require.NoError(t, err)
		require.Equal(t, uint64(9), got.ID)
	})

	t.Run("get by name not found", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().First(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(gorm.ErrRecordNotFound)

		_, err := repo.GetByName(t.Context(), "x", 2)
		require.EqualError(t, err, repository.ErrCategoryNotFoundMsg)
	})

	t.Run("get by name generic", func(t *testing.T) {
		errExp := errors.New("db fail")
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().First(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errExp)

		_, err := repo.GetByName(t.Context(), "x", 2)
		require.ErrorIs(t, err, errExp)
	})

	t.Run("list all success", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Select(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Find(gomock.Any()).DoAndReturn(func(dest any, _ ...any) db.DB {
			rows, ok := dest.(*[]model.CategoryDB)
			require.True(t, ok)
			*rows = []model.CategoryDB{{ID: 1, Name: "h"}}
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		got, err := repo.ListAll(t.Context(), 2)
		require.NoError(t, err)
		require.Len(t, got, 1)
	})

	t.Run("list all error", func(t *testing.T) {
		errExp := errors.New("db fail")
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Select(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Find(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errExp)

		got, err := repo.ListAll(t.Context(), 2)
		require.ErrorIs(t, err, errExp)
		require.Nil(t, got)
	})
}

func TestCategoryUpdateAndSoftDelete(t *testing.T) {
	repo, dbMock := newCategoryRepo(t)

	t.Run("update no fields", func(t *testing.T) {
		_, err := repo.UpdateCategory(t.Context(), 1, 2, map[string]interface{}{commonkeys.CategoryCreatedAt: "x"})
		require.EqualError(t, err, "no fields to update")
	})

	t.Run("update query error", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Model(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Updates(gomock.Any()).Return(dbMock)
		errExp := errors.New("update fail")
		dbMock.EXPECT().Error().Return(errExp)

		_, err := repo.UpdateCategory(t.Context(), 1, 2, map[string]interface{}{"name": "new"})
		require.ErrorIs(t, err, errExp)
	})

	t.Run("update not found", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Model(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Updates(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(nil)
		dbMock.EXPECT().RowsAffected().Return(int64(0))

		_, err := repo.UpdateCategory(t.Context(), 1, 2, map[string]interface{}{"name": "new"})
		require.EqualError(t, err, repository.ErrCategoryNotFoundMsg)
	})

	t.Run("update success", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Model(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Updates(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(nil)
		dbMock.EXPECT().RowsAffected().Return(int64(1))

		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().First(gomock.Any()).DoAndReturn(func(dest any, _ ...any) db.DB {
			row, ok := dest.(*model.CategoryDB)
			require.True(t, ok)
			*row = model.CategoryDB{ID: 1, UserID: 2, Name: "new"}
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		got, err := repo.UpdateCategory(t.Context(), 1, 2, map[string]interface{}{"name": "new"})
		require.NoError(t, err)
		require.Equal(t, "new", got.Name)
	})

	t.Run("soft delete update error", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Model(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Updates(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("delete fail"))

		err := repo.SoftDelete(t.Context(), 1, 2)
		require.Error(t, err)
	})

	t.Run("soft delete not found", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Model(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Updates(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(nil)
		dbMock.EXPECT().RowsAffected().Return(int64(0))

		err := repo.SoftDelete(t.Context(), 1, 2)
		require.EqualError(t, err, repository.ErrCategoryNotFoundMsg)
	})

	t.Run("soft delete success", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Model(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Updates(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(nil)
		dbMock.EXPECT().RowsAffected().Return(int64(1))

		err := repo.SoftDelete(t.Context(), 1, 2)
		require.NoError(t, err)
	})
}
