package cors_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lechitz/aion-api/internal/platform/server/http/middleware/cors"
)

func TestCORSMiddleware_AllowsConfiguredOriginAndCredentials(t *testing.T) {
	mw := cors.New()
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://localhost:5000")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if got := res.Header.Get("Access-Control-Allow-Origin"); got != "http://localhost:5000" {
		t.Fatalf("unexpected allow-origin header: %q", got)
	}
	if got := res.Header.Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Fatalf("unexpected allow-credentials header: %q", got)
	}
}

func TestCORSMiddleware_HandlesPreflight(t *testing.T) {
	mw := cors.New()
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodOptions, "/", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type, Authorization")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected preflight status 200, got %d", res.StatusCode)
	}
	if got := res.Header.Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Fatalf("unexpected allow-origin header on preflight: %q", got)
	}
}
