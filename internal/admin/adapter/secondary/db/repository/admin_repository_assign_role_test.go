package repository_test

import (
	"errors"
	"testing"

	repository "github.com/lechitz/AionApi/internal/admin/adapter/secondary/db/repository"
	"github.com/lechitz/AionApi/internal/platform/ports/output/db"
	"github.com/lechitz/AionApi/tests/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAssignDefaultRole(t *testing.T) {
	t.Run("count query error", func(t *testing.T) {
		repo, dbMock := newAdminRepo(t)
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Raw(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Scan(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("count fail"))

		err := repo.AssignDefaultRole(t.Context(), 10)
		require.Error(t, err)
		require.ErrorIs(t, err, repository.ErrCheckDefaultRoleExists)
	})

	t.Run("default role not found", func(t *testing.T) {
		repo, dbMock := newAdminRepo(t)
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Raw(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest any) db.DB {
			count, ok := dest.(*int64)
			require.True(t, ok)
			*count = 0
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		err := repo.AssignDefaultRole(t.Context(), 10)
		require.ErrorIs(t, err, repository.ErrDefaultRoleNotFound)
	})

	t.Run("exec insert error", func(t *testing.T) {
		repo, dbMock := newAdminRepo(t)
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Raw(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest any) db.DB {
			count, ok := dest.(*int64)
			require.True(t, ok)
			*count = 1
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("insert fail"))

		err := repo.AssignDefaultRole(t.Context(), 10)
		require.Error(t, err)
		require.ErrorIs(t, err, repository.ErrAssignDefaultRole)
	})

	t.Run("success", func(t *testing.T) {
		repo, dbMock := newAdminRepo(t)
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Raw(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest any) db.DB {
			count, ok := dest.(*int64)
			require.True(t, ok)
			*count = 1
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(nil)

		require.NoError(t, repo.AssignDefaultRole(t.Context(), 10))
	})
}

func TestAssignDefaultRoleWithTx(t *testing.T) {
	repo, dbMock := newAdminRepo(t)
	ctrl := gomock.NewController(t)
	txMock := mocks.NewMockDB(ctrl)
	t.Cleanup(ctrl.Finish)

	txMock.EXPECT().WithContext(gomock.Any()).Return(txMock)
	txMock.EXPECT().Raw(gomock.Any(), gomock.Any()).Return(txMock)
	txMock.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest any) db.DB {
		count, ok := dest.(*int64)
		require.True(t, ok)
		*count = 1
		return txMock
	})
	txMock.EXPECT().Error().Return(nil)

	txMock.EXPECT().WithContext(gomock.Any()).Return(txMock)
	txMock.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(txMock)
	txMock.EXPECT().Error().Return(nil)

	require.NoError(t, repo.AssignDefaultRoleWithTx(t.Context(), txMock, 11))

	dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
	dbMock.EXPECT().Raw(gomock.Any(), gomock.Any()).Return(dbMock)
	dbMock.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest any) db.DB {
		count, ok := dest.(*int64)
		require.True(t, ok)
		*count = 1
		return dbMock
	})
	dbMock.EXPECT().Error().Return(nil)
	dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
	dbMock.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
	dbMock.EXPECT().Error().Return(nil)

	require.NoError(t, repo.AssignDefaultRoleWithTx(t.Context(), nil, 12))
}
