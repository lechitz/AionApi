package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/lechitz/AionApi/internal/realtime/core/domain"
)

func TestServicePublishAndSubscribe(t *testing.T) {
	svc := NewService(noopRealtimeLogger{}, 2)

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	stream, cleanup := svc.Subscribe(ctx, 14)
	defer cleanup()

	expected := domain.Event{
		Type:     "record_projection_changed",
		UserID:   14,
		RecordID: 42,
		Action:   "created",
	}

	svc.Publish(t.Context(), expected)

	select {
	case got := <-stream:
		if got.RecordID != expected.RecordID || got.UserID != expected.UserID || got.Action != expected.Action {
			t.Fatalf("unexpected event: %#v", got)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for realtime event")
	}
}

func TestServiceSubscribeFiltersByUser(t *testing.T) {
	svc := NewService(noopRealtimeLogger{}, 1)
	stream, cleanup := svc.Subscribe(t.Context(), 14)
	defer cleanup()

	svc.Publish(t.Context(), domain.Event{Type: "record_projection_changed", UserID: 99, RecordID: 1})

	select {
	case got := <-stream:
		t.Fatalf("unexpected event: %#v", got)
	case <-time.After(50 * time.Millisecond):
	}
}

type noopRealtimeLogger struct{}

func (noopRealtimeLogger) Infof(string, ...any)                      {}
func (noopRealtimeLogger) Errorf(string, ...any)                     {}
func (noopRealtimeLogger) Debugf(string, ...any)                     {}
func (noopRealtimeLogger) Warnf(string, ...any)                      {}
func (noopRealtimeLogger) Infow(string, ...any)                      {}
func (noopRealtimeLogger) Errorw(string, ...any)                     {}
func (noopRealtimeLogger) Debugw(string, ...any)                     {}
func (noopRealtimeLogger) Warnw(string, ...any)                      {}
func (noopRealtimeLogger) InfowCtx(context.Context, string, ...any)  {}
func (noopRealtimeLogger) ErrorwCtx(context.Context, string, ...any) {}
func (noopRealtimeLogger) WarnwCtx(context.Context, string, ...any)  {}
func (noopRealtimeLogger) DebugwCtx(context.Context, string, ...any) {}
