package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	handlerpkg "github.com/lechitz/AionApi/internal/auth/adapter/primary/http/handler"
	authdomain "github.com/lechitz/AionApi/internal/auth/core/domain"
	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/shared/constants/claimskeys"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
	"github.com/lechitz/AionApi/tests/mocks"
	"github.com/lechitz/AionApi/tests/setup"
	"go.uber.org/mock/gomock"
)

type authServiceStub struct {
	loginFn    func(ctx context.Context, username, password string) (authdomain.AuthenticatedUser, string, string, error)
	validateFn func(ctx context.Context, token string) (uint64, map[string]any, error)
	logoutFn   func(ctx context.Context, userID uint64) error
	refreshFn  func(ctx context.Context, refreshToken string) (string, string, error)
}

func (s *authServiceStub) Login(ctx context.Context, username, password string) (authdomain.AuthenticatedUser, string, string, error) {
	if s.loginFn != nil {
		return s.loginFn(ctx, username, password)
	}
	return authdomain.AuthenticatedUser{}, "", "", nil
}

func (s *authServiceStub) Validate(ctx context.Context, token string) (uint64, map[string]any, error) {
	if s.validateFn != nil {
		return s.validateFn(ctx, token)
	}
	return 0, nil, nil
}

func (s *authServiceStub) Logout(ctx context.Context, userID uint64) error {
	if s.logoutFn != nil {
		return s.logoutFn(ctx, userID)
	}
	return nil
}

func (s *authServiceStub) RefreshTokenRenewal(ctx context.Context, refreshToken string) (string, string, error) {
	if s.refreshFn != nil {
		return s.refreshFn(ctx, refreshToken)
	}
	return "", "", nil
}

func newTestHandler(t *testing.T, svc *authServiceStub) *handlerpkg.Handler {
	t.Helper()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	log := mocks.NewMockContextLogger(ctrl)
	setup.ExpectLoggerDefaultBehavior(log)
	log.EXPECT().Infow(gomock.Any(), gomock.Any()).AnyTimes()
	log.EXPECT().Errorw(gomock.Any(), gomock.Any()).AnyTimes()
	log.EXPECT().Warnw(gomock.Any(), gomock.Any()).AnyTimes()
	log.EXPECT().Debugw(gomock.Any(), gomock.Any()).AnyTimes()

	cfg := &config.Config{
		Cookie: config.CookieConfig{
			Path:   "/",
			Secure: false,
		},
	}
	return handlerpkg.New(svc, cfg, log)
}

