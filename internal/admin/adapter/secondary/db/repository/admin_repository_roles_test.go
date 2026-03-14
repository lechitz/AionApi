package repository_test

import (
	"errors"
	"testing"

	adminmodel "github.com/lechitz/AionApi/internal/admin/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/platform/ports/output/db"
	usermodel "github.com/lechitz/AionApi/internal/user/adapter/secondary/db/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetRolesByUserID(t *testing.T) {
	repo, dbMock := newAdminRepo(t)

	t.Run("query error", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Raw(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Scan(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("db fail"))

		roles, err := repo.GetRolesByUserID(t.Context(), 1)
		require.Error(t, err)
		require.Nil(t, roles)
	})

	t.Run("defaults to user role", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Raw(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest any) db.DB {
			rows, ok := dest.(*[]adminmodel.RoleDB)
			require.True(t, ok)
			*rows = []adminmodel.RoleDB{}
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		roles, err := repo.GetRolesByUserID(t.Context(), 1)
		require.NoError(t, err)
		require.Equal(t, []string{"user"}, roles)
	})

	t.Run("returns fetched roles", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Raw(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest any) db.DB {
			rows, ok := dest.(*[]adminmodel.RoleDB)
			require.True(t, ok)
			*rows = []adminmodel.RoleDB{{Name: "admin"}, {Name: "owner"}}
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		roles, err := repo.GetRolesByUserID(t.Context(), 1)
		require.NoError(t, err)
		require.Equal(t, []string{"admin", "owner"}, roles)
	})
}

func TestGetByID(t *testing.T) {
	repo, dbMock := newAdminRepo(t)

	t.Run("user lookup error", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().First(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("user fail"))

		_, err := repo.GetByID(t.Context(), 10)
		require.Error(t, err)
	})

	t.Run("roles lookup error", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().First(gomock.Any()).DoAndReturn(func(dest any, _ ...any) db.DB {
			user, ok := dest.(*usermodel.UserDB)
			require.True(t, ok)
			user.ID = 10
			user.Username = "u"
			user.Email = "u@example.com"
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Raw(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Scan(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("roles fail"))

		_, err := repo.GetByID(t.Context(), 10)
		require.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().First(gomock.Any()).DoAndReturn(func(dest any, _ ...any) db.DB {
			user, ok := dest.(*usermodel.UserDB)
			require.True(t, ok)
			user.ID = 10
			user.Username = "u"
			user.Email = "u@example.com"
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Raw(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest any) db.DB {
			rows, ok := dest.(*[]adminmodel.RoleDB)
			require.True(t, ok)
			*rows = []adminmodel.RoleDB{{Name: "admin"}}
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		got, err := repo.GetByID(t.Context(), 10)
		require.NoError(t, err)
		require.Equal(t, uint64(10), got.ID)
		require.Equal(t, []string{"admin"}, got.Roles)
	})
}

func TestUpdateRoles(t *testing.T) {
	repo, dbMock := newAdminRepo(t)
	roles := []string{"admin"}

	t.Run("transaction error", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(fc func(db.DB) error) error {
			tx := dbMock
			tx.EXPECT().Where(gomock.Any(), gomock.Any()).Return(tx)
			tx.EXPECT().Delete(gomock.Any()).Return(tx)
			tx.EXPECT().Error().Return(errors.New("delete fail"))
			return fc(tx)
		})

		_, err := repo.UpdateRoles(t.Context(), 7, roles)
		require.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(fc func(db.DB) error) error {
			tx := dbMock
			tx.EXPECT().Where("user_id = ?", uint64(7)).Return(tx)
			tx.EXPECT().Delete(gomock.Any()).Return(tx)
			tx.EXPECT().Error().Return(nil)

			tx.EXPECT().Where("name IN ?", roles).Return(tx)
			tx.EXPECT().Where("is_active = ?", true).Return(tx)
			tx.EXPECT().Find(gomock.Any()).DoAndReturn(func(dest any, _ ...any) db.DB {
				rows, ok := dest.(*[]adminmodel.RoleDB)
				require.True(t, ok)
				*rows = []adminmodel.RoleDB{{ID: 2, Name: "admin"}}
				return tx
			})
			tx.EXPECT().Error().Return(nil)

			tx.EXPECT().Create(gomock.Any()).Return(tx)
			tx.EXPECT().Error().Return(nil)

			tx.EXPECT().Model(gomock.Any()).Return(tx)
			tx.EXPECT().Where("user_id = ?", uint64(7)).Return(tx)
			tx.EXPECT().Update(gomock.Any(), gomock.Any()).Return(tx)
			tx.EXPECT().Error().Return(nil)
			return fc(tx)
		})

		// GetByID after transaction
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().First(gomock.Any()).DoAndReturn(func(dest any, _ ...any) db.DB {
			user, ok := dest.(*usermodel.UserDB)
			require.True(t, ok)
			user.ID = 7
			user.Username = "u"
			user.Email = "u@example.com"
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Raw(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest any) db.DB {
			rows, ok := dest.(*[]adminmodel.RoleDB)
			require.True(t, ok)
			*rows = []adminmodel.RoleDB{{Name: "admin"}}
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		got, err := repo.UpdateRoles(t.Context(), 7, roles)
		require.NoError(t, err)
		require.Equal(t, []string{"admin"}, got.Roles)
	})
}
