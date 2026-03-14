package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lechitz/AionApi/internal/platform/config"
	handler "github.com/lechitz/AionApi/internal/platform/server/http/generic/handler"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
)

type fakeLogger struct{}

func (f *fakeLogger) Infof(string, ...any)                      {}
func (f *fakeLogger) Errorf(string, ...any)                     {}
func (f *fakeLogger) Debugf(string, ...any)                     {}
func (f *fakeLogger) Warnf(string, ...any)                      {}
func (f *fakeLogger) Infow(string, ...any)                      {}
func (f *fakeLogger) Errorw(string, ...any)                     {}
func (f *fakeLogger) Debugw(string, ...any)                     {}
func (f *fakeLogger) Warnw(string, ...any)                      {}
func (f *fakeLogger) InfowCtx(context.Context, string, ...any)  {}
func (f *fakeLogger) ErrorwCtx(context.Context, string, ...any) {}
func (f *fakeLogger) WarnwCtx(context.Context, string, ...any)  {}
func (f *fakeLogger) DebugwCtx(context.Context, string, ...any) {}

func newGenericHandler() *handler.Handler {
	return handler.New(&fakeLogger{}, config.GeneralConfig{
		Name:    "AionApi",
		Env:     "test",
		Version: "v0",
	})
}

func TestHealthCheck_Success(t *testing.T) {
	h := newGenericHandler()
	req := httptest.NewRequestWithContext(
		context.WithValue(t.Context(), ctxkeys.RequestID, "r-1"),
		http.MethodGet,
		"/health",
		nil,
	)
	w := httptest.NewRecorder()

	h.HealthCheck(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body map[string]any
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	code, ok := body["code"].(float64)
	if !ok {
		t.Fatalf("expected numeric code, got %#v", body["code"])
	}
	if int(code) != http.StatusOK {
		t.Fatalf("unexpected code payload: %v", body["code"])
	}
	result, ok := body["result"].(map[string]any)
	if !ok {
		t.Fatalf("expected object result, got %#v", body["result"])
	}
	if result["status"] != "healthy" || result["name"] != "AionApi" {
		t.Fatalf("unexpected health payload: %+v", result)
	}
}

func TestHealthCheck_MethodNotAllowed(t *testing.T) {
	h := newGenericHandler()
	req := httptest.NewRequestWithContext(
		context.WithValue(t.Context(), ctxkeys.RequestID, "r-2"),
		http.MethodPost,
		"/health",
		nil,
	)
	w := httptest.NewRecorder()

	h.HealthCheck(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}

func TestMethodNotAllowedHandler(t *testing.T) {
	h := newGenericHandler()
	req := httptest.NewRequestWithContext(
		context.WithValue(t.Context(), ctxkeys.RequestID, "r-3"),
		http.MethodPost,
		"/x",
		nil,
	)
	w := httptest.NewRecorder()

	h.MethodNotAllowedHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}

func TestNotFoundHandler(t *testing.T) {
	h := newGenericHandler()
	req := httptest.NewRequestWithContext(
		context.WithValue(t.Context(), ctxkeys.RequestID, "r-4"),
		http.MethodGet,
		"/not-found",
		nil,
	)
	w := httptest.NewRecorder()

	h.NotFoundHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestErrorHandler(t *testing.T) {
	h := newGenericHandler()
	req := httptest.NewRequestWithContext(
		context.WithValue(t.Context(), ctxkeys.RequestID, "r-5"),
		http.MethodGet,
		"/err",
		nil,
	)
	w := httptest.NewRecorder()

	h.ErrorHandler(w, req, errors.New("boom"))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestRecoveryHandler(t *testing.T) {
	h := newGenericHandler()
	req := httptest.NewRequestWithContext(
		context.WithValue(t.Context(), ctxkeys.RequestID, "r-6"),
		http.MethodGet,
		"/panic",
		nil,
	)
	w := httptest.NewRecorder()

	h.RecoveryHandler(w, req, "panic payload", "err-id")

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
