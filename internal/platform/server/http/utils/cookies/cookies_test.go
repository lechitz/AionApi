package cookies_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lechitz/AionApi/internal/platform/config"
	cookies "github.com/lechitz/AionApi/internal/platform/server/http/utils/cookies"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
)

func TestSetAuthCookie_AndClearAuthCookie(t *testing.T) {
	cfg := config.CookieConfig{
		Domain:   "example.com",
		Path:     "/",
		SameSite: "Strict",
		Secure:   true,
		MaxAge:   3600,
	}

	w := httptest.NewRecorder()
	cookies.SetAuthCookie(w, "access-token", cfg)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	gotCookies := res.Cookies()
	if len(gotCookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(gotCookies))
	}
	c := gotCookies[0]
	if c.Name != commonkeys.AuthTokenCookieName || c.Value != "access-token" {
		t.Fatalf("unexpected auth cookie: %+v", c)
	}
	if !c.HttpOnly || !c.Secure || c.SameSite != http.SameSiteStrictMode {
		t.Fatalf("unexpected auth cookie flags: %+v", c)
	}
	if c.MaxAge != 3600 {
		t.Fatalf("expected max-age 3600, got %d", c.MaxAge)
	}

	w = httptest.NewRecorder()
	cookies.ClearAuthCookie(w, cfg)
	res = w.Result()
	defer func() { _ = res.Body.Close() }()
	c = res.Cookies()[0]
	if c.Name != commonkeys.AuthTokenCookieName || c.MaxAge != -1 {
		t.Fatalf("unexpected cleared auth cookie: %+v", c)
	}
}

func TestSetRefreshCookie_AndClearRefreshCookie(t *testing.T) {
	cfg := config.CookieConfig{
		Domain:   "example.com",
		Path:     "/",
		SameSite: "Lax",
		Secure:   false,
		MaxAge:   100,
	}

	w := httptest.NewRecorder()
	cookies.SetRefreshCookie(w, "refresh-token", cfg)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	gotCookies := res.Cookies()
	if len(gotCookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(gotCookies))
	}
	c := gotCookies[0]
	if c.Name != "refresh_token" || c.Value != "refresh-token" {
		t.Fatalf("unexpected refresh cookie: %+v", c)
	}
	if !c.HttpOnly || c.Secure {
		t.Fatalf("unexpected refresh cookie flags: %+v", c)
	}
	if c.SameSite != http.SameSiteLaxMode {
		t.Fatalf("expected SameSite=Lax, got %v", c.SameSite)
	}
	if c.MaxAge != 700 {
		t.Fatalf("expected max-age 700, got %d", c.MaxAge)
	}

	w = httptest.NewRecorder()
	cookies.ClearRefreshCookie(w, cfg)
	res = w.Result()
	defer func() { _ = res.Body.Close() }()
	c = res.Cookies()[0]
	if c.Name != "refresh_token" || c.MaxAge != -1 {
		t.Fatalf("unexpected cleared refresh cookie: %+v", c)
	}
}

func TestExtractTokens(t *testing.T) {
	t.Run("refresh success", func(t *testing.T) {
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "abc"})
		got, err := cookies.ExtractRefreshToken(req)
		if err != nil || got != "abc" {
			t.Fatalf("expected refresh token abc, got %q err=%v", got, err)
		}
	})

	t.Run("refresh missing", func(t *testing.T) {
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
		got, err := cookies.ExtractRefreshToken(req)
		if err == nil || got != "" {
			t.Fatalf("expected missing refresh token error")
		}
	})

	t.Run("auth success", func(t *testing.T) {
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: commonkeys.AuthTokenCookieName, Value: "jwt"})
		got, err := cookies.ExtractAuthToken(req)
		if err != nil || got != "jwt" {
			t.Fatalf("expected auth token jwt, got %q err=%v", got, err)
		}
	})

	t.Run("auth empty", func(t *testing.T) {
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: commonkeys.AuthTokenCookieName, Value: ""})
		got, err := cookies.ExtractAuthToken(req)
		if err != nil || got != "" {
			t.Fatalf("expected empty token with nil error, got token=%q err=%v", got, err)
		}
	})
}
