package repository

import (
	"context"
	"testing"
	"time"

	"github.com/lechitz/AionApi/internal/eventoutbox/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/platform/ports/output/db"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/tests/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func newStatsRepo(t *testing.T) (*EventRepository, *mocks.MockDB) {
	t.Helper()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	dbMock := mocks.NewMockDB(ctrl)
	return NewEventRepository(dbMock, noopStatsLogger{}), dbMock
}

func TestEventRepository_GetStats(t *testing.T) {
	repo, dbMock := newStatsRepo(t)
	oldest := time.Date(2026, 3, 13, 20, 0, 0, 0, time.UTC)

	dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
	dbMock.EXPECT().Raw(gomock.Any()).Return(dbMock)
	dbMock.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest any) db.DB {
		row, ok := dest.(*outboxStatsRow)
		require.True(t, ok)
		row.PendingCount = 3
		row.PublishedCount = 9
		row.FailedCount = 1
		row.OldestPendingAtUTC = &oldest
		return dbMock
	})
	dbMock.EXPECT().Error().Return(nil)

	stats, err := repo.GetStats(t.Context())
	require.NoError(t, err)
	require.Equal(t, int64(3), stats.PendingCount)
	require.Equal(t, int64(9), stats.PublishedCount)
	require.Equal(t, int64(1), stats.FailedCount)
	require.NotNil(t, stats.OldestPendingAtUTC)
	require.True(t, oldest.Equal(*stats.OldestPendingAtUTC))
}

func TestEventRepository_ListFailed(t *testing.T) {
	repo, dbMock := newStatsRepo(t)

	dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
	dbMock.EXPECT().Where("status = ? AND last_error IS NOT NULL AND last_error <> ''", "pending").Return(dbMock)
	dbMock.EXPECT().Order("available_at_utc ASC, id ASC").Return(dbMock)
	dbMock.EXPECT().Limit(2).Return(dbMock)
	dbMock.EXPECT().Find(gomock.Any()).DoAndReturn(func(dest any, _ ...any) db.DB {
		rows, ok := dest.(*[]model.EventDB)
		require.True(t, ok)
		*rows = []model.EventDB{{
			EventID:       "evt-1",
			AggregateType: "record",
			AggregateID:   "1",
			EventType:     "record.created",
			EventVersion:  "v1",
			Status:        "pending",
			LastError:     "boom",
		}}
		return dbMock
	})
	dbMock.EXPECT().Error().Return(nil)

	events, err := repo.ListFailed(t.Context(), 2)
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Equal(t, "evt-1", events[0].EventID)
	require.Equal(t, "boom", events[0].LastError)
}

type noopStatsLogger struct{}

func (noopStatsLogger) Infof(string, ...any)                      {}
func (noopStatsLogger) Errorf(string, ...any)                     {}
func (noopStatsLogger) Debugf(string, ...any)                     {}
func (noopStatsLogger) Warnf(string, ...any)                      {}
func (noopStatsLogger) Infow(string, ...any)                      {}
func (noopStatsLogger) Errorw(string, ...any)                     {}
func (noopStatsLogger) Debugw(string, ...any)                     {}
func (noopStatsLogger) Warnw(string, ...any)                      {}
func (noopStatsLogger) InfowCtx(context.Context, string, ...any)  {}
func (noopStatsLogger) ErrorwCtx(context.Context, string, ...any) {}
func (noopStatsLogger) WarnwCtx(context.Context, string, ...any)  {}
func (noopStatsLogger) DebugwCtx(context.Context, string, ...any) {}

var _ logger.ContextLogger = noopStatsLogger{}
