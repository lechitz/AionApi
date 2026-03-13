package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lechitz/AionApi/internal/eventoutbox/adapter/secondary/db/model"
	repository "github.com/lechitz/AionApi/internal/eventoutbox/adapter/secondary/db/repository"
	"github.com/lechitz/AionApi/internal/eventoutbox/core/domain"
	"github.com/lechitz/AionApi/internal/platform/ports/output/db"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/tests/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func newOutboxRepo(t *testing.T) (*repository.EventRepository, *mocks.MockDB) {
	t.Helper()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	dbMock := mocks.NewMockDB(ctrl)

	return repository.NewEventRepository(dbMock, noopLogger{}), dbMock
}

func TestEventRepository_Save(t *testing.T) {
	repo, dbMock := newOutboxRepo(t)
	event := domain.Event{
		EventID:        "evt-1",
		AggregateType:  "record",
		AggregateID:    "1",
		EventType:      "record.created",
		EventVersion:   "v1",
		Source:         "aionapi",
		TraceID:        "trace-1",
		RequestID:      "req-1",
		Status:         "pending",
		AvailableAtUTC: time.Now().UTC(),
		PayloadJSON:    []byte(`{"record_id":1}`),
		CreatedAt:      time.Now().UTC(),
	}

	t.Run("success", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Create(gomock.Any()).DoAndReturn(func(v any) db.DB {
			row, ok := v.(*model.EventDB)
			require.True(t, ok)
			require.Equal(t, "evt-1", row.EventID)
			require.Equal(t, "record", row.AggregateType)
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

type noopLogger struct{}

func (noopLogger) Infof(string, ...any)                      {}
func (noopLogger) Errorf(string, ...any)                     {}
func (noopLogger) Debugf(string, ...any)                     {}
func (noopLogger) Warnf(string, ...any)                      {}
func (noopLogger) Infow(string, ...any)                      {}
func (noopLogger) Errorw(string, ...any)                     {}
func (noopLogger) Debugw(string, ...any)                     {}
func (noopLogger) Warnw(string, ...any)                      {}
func (noopLogger) InfowCtx(context.Context, string, ...any)  {}
func (noopLogger) ErrorwCtx(context.Context, string, ...any) {}
func (noopLogger) WarnwCtx(context.Context, string, ...any)  {}
func (noopLogger) DebugwCtx(context.Context, string, ...any) {}

var _ logger.ContextLogger = noopLogger{}
