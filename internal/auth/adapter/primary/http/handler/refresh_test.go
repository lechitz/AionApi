package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	handlerpkg "github.com/lechitz/AionApi/internal/auth/adapter/primary/http/handler"
	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/user/core/domain"
)

// minimalAuthService is a test stub implementing the expected AuthService interface.
// It allows controlling the return values of RefreshTokenRenewal.
// We implement all methods of input.AuthService as simple stubs so the test can construct
// a Handler without requiring the full production service.
type minimalAuthService struct {
	respAccess  string
	respRefresh string
	respErr     error
	// capture the received token for assertion if needed.
	receivedToken string
}

func (m *minimalAuthService) RefreshTokenRenewal(_ context.Context, refreshToken string) (string, string, error) {
	m.receivedToken = refreshToken
	return m.respAccess, m.respRefresh, m.respErr
}

// ---- stubs to satisfy the rest of input.AuthService.
func (m *minimalAuthService) Login(_ context.Context, _ string, _ string) (domain.User, string, string, error) {
	// not used by these tests.
	return domain.User{}, "", "", nil
}

func (m *minimalAuthService) Validate(_ context.Context, _ string) (uint64, map[string]any, error) {
	return 0, nil, nil
}

func (m *minimalAuthService) Logout(_ context.Context, _ uint64) error {
	return nil
}

// simpleLogger satisfies logger.ContextLogger used by handler.New; it discards logs.
// Implement all methods from the interface to keep tests focused on handler behavior.
type simpleLogger struct{}

func (s *simpleLogger) Infof(_ string, _ ...any)  {}
func (s *simpleLogger) Errorf(_ string, _ ...any) {}
func (s *simpleLogger) Debugf(_ string, _ ...any) {}
func (s *simpleLogger) Warnf(_ string, _ ...any)  {}

func (s *simpleLogger) Infow(_ string, _ ...any)  {}
func (s *simpleLogger) Errorw(_ string, _ ...any) {}
func (s *simpleLogger) Debugw(_ string, _ ...any) {}
func (s *simpleLogger) Warnw(_ string, _ ...any)  {}

func (s *simpleLogger) InfowCtx(_ context.Context, _ string, _ ...any)  {}
func (s *simpleLogger) ErrorwCtx(_ context.Context, _ string, _ ...any) {}
func (s *simpleLogger) WarnwCtx(_ context.Context, _ string, _ ...any)  {}
func (s *simpleLogger) DebugwCtx(_ context.Context, _ string, _ ...any) {}

func TestRefresh_CookieMissing(t *testing.T) {
	// Setup a simple in-memory tracer provider so handler can start spans and we can inspect them.
	tt := NewTestTracer(t)
	defer func() {
		if err := tt.Shutdown(context.Background()); err != nil {
			t.Fatalf("failed to shutdown tracer: %v", err)
		}
	}()

	cfg := &config.Config{}
	svc := &minimalAuthService{}
	h := handlerpkg.New(svc, cfg, &simpleLogger{})

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()

	h.Refresh(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized when cookie missing, got %d", res.StatusCode)
	}

	// ensure service was not called
	if svc.receivedToken != "" {
		t.Fatalf("service should not be called when cookie is missing")
	}

	// verify no spans leaked tokens
	spans := tt.Spans()
	if len(spans) == 0 {
		// still acceptable; nothing to assert further.
		return
	}
	for _, sp := range spans {
		for _, attr := range sp.Attributes {
			if v := attr.Value.AsString(); strings.Contains(v, "token") {
				t.Fatalf("span attribute leaked token-like value: %q", v)
			}
		}
		for _, ev := range sp.Events {
			for _, attr := range ev.Attributes {
				if v := attr.Value.AsString(); strings.Contains(v, "token") {
					t.Fatalf("span event leaked token-like value: %q", v)
				}
			}
		}
	}
}

func TestRefresh_CookiePresent_Success(t *testing.T) {
	// simple in-memory tracer provider
	tt := NewTestTracer(t)
	defer func() {
		if err := tt.Shutdown(context.Background()); err != nil {
			t.Fatalf("failed to shutdown tracer: %v", err)
		}
	}()

	cfg := &config.Config{}
	// prepare service to return tokens
	svc := &minimalAuthService{respAccess: "access-123", respRefresh: "refresh-456", respErr: nil}
	h := handlerpkg.New(svc, cfg, &simpleLogger{})

	// build request with cookie
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	token := "r-token-xyz"
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: token})
	w := httptest.NewRecorder()

	h.Refresh(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK on success, got %d", res.StatusCode)
	}

	// ensure service received the correct token
	if svc.receivedToken != token {
		t.Fatalf("service did not receive expected refresh token, got %q", svc.receivedToken)
	}

	// We assert that cookies were set on response for both auth and refresh.
	assertCookiesSet(t, res)

	// verify spans do not contain the raw token
	assertNoTokenInSpans(t, tt, token)
}

func TestRefresh_ServiceError(t *testing.T) {
	// in-memory tracer provider
	tt := NewTestTracer(t)
	defer func() {
		if err := tt.Shutdown(context.Background()); err != nil {
			t.Fatalf("failed to shutdown tracer: %v", err)
		}
	}()

	cfg := &config.Config{}
	svc := &minimalAuthService{respErr: errors.New("service failure")}
	h := handlerpkg.New(svc, cfg, &simpleLogger{})

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	token := "some"
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: token})
	w := httptest.NewRecorder()

	h.Refresh(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized on service error, got %d", res.StatusCode)
	}
	// ensure service was called
	if svc.receivedToken != token {
		t.Fatalf("service should have been called with cookie value")
	}

	// ensure spans do not include the token
	assertNoTokenInSpans(t, tt, token)
}

// assertCookiesSet asserts both auth and refresh cookies are present in the response.
func assertCookiesSet(t *testing.T, res *http.Response) {
	t.Helper()
	foundAuth := false
	foundRefresh := false
	for _, c := range res.Cookies() {
		if c.Name == "refresh_token" {
			foundRefresh = true
		}
		if c.Name == "auth_token" {
			foundAuth = true
		}
	}
	if !foundRefresh || !foundAuth {
		t.Fatalf("expected both auth and refresh cookies to be set on success; got auth=%v refresh=%v", foundAuth, foundRefresh)
	}
}

// assertNoTokenInSpans inspects exported spans to ensure the raw token does not appear.
func assertNoTokenInSpans(t *testing.T, tt *TestTracer, token string) {
	t.Helper()
	spans := tt.Spans()
	if len(spans) == 0 {
		t.Fatal("expected spans to be recorded")
	}
	for _, sp := range spans {
		for _, attr := range sp.Attributes {
			if attr.Value.AsString() == token {
				t.Fatalf("span attribute leaked raw token: %q", token)
			}
			if strings.Contains(string(attr.Key), "token") && attr.Value.AsString() == token {
				t.Fatalf("span attribute key/value leaked token: %s=%q", attr.Key, token)
			}
		}
		for _, ev := range sp.Events {
			for _, attr := range ev.Attributes {
				if attr.Value.AsString() == token {
					t.Fatalf("span event attribute leaked raw token: %q", token)
				}
			}
		}
	}
}
