package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lechitz/aion-api/internal/audit/adapter/primary/http/handler"
	"github.com/lechitz/aion-api/internal/audit/core/domain"
	"github.com/lechitz/aion-api/internal/platform/config"
	"github.com/lechitz/aion-api/internal/shared/constants/ctxkeys"
	"github.com/lechitz/aion-api/tests/mocks"
	"github.com/lechitz/aion-api/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type envelope struct {
	Result json.RawMessage `json:"result"`
	Code   int             `json:"code"`
}

type eventsResult struct {
	Items []struct {
		EventID string `json:"event_id"`
	} `json:"items"`
	Count int `json:"count"`
}

func newAuditHandler(t *testing.T) (*handler.Handler, *mocks.MockService) {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mocks.NewMockService(ctrl)
	lg := mocks.NewMockContextLogger(ctrl)
	setup.ExpectLoggerDefaultBehavior(lg)
	lg.EXPECT().Errorw(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	lg.EXPECT().
		InfowCtx(
			gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		).
		AnyTimes()
	lg.EXPECT().
		ErrorwCtx(
			gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		).
		AnyTimes()

	return handler.New(svc, &config.Config{}, lg), svc
}

func TestListEvents_Success(t *testing.T) {
	h, svc := newAuditHandler(t)
	from := "2026-02-25T00:00:00Z"
	to := "2026-02-25T23:59:59Z"

	svc.EXPECT().ListEvents(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, filter domain.AuditActionEventFilter) ([]domain.AuditActionEvent, error) {
			require.NotNil(t, filter.UserID)
			require.Equal(t, uint64(7), *filter.UserID)
			require.Equal(t, "trace-1", filter.TraceID)
			require.Equal(t, "draft-1", filter.DraftID)
			require.Equal(t, []string{"failed", "blocked"}, filter.Statuses)
			require.Equal(t, 20, filter.Limit)
			require.Equal(t, 5, filter.Offset)
			require.NotNil(t, filter.FromUTC)
			require.NotNil(t, filter.ToUTC)
			return []domain.AuditActionEvent{{EventID: "evt-1", TimestampUTC: time.Now().UTC(), UserID: 7}}, nil
		},
	)

	req := httptest.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		"/audit/events?trace_id=trace-1&draft_id=draft-1&status=failed&status=blocked&limit=20&offset=5&from_utc="+from+"&to_utc="+to,
		nil,
	)
	req = req.WithContext(context.WithValue(req.Context(), ctxkeys.UserID, uint64(7)))
	rec := httptest.NewRecorder()

	h.ListEvents(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var env envelope
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &env))
	require.Equal(t, http.StatusOK, env.Code)

	var result eventsResult
	require.NoError(t, json.Unmarshal(env.Result, &result))
	require.Equal(t, 1, result.Count)
	require.Equal(t, "evt-1", result.Items[0].EventID)
}

func TestListEvents_InvalidLimit(t *testing.T) {
	h, _ := newAuditHandler(t)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/audit/events?limit=abc", nil)
	req = req.WithContext(context.WithValue(req.Context(), ctxkeys.UserID, uint64(7)))
	rec := httptest.NewRecorder()

	h.ListEvents(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListEvents_UserIDParamNotAllowed(t *testing.T) {
	h, _ := newAuditHandler(t)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/audit/events?user_id=999", nil)
	req = req.WithContext(context.WithValue(req.Context(), ctxkeys.UserID, uint64(7)))
	rec := httptest.NewRecorder()

	h.ListEvents(rec, req)
	require.Equal(t, http.StatusForbidden, rec.Code)
}

func TestListEvents_AdminCanQueryOtherUser(t *testing.T) {
	h, svc := newAuditHandler(t)

	svc.EXPECT().ListEvents(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, filter domain.AuditActionEventFilter) ([]domain.AuditActionEvent, error) {
			require.NotNil(t, filter.UserID)
			require.Equal(t, uint64(999), *filter.UserID)
			return []domain.AuditActionEvent{}, nil
		},
	)

	ctx := context.WithValue(t.Context(), ctxkeys.UserID, uint64(7))
	ctx = context.WithValue(ctx, ctxkeys.Claims, map[string]any{"roles": []any{"admin"}})
	req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/audit/events?user_id=999", nil)
	rec := httptest.NewRecorder()

	h.ListEvents(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestListEvents_ServiceError(t *testing.T) {
	h, svc := newAuditHandler(t)
	svc.EXPECT().ListEvents(gomock.Any(), gomock.Any()).Return(nil, errors.New("db down"))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/audit/events", nil)
	req = req.WithContext(context.WithValue(req.Context(), ctxkeys.UserID, uint64(7)))
	rec := httptest.NewRecorder()

	h.ListEvents(rec, req)
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}
