package middleware_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lechitz/AionApi/internal/auth/adapter/primary/http/middleware"
	authdomain "github.com/lechitz/AionApi/internal/auth/core/domain"
	input "github.com/lechitz/AionApi/internal/auth/core/ports/input"
	"github.com/lechitz/AionApi/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
)

type fakeAuthService struct {
	validateFn func(ctx context.Context, token string) (uint64, map[string]any, error)
}

var _ input.AuthService = (*fakeAuthService)(nil)

func (f *fakeAuthService) Validate(ctx context.Context, token string) (uint64, map[string]any, error) {
	if f.validateFn != nil {
		return f.validateFn(ctx, token)
	}
	return 0, nil, nil
}

func (f *fakeAuthService) Login(context.Context, string, string) (authdomain.AuthenticatedUser, string, string, error) {
	return authdomain.AuthenticatedUser{}, "", "", nil
}

func (f *fakeAuthService) Logout(context.Context, uint64) error { return nil }

func (f *fakeAuthService) RefreshTokenRenewal(context.Context, string) (string, string, error) {
	return "", "", nil
}

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

func TestAuthMiddleware_ServiceAccountBypass(t *testing.T) {
	m := middleware.New(&fakeAuthService{}, &fakeLogger{})
	h := m.Auth(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	ctx := context.WithValue(req.Context(), ctxkeys.ServiceAccount, true)
	ctx = context.WithValue(ctx, ctxkeys.UserID, uint64(9))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected bypass success, got %d", w.Code)
	}
}

func TestAuthMiddleware_MissingToken_ReturnsError(t *testing.T) {
	m := middleware.New(&fakeAuthService{}, &fakeLogger{})
	h := m.Auth(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for missing token path, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidToken_ReturnsUnauthorized(t *testing.T) {
	m := middleware.New(&fakeAuthService{
		validateFn: func(context.Context, string) (uint64, map[string]any, error) {
			return 0, nil, sharederrors.ErrUnauthorized("invalid")
		},
	}, &fakeLogger{})
	h := m.Auth(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer bad-token")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for invalid token, got %d", w.Code)
	}
}

func TestAuthMiddleware_ValidateErrorRaw_ReturnsServerError(t *testing.T) {
	m := middleware.New(&fakeAuthService{
		validateFn: func(context.Context, string) (uint64, map[string]any, error) {
			return 0, nil, errors.New("boom")
		},
	}, &fakeLogger{})
	h := m.Auth(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer token")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for generic validation error, got %d", w.Code)
	}
}

func TestAuthMiddleware_ValidBearerSetsContext(t *testing.T) {
	m := middleware.New(&fakeAuthService{
		validateFn: func(_ context.Context, token string) (uint64, map[string]any, error) {
			if token != "good-token" {
				t.Fatalf("unexpected token: %s", token)
			}
			return 42, map[string]any{"role": "user"}, nil
		},
	}, &fakeLogger{})

	var gotUserID uint64
	var gotToken string
	var gotClaims map[string]any
	h := m.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUserID, _ = r.Context().Value(ctxkeys.UserID).(uint64)
		gotToken, _ = r.Context().Value(ctxkeys.Token).(string)
		gotClaims, _ = r.Context().Value(ctxkeys.Claims).(map[string]any)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer good-token")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if gotUserID != 42 || gotToken != "good-token" || gotClaims["role"] != "user" {
		t.Fatalf("unexpected context values: userID=%d token=%s claims=%v", gotUserID, gotToken, gotClaims)
	}
}

func TestAuthMiddleware_ValidCookieToken(t *testing.T) {
	m := middleware.New(&fakeAuthService{
		validateFn: func(_ context.Context, token string) (uint64, map[string]any, error) {
			if token != "cookie-token" {
				t.Fatalf("unexpected token: %s", token)
			}
			return 7, nil, nil
		},
	}, &fakeLogger{})

	var gotUserID uint64
	h := m.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUserID, _ = r.Context().Value(ctxkeys.UserID).(uint64)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: commonkeys.AuthTokenCookieName, Value: "cookie-token"})
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if gotUserID != 7 {
		t.Fatalf("expected user id 7, got %d", gotUserID)
	}
}
