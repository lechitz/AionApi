package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lechitz/AionApi/internal/audit/core/domain"
	"github.com/lechitz/AionApi/internal/audit/core/usecase"
	"github.com/lechitz/AionApi/tests/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func newAuditService(t *testing.T) (*usecase.Service, *mocks.MockAuditActionEventRepository) {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mocks.NewMockAuditActionEventRepository(ctrl)
	logger := mocks.NewMockContextLogger(ctrl)
	logger.EXPECT().InfowCtx(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().InfowCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().InfowCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().InfowCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().ErrorwCtx(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().ErrorwCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().ErrorwCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().ErrorwCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	svc, ok := usecase.NewService(repo, logger).(*usecase.Service)
	require.True(t, ok)
	return svc, repo
}

func TestWriteEvent(t *testing.T) {
	svc, repo := newAuditService(t)

	event := domain.AuditActionEvent{
		UserID:       10,
		TraceID:      "trace-1",
		UIActionType: "draft_accept",
		DraftID:      "draft-1",
		Status:       "success",
		Action:       "draft_accept",
		Entity:       "category",
		Operation:    "create",
	}

	t.Run("success normalizes and writes", func(t *testing.T) {
		repo.EXPECT().Save(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, got domain.AuditActionEvent) error {
			require.NotEmpty(t, got.EventID)
			require.False(t, got.TimestampUTC.IsZero())
			require.Equal(t, "aionapi", got.Source)
			return nil
		})

		err := svc.WriteEvent(t.Context(), event)
		require.NoError(t, err)
	})

	t.Run("invalid event fails", func(t *testing.T) {
		err := svc.WriteEvent(t.Context(), domain.AuditActionEvent{})
		require.Error(t, err)
	})

	t.Run("repository error propagates", func(t *testing.T) {
		repo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(errors.New("db fail"))
		err := svc.WriteEvent(t.Context(), event)
		require.Error(t, err)
	})
}

func TestListEvents(t *testing.T) {
	svc, repo := newAuditService(t)

	t.Run("applies defaults", func(t *testing.T) {
		repo.EXPECT().List(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, filter domain.AuditActionEventFilter) ([]domain.AuditActionEvent, error) {
			require.Equal(t, 100, filter.Limit)
			require.Equal(t, 0, filter.Offset)
			return []domain.AuditActionEvent{{EventID: "evt-1", TimestampUTC: time.Now().UTC()}}, nil
		})

		out, err := svc.ListEvents(t.Context(), domain.AuditActionEventFilter{Limit: 0, Offset: -10})
		require.NoError(t, err)
		require.Len(t, out, 1)
	})

	t.Run("repository error propagates", func(t *testing.T) {
		repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, errors.New("query fail"))
		_, err := svc.ListEvents(t.Context(), domain.AuditActionEventFilter{Limit: 5})
		require.Error(t, err)
	})
}
