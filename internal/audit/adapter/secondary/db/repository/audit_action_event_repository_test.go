package repository_test

import (
	"errors"
	"testing"
	"time"

	"github.com/lechitz/aion-api/internal/audit/adapter/secondary/db/model"
	repository "github.com/lechitz/aion-api/internal/audit/adapter/secondary/db/repository"
	"github.com/lechitz/aion-api/internal/audit/core/domain"
	"github.com/lechitz/aion-api/internal/platform/ports/output/db"
	"github.com/lechitz/aion-api/tests/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func newAuditRepo(t *testing.T) (*repository.AuditActionEventRepository, *mocks.MockDB) {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	dbMock := mocks.NewMockDB(ctrl)
	logger := mocks.NewMockContextLogger(ctrl)
	logger.EXPECT().ErrorwCtx(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().ErrorwCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().ErrorwCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().ErrorwCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().
		ErrorwCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes()
	logger.EXPECT().InfowCtx(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().InfowCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().InfowCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().InfowCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().
		InfowCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes()

	return repository.NewAuditActionEventRepository(dbMock, logger), dbMock
}

func TestAuditActionEventRepository_Save(t *testing.T) {
	repo, dbMock := newAuditRepo(t)
	event := domain.AuditActionEvent{
		EventID:      "evt-1",
		TimestampUTC: time.Now().UTC(),
		UserID:       10,
		Source:       "aion-api",
		TraceID:      "trace-1",
		UIActionType: "draft_accept",
		DraftID:      "draft-1",
		Action:       "draft_accept",
		Entity:       "category",
		Operation:    "create",
		Status:       "success",
	}

	t.Run("success", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Create(gomock.Any()).DoAndReturn(func(v any) db.DB {
			row, ok := v.(*model.AuditActionEventDB)
			require.True(t, ok)
			require.Equal(t, "evt-1", row.EventID)
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		err := repo.Save(t.Context(), event)
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Create(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("insert fail"))

		err := repo.Save(t.Context(), event)
		require.Error(t, err)
	})
}

func TestAuditActionEventRepository_List(t *testing.T) {
	repo, dbMock := newAuditRepo(t)

	from := time.Now().UTC().Add(-1 * time.Hour)
	to := time.Now().UTC()
	userID := uint64(9)
	filter := domain.AuditActionEventFilter{
		UserID:   &userID,
		TraceID:  "trace-9",
		DraftID:  "draft-9",
		Statuses: []string{"failed", "blocked"},
		FromUTC:  &from,
		ToUTC:    &to,
		Limit:    20,
		Offset:   5,
	}

	t.Run("success", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Model(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock).Times(6)
		dbMock.EXPECT().Order(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Limit(20).Return(dbMock)
		dbMock.EXPECT().Offset(5).Return(dbMock)
		dbMock.EXPECT().Find(gomock.Any()).DoAndReturn(func(dest any, _ ...any) db.DB {
			rows, ok := dest.(*[]model.AuditActionEventDB)
			require.True(t, ok)
			*rows = []model.AuditActionEventDB{
				{EventID: "evt-1", TraceID: "trace-9", DraftID: "draft-9", Status: "failed"},
				{EventID: "evt-2", TraceID: "trace-9", DraftID: "draft-9", Status: "blocked"},
			}
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		got, err := repo.List(t.Context(), filter)
		require.NoError(t, err)
		require.Len(t, got, 2)
		require.Equal(t, "evt-1", got[0].EventID)
	})

	t.Run("error", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Model(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock).Times(6)
		dbMock.EXPECT().Order(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Limit(20).Return(dbMock)
		dbMock.EXPECT().Offset(5).Return(dbMock)
		dbMock.EXPECT().Find(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("query fail"))

		_, err := repo.List(t.Context(), filter)
		require.Error(t, err)
	})
}
