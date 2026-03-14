package repository_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/chat/adapter/secondary/db/model"
	repository "github.com/lechitz/AionApi/internal/chat/adapter/secondary/db/repository"
	"github.com/lechitz/AionApi/internal/chat/core/domain"
	"github.com/lechitz/AionApi/internal/platform/ports/output/db"
	"github.com/lechitz/AionApi/tests/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func newChatRepo(t *testing.T) (*repository.ChatHistoryRepository, *mocks.MockDB) {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	dbMock := mocks.NewMockDB(ctrl)
	logger := mocks.NewMockContextLogger(ctrl)
	logger.EXPECT().ErrorwCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().ErrorwCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().
		ErrorwCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes()
	logger.EXPECT().InfowCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().InfowCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().
		InfowCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes()

	return repository.New(dbMock, logger), dbMock
}

func TestChatHistoryRepository(t *testing.T) {
	repo, dbMock := newChatRepo(t)

	t.Run("save success", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Create(gomock.Any()).DoAndReturn(func(v any) db.DB {
			row, ok := v.(*model.ChatHistoryDB)
			require.True(t, ok)
			row.ChatID = 77
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		got, err := repo.Save(t.Context(), domain.ChatHistory{UserID: 10, Message: "hi", Response: "ok"})
		require.NoError(t, err)
		require.Equal(t, uint64(77), got.ChatID)
	})

	t.Run("save error", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Create(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("insert fail"))

		_, err := repo.Save(t.Context(), domain.ChatHistory{UserID: 10, Message: "hi", Response: "ok"})
		require.Error(t, err)
	})

	t.Run("get latest success", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Order(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Limit(2).Return(dbMock)
		dbMock.EXPECT().Find(gomock.Any()).DoAndReturn(func(dest any, _ ...any) db.DB {
			rows, ok := dest.(*[]model.ChatHistoryDB)
			require.True(t, ok)
			*rows = []model.ChatHistoryDB{{ChatID: 1}, {ChatID: 2}}
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		got, err := repo.GetLatest(t.Context(), 10, 2)
		require.NoError(t, err)
		require.Len(t, got, 2)
	})

	t.Run("get latest error", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Order(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Limit(2).Return(dbMock)
		dbMock.EXPECT().Find(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("query fail"))

		_, err := repo.GetLatest(t.Context(), 10, 2)
		require.Error(t, err)
	})

	t.Run("get by user id success", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Order(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Limit(3).Return(dbMock)
		dbMock.EXPECT().Offset(1).Return(dbMock)
		dbMock.EXPECT().Find(gomock.Any()).DoAndReturn(func(dest any, _ ...any) db.DB {
			rows, ok := dest.(*[]model.ChatHistoryDB)
			require.True(t, ok)
			*rows = []model.ChatHistoryDB{{ChatID: 9}}
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		got, err := repo.GetByUserID(t.Context(), 10, 3, 1)
		require.NoError(t, err)
		require.Len(t, got, 1)
	})

	t.Run("get by user id error", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Order(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Limit(3).Return(dbMock)
		dbMock.EXPECT().Offset(1).Return(dbMock)
		dbMock.EXPECT().Find(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("query fail"))

		_, err := repo.GetByUserID(t.Context(), 10, 3, 1)
		require.Error(t, err)
	})
}
