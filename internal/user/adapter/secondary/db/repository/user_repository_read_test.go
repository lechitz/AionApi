package repository_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/platform/ports/output/db"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	usermodel "github.com/lechitz/AionApi/internal/user/adapter/secondary/db/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestGetByIDAndListAll(t *testing.T) {
	repo, dbMock, _ := newUserRepo(t)

	t.Run("get by id success", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Model(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().First(gomock.Any()).DoAndReturn(func(dest any, _ ...any) db.DB {
			row, ok := dest.(*usermodel.UserDB)
			require.True(t, ok)
			*row = usermodel.UserDB{ID: 10, Username: "john", Email: "john@example.com"}
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		got, err := repo.GetByID(t.Context(), 10)
		require.NoError(t, err)
		require.Equal(t, uint64(10), got.ID)
	})

	t.Run("get by id error", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Model(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().First(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("db fail"))

		_, err := repo.GetByID(t.Context(), 10)
		require.Error(t, err)
	})

	t.Run("list all success", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Model(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Select(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Find(gomock.Any()).DoAndReturn(func(dest any, _ ...any) db.DB {
			rows, ok := dest.(*[]usermodel.UserDB)
			require.True(t, ok)
			*rows = []usermodel.UserDB{{ID: 1}, {ID: 2}}
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		got, err := repo.ListAll(t.Context())
		require.NoError(t, err)
		require.Len(t, got, 2)
	})

	t.Run("list all error", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Model(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Select(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Find(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("list fail"))

		got, err := repo.ListAll(t.Context())
		require.Error(t, err)
		require.Nil(t, got)
	})
}

func TestGetByUsername(t *testing.T) {
	repo, dbMock, _ := newUserRepo(t)

	t.Run("success", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Select(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().First(gomock.Any()).DoAndReturn(func(dest any, _ ...any) db.DB {
			row, ok := dest.(*usermodel.UserDB)
			require.True(t, ok)
			*row = usermodel.UserDB{ID: 11, Username: "john"}
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		got, err := repo.GetByUsername(t.Context(), "john")
		require.NoError(t, err)
		require.Equal(t, uint64(11), got.ID)
	})

	t.Run("not found", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Select(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().First(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(gorm.ErrRecordNotFound)

		got, err := repo.GetByUsername(t.Context(), "john")
		require.NoError(t, err)
		require.Empty(t, got)
	})

	t.Run("generic error", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Select(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().First(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("db fail"))

		_, err := repo.GetByUsername(t.Context(), "john")
		require.Error(t, err)
	})
}

func TestUpdateAndSoftDelete(t *testing.T) {
	repo, dbMock, _ := newUserRepo(t)

	t.Run("update success", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Model(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Updates(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(nil)

		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Model(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().First(gomock.Any()).DoAndReturn(func(dest any, _ ...any) db.DB {
			row, ok := dest.(*usermodel.UserDB)
			require.True(t, ok)
			*row = usermodel.UserDB{ID: 5, Username: "newname"}
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		got, err := repo.Update(t.Context(), 5, map[string]interface{}{"username": "newname", commonkeys.UserCreatedAt: "ignore"})
		require.NoError(t, err)
		require.Equal(t, "newname", got.Username)
	})

	t.Run("soft delete error", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Model(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Updates(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("soft fail"))

		err := repo.SoftDelete(t.Context(), 5)
		require.Error(t, err)
	})

	t.Run("soft delete success", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Model(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Updates(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(nil)

		err := repo.SoftDelete(t.Context(), 5)
		require.NoError(t, err)
	})

	t.Run("update error", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Model(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Updates(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("update fail"))

		got, err := repo.Update(t.Context(), 5, map[string]interface{}{"username": "newname"})
		require.Error(t, err)
		require.Empty(t, got)
	})
}
