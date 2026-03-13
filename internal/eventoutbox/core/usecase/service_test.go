package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lechitz/AionApi/internal/eventoutbox/core/domain"
	"github.com/lechitz/AionApi/internal/eventoutbox/core/usecase"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
	"github.com/stretchr/testify/require"
)

func TestEnqueue(t *testing.T) {
	t.Run("success writes normalized event", func(t *testing.T) {
		repo := &eventRepositoryStub{}
		svc, ok := usecase.NewService(repo, noopLogger{}).(*usecase.Service)
		require.True(t, ok)

		ctx := context.WithValue(t.Context(), ctxkeys.TraceID, "trace-1")
		ctx = context.WithValue(ctx, ctxkeys.RequestID, "req-1")

		err := svc.Enqueue(ctx, domain.Event{
			AggregateType: "record",
			AggregateID:   "123",
			EventType:     "record.created",
			PayloadJSON:   []byte(`{"record_id":123}`),
		})
		require.NoError(t, err)
		require.Len(t, repo.saved, 1)
		require.NotEmpty(t, repo.saved[0].EventID)
		require.Equal(t, usecase.EventVersionV1, repo.saved[0].EventVersion)
		require.Equal(t, usecase.DefaultSource, repo.saved[0].Source)
		require.Equal(t, usecase.DefaultEventStatus, repo.saved[0].Status)
		require.Equal(t, "trace-1", repo.saved[0].TraceID)
		require.Equal(t, "req-1", repo.saved[0].RequestID)
		require.False(t, repo.saved[0].AvailableAtUTC.IsZero())
	})

	t.Run("invalid command fails", func(t *testing.T) {
		svc, ok := usecase.NewService(&eventRepositoryStub{}, noopLogger{}).(*usecase.Service)
		require.True(t, ok)

		err := svc.Enqueue(t.Context(), domain.Event{})
		require.ErrorIs(t, err, usecase.ErrAggregateTypeRequired)
	})

	t.Run("repository error propagates", func(t *testing.T) {
		svc, ok := usecase.NewService(&eventRepositoryStub{err: errors.New("db fail")}, noopLogger{}).(*usecase.Service)
		require.True(t, ok)

		err := svc.Enqueue(t.Context(), domain.Event{
			AggregateType: "record",
			AggregateID:   "321",
			EventType:     "record.updated",
			PayloadJSON:   []byte(`{"record_id":321}`),
		})
		require.Error(t, err)
	})
}

type eventRepositoryStub struct {
	err   error
	saved []domain.Event
}

func (s *eventRepositoryStub) Save(_ context.Context, event domain.Event) error {
	if s.err != nil {
		return s.err
	}
	s.saved = append(s.saved, event)
	return nil
}

func (s *eventRepositoryStub) ListPending(context.Context, int) ([]domain.Event, error) {
	return nil, nil
}

func (s *eventRepositoryStub) MarkPublished(context.Context, string, time.Time) error {
	return nil
}

func (s *eventRepositoryStub) Reschedule(context.Context, string, time.Time, string) error {
	return nil
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
