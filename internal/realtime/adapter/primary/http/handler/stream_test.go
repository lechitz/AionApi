//nolint:testpackage // Tests exercise package-private SSE details and constants.
package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/lechitz/aion-api/internal/platform/config"
	"github.com/lechitz/aion-api/internal/realtime/core/domain"
	realtimeUsecase "github.com/lechitz/aion-api/internal/realtime/core/usecase"
	"github.com/lechitz/aion-api/internal/shared/constants/ctxkeys"
)

func TestStreamWritesSSEEvent(t *testing.T) {
	service := realtimeUsecase.NewService(noopRealtimeHandlerLogger{}, 4)
	handler := New(service, &config.Config{
		Realtime: config.RealtimeConfig{
			HeartbeatInterval: time.Minute,
		},
	}, noopRealtimeHandlerLogger{})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/events/stream", nil)
	ctx, cancel := context.WithCancel(req.Context())
	req = req.WithContext(context.WithValue(ctx, ctxkeys.UserID, uint64(14)))
	rec := httptest.NewRecorder()

	done := make(chan struct{})
	go func() {
		defer close(done)
		handler.Stream(rec, req)
	}()

	time.Sleep(20 * time.Millisecond)
	service.Publish(t.Context(), domain.Event{
		Type:           "record_projection_changed",
		UserID:         14,
		RecordID:       42,
		Action:         "created",
		ProjectedAtUTC: time.Now().UTC(),
	})
	time.Sleep(20 * time.Millisecond)
	cancel()
	<-done

	body := rec.Body.String()
	if !strings.Contains(body, "event: connected") {
		t.Fatalf("expected connected event, got %q", body)
	}
	if !strings.Contains(body, "event: record_projection_changed") {
		t.Fatalf("expected record projection event, got %q", body)
	}
	if !strings.Contains(body, "\"recordId\":42") {
		t.Fatalf("expected record id payload, got %q", body)
	}
	if got := rec.Header().Get(headerContentType); got != contentTypeEventStream {
		t.Fatalf("expected %s header, got %q", headerContentType, got)
	}
}

type noopRealtimeHandlerLogger struct{}

func (noopRealtimeHandlerLogger) Infof(string, ...any)                      {}
func (noopRealtimeHandlerLogger) Errorf(string, ...any)                     {}
func (noopRealtimeHandlerLogger) Debugf(string, ...any)                     {}
func (noopRealtimeHandlerLogger) Warnf(string, ...any)                      {}
func (noopRealtimeHandlerLogger) Infow(string, ...any)                      {}
func (noopRealtimeHandlerLogger) Errorw(string, ...any)                     {}
func (noopRealtimeHandlerLogger) Debugw(string, ...any)                     {}
func (noopRealtimeHandlerLogger) Warnw(string, ...any)                      {}
func (noopRealtimeHandlerLogger) InfowCtx(context.Context, string, ...any)  {}
func (noopRealtimeHandlerLogger) ErrorwCtx(context.Context, string, ...any) {}
func (noopRealtimeHandlerLogger) WarnwCtx(context.Context, string, ...any)  {}
func (noopRealtimeHandlerLogger) DebugwCtx(context.Context, string, ...any) {}
