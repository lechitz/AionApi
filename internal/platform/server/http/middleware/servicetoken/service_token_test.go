package servicetoken_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/server/http/middleware/servicetoken"
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

func TestServiceTokenMiddleware_NoHeader_PassesThrough(t *testing.T) {
	mw := servicetoken.New(&config.Config{AionChat: config.AionChatConfig{ServiceKey: "k"}}, &fakeLogger{})

	var serviceAccount bool
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v, _ := r.Context().Value(ctxkeys.ServiceAccount).(bool)
		serviceAccount = v
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected pass-through status, got %d", w.Code)
	}
	if serviceAccount {
		t.Fatalf("expected service account=false when header not provided")
	}
}

func TestServiceTokenMiddleware_HeaderButNoConfiguredKey_ReturnsUnauthorized(t *testing.T) {
	mw := servicetoken.New(&config.Config{AionChat: config.AionChatConfig{ServiceKey: ""}}, &fakeLogger{})
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	req.Header.Set(servicetoken.HeaderServiceKey, "abc")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestServiceTokenMiddleware_InvalidKey_ReturnsUnauthorized(t *testing.T) {
	mw := servicetoken.New(&config.Config{AionChat: config.AionChatConfig{ServiceKey: "expected"}}, &fakeLogger{})
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	req.Header.Set(servicetoken.HeaderServiceKey, "wrong")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestServiceTokenMiddleware_ValidKey_SetsServiceAccount(t *testing.T) {
	mw := servicetoken.New(&config.Config{AionChat: config.AionChatConfig{ServiceKey: "expected"}}, &fakeLogger{})

	var serviceAccount bool
	var userID uint64
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serviceAccount, _ = r.Context().Value(ctxkeys.ServiceAccount).(bool)
		userID, _ = r.Context().Value(ctxkeys.UserID).(uint64)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/path", nil)
	req.Header.Set(servicetoken.HeaderServiceKey, "expected")
	req.Header.Set(servicetoken.HeaderServiceUser, "123")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !serviceAccount {
		t.Fatalf("expected service account=true for valid key")
	}
	if userID != 123 {
		t.Fatalf("expected propagated user id 123, got %d", userID)
	}
}

func TestServiceTokenMiddleware_ValidKeyInvalidUserID_StillPasses(t *testing.T) {
	mw := servicetoken.New(&config.Config{AionChat: config.AionChatConfig{ServiceKey: "expected"}}, &fakeLogger{})

	var hasUserID bool
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, hasUserID = r.Context().Value(ctxkeys.UserID).(uint64)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/path", nil)
	req.Header.Set(servicetoken.HeaderServiceKey, "expected")
	req.Header.Set(servicetoken.HeaderServiceUser, "invalid-id")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if hasUserID {
		t.Fatalf("expected no user id in context for invalid header value")
	}
}