func TestLogin_InvalidJSON(t *testing.T) {
	h := newTestHandler(t, &authServiceStub{})

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader("{"))
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestLogin_ValidationError(t *testing.T) {
	h := newTestHandler(t, &authServiceStub{})

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"username":"ab","password":"123"}`))
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestLogin_ServiceError(t *testing.T) {
	h := newTestHandler(t, &authServiceStub{
		loginFn: func(context.Context, string, string) (authdomain.AuthenticatedUser, string, string, error) {
			return authdomain.AuthenticatedUser{}, "", "", errors.New("invalid credentials")
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"username":"validuser","password":"validpassword"}`))
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestLogin_Success(t *testing.T) {
	h := newTestHandler(t, &authServiceStub{
		loginFn: func(_ context.Context, username, password string) (authdomain.AuthenticatedUser, string, string, error) {
			if username != "validuser" || password != "validpassword" {
				t.Fatalf("unexpected credentials: %s / %s", username, password)
			}
			return authdomain.AuthenticatedUser{
				ID:   42,
				Name: "Lechitz",
			}, "access-token", "refresh-token", nil
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"username":"validuser","password":"validpassword"}`))
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	res := w.Result()
	foundAuth := false
	foundRefresh := false
	for _, c := range res.Cookies() {
		if c.Name == commonkeys.AuthTokenCookieName {
			foundAuth = true
		}
		if c.Name == "refresh_token" {
			foundRefresh = true
		}
	}
	if !foundAuth || !foundRefresh {
		t.Fatalf("expected auth and refresh cookies on login success")
	}

	var body struct {
		Code   int `json:"code"`
		Result struct {
			ID    uint64 `json:"id"`
			Name  string `json:"name"`
			Token string `json:"token"`
		} `json:"result"`
	}
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body.Code != http.StatusOK || body.Result.ID != 42 || body.Result.Name != "Lechitz" {
		t.Fatalf("unexpected response payload: %+v", body)
	}
}

func TestLogout_MissingUserIDInContext(t *testing.T) {
	h := newTestHandler(t, &authServiceStub{})

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	w := httptest.NewRecorder()

	h.Logout(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestLogout_ServiceError(t *testing.T) {
	h := newTestHandler(t, &authServiceStub{
		logoutFn: func(_ context.Context, userID uint64) error {
			if userID != 99 {
				t.Fatalf("unexpected userID: %d", userID)
			}
			return errors.New("db unavailable")
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	ctx := context.WithValue(req.Context(), ctxkeys.UserID, uint64(99))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.Logout(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestLogout_Success(t *testing.T) {
	h := newTestHandler(t, &authServiceStub{
		logoutFn: func(_ context.Context, userID uint64) error {
			if userID != 77 {
				t.Fatalf("unexpected userID: %d", userID)
			}
			return nil
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	ctx := context.WithValue(req.Context(), ctxkeys.UserID, uint64(77))
	ctx = context.WithValue(ctx, ctxkeys.Token, "abcdefghij12345token")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.Logout(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}

	res := w.Result()
	foundAuth := false
	foundRefresh := false
	for _, c := range res.Cookies() {
		if c.Name == commonkeys.AuthTokenCookieName && c.MaxAge < 0 {
			foundAuth = true
		}
		if c.Name == "refresh_token" && c.MaxAge < 0 {
			foundRefresh = true
		}
	}
	if !foundAuth || !foundRefresh {
		t.Fatalf("expected cleared auth and refresh cookies on logout")
	}
}

func TestSession_MissingAuthCookie(t *testing.T) {
	h := newTestHandler(t, &authServiceStub{})

	req := httptest.NewRequest(http.MethodGet, "/auth/session", nil)
	w := httptest.NewRecorder()

	h.Session(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestSession_ValidateError(t *testing.T) {
	h := newTestHandler(t, &authServiceStub{
		validateFn: func(context.Context, string) (uint64, map[string]any, error) {
			return 0, nil, errors.New("token invalid")
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/auth/session", nil)
	req.AddCookie(&http.Cookie{Name: commonkeys.AuthTokenCookieName, Value: "access-token"})
	w := httptest.NewRecorder()

	h.Session(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestSession_SuccessWithClaimsNormalization(t *testing.T) {
	expectedExp := time.Unix(1735689600, 0).UTC()
	h := newTestHandler(t, &authServiceStub{
		validateFn: func(_ context.Context, token string) (uint64, map[string]any, error) {
			if token != "access-token" {
				t.Fatalf("unexpected token: %s", token)
			}
			return 101, map[string]any{
				claimskeys.Username: "john",
				claimskeys.Email:    "john@example.com",
				claimskeys.Name:     "John",
				claimskeys.Roles:    []any{"user", "admin"},
				claimskeys.Exp:      "1735689600",
			}, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/auth/session", nil)
	req.AddCookie(&http.Cookie{Name: commonkeys.AuthTokenCookieName, Value: "access-token"})
	w := httptest.NewRecorder()

	h.Session(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body struct {
		Code   int `json:"code"`
		Result struct {
			Authenticated bool       `json:"authenticated"`
			UserID        uint64     `json:"user_id"`
			Username      string     `json:"username"`
			Email         string     `json:"email"`
			Name          string     `json:"name"`
			Roles         []string   `json:"roles"`
			ExpiresAt     *time.Time `json:"expires_at"`
		} `json:"result"`
	}
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if body.Code != http.StatusOK || !body.Result.Authenticated || body.Result.UserID != 101 {
		t.Fatalf("unexpected response status/body: %+v", body)
	}
	if len(body.Result.Roles) != 2 || body.Result.Roles[0] != "user" || body.Result.Roles[1] != "admin" {
		t.Fatalf("unexpected roles: %+v", body.Result.Roles)
	}
	if body.Result.ExpiresAt == nil || !body.Result.ExpiresAt.Equal(expectedExp) {
		t.Fatalf("unexpected expires_at: %v expected %v", body.Result.ExpiresAt, expectedExp)
	}
}

func TestLogin_ValidationRequiredFields(t *testing.T) {
	h := newTestHandler(t, &authServiceStub{})

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"username":"   ","password":"   "}`))
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestSession_ClaimVariant_RolesFromCommonKeyAndExpFloat64(t *testing.T) {
	body := runSessionWithClaims(t, map[string]any{
		commonkeys.Roles: []string{"user"},
		claimskeys.Exp:   float64(1735689600),
	})
	if strings.Join(body.Result.Roles, ",") != "user" {
		t.Fatalf("unexpected roles: %+v", body.Result.Roles)
	}
	if body.Result.ExpiresAt == nil || body.Result.ExpiresAt.Unix() != 1735689600 {
		t.Fatalf("unexpected expires_at: %v", body.Result.ExpiresAt)
	}
}

func TestSession_ClaimVariant_RolesAsStringAndExpInt64(t *testing.T) {
	body := runSessionWithClaims(t, map[string]any{
		claimskeys.Roles: "admin",
		claimskeys.Exp:   int64(1735689601),
	})
	if strings.Join(body.Result.Roles, ",") != "admin" {
		t.Fatalf("unexpected roles: %+v", body.Result.Roles)
	}
	if body.Result.ExpiresAt == nil || body.Result.ExpiresAt.Unix() != 1735689601 {
		t.Fatalf("unexpected expires_at: %v", body.Result.ExpiresAt)
	}
}

func TestSession_ClaimVariant_RolesInvalidAndExpInt(t *testing.T) {
	body := runSessionWithClaims(t, map[string]any{
		claimskeys.Roles: 99,
		claimskeys.Exp:   int(1735689602),
	})
	if len(body.Result.Roles) != 0 {
		t.Fatalf("unexpected roles: %+v", body.Result.Roles)
	}
	if body.Result.ExpiresAt == nil || body.Result.ExpiresAt.Unix() != 1735689602 {
		t.Fatalf("unexpected expires_at: %v", body.Result.ExpiresAt)
	}
}

func TestSession_ClaimVariant_InvalidExpStringReturnsNil(t *testing.T) {
	body := runSessionWithClaims(t, map[string]any{
		claimskeys.Roles: []any{"user"},
		claimskeys.Exp:   "173x",
	})
	if strings.Join(body.Result.Roles, ",") != "user" {
		t.Fatalf("unexpected roles: %+v", body.Result.Roles)
	}
	if body.Result.ExpiresAt != nil {
		t.Fatalf("expected nil expires_at, got %v", body.Result.ExpiresAt)
	}
}

func TestSession_EmptyAccessTokenCookie(t *testing.T) {
	h := newTestHandler(t, &authServiceStub{})

	req := httptest.NewRequest(http.MethodGet, "/auth/session", nil)
	req.AddCookie(&http.Cookie{Name: commonkeys.AuthTokenCookieName, Value: ""})
	w := httptest.NewRecorder()

	h.Session(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func runSessionWithClaims(t *testing.T, claims map[string]any) struct {
	Result struct {
		Roles     []string   `json:"roles"`
		ExpiresAt *time.Time `json:"expires_at"`
	} `json:"result"`
} {
	t.Helper()
	h := newTestHandler(t, &authServiceStub{
		validateFn: func(_ context.Context, token string) (uint64, map[string]any, error) {
			if token != "access-token" {
				t.Fatalf("unexpected token: %s", token)
			}
			return 55, claims, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/auth/session", nil)
	req.AddCookie(&http.Cookie{Name: commonkeys.AuthTokenCookieName, Value: "access-token"})
	w := httptest.NewRecorder()
	h.Session(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body struct {
		Result struct {
			Roles     []string   `json:"roles"`
			ExpiresAt *time.Time `json:"expires_at"`
		} `json:"result"`
	}
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	return body
}
