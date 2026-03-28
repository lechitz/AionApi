package chi_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lechitz/aion-api/internal/platform/server/http/ports"
	chirouter "github.com/lechitz/aion-api/internal/platform/server/http/router/chi"
)

func TestChiRouter_BasicVerbsAndServeHTTP(t *testing.T) {
	r := chirouter.New()

	r.GET("/g", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	r.POST("/p", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	r.PUT("/u", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))
	r.DELETE("/d", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	cases := []struct {
		method string
		path   string
		want   int
	}{
		{http.MethodGet, "/g", http.StatusOK},
		{http.MethodPost, "/p", http.StatusCreated},
		{http.MethodPut, "/u", http.StatusAccepted},
		{http.MethodDelete, "/d", http.StatusNoContent},
	}
	for _, tc := range cases {
		req := httptest.NewRequestWithContext(t.Context(), tc.method, tc.path, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != tc.want {
			t.Fatalf("%s %s expected %d got %d", tc.method, tc.path, tc.want, w.Code)
		}
	}
}

func TestChiRouter_UseAppliesMiddleware(t *testing.T) {
	r := chirouter.New()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("X-MW", "applied")
			next.ServeHTTP(w, req)
		})
	}, nil)

	r.GET("/x", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if got := w.Result().Header.Get("X-MW"); got != "applied" {
		t.Fatalf("expected middleware header, got %q", got)
	}
}

func TestChiRouter_GroupAndGroupWith(t *testing.T) {
	r := chirouter.New()

	r.Group("/api", func(gr ports.Router) {
		gr.GET("/health", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
	})

	r.GroupWith(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("X-Scoped", "1")
			next.ServeHTTP(w, req)
		})
	}, func(gr ports.Router) {
		gr.GET("/scoped", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected grouped route status 200, got %d", w.Code)
	}

	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/scoped", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected groupWith route status 200, got %d", w.Code)
	}
	if got := w.Result().Header.Get("X-Scoped"); got != "1" {
		t.Fatalf("expected scoped middleware header, got %q", got)
	}
}

func TestChiRouter_MountAndCustom404405(t *testing.T) {
	r := chirouter.New()

	r.Mount("/m", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r.SetNotFound(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))
	r.SetMethodNotAllowed(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusConflict)
	}))
	r.SetError(func(http.ResponseWriter, *http.Request, error) {})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/m/anything", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected mounted route status 200, got %d", w.Code)
	}

	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/not-found", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusTeapot {
		t.Fatalf("expected custom 404 status 418, got %d", w.Code)
	}

	r.POST("/only-post", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/only-post", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusConflict {
		t.Fatalf("expected custom 405 status 409, got %d", w.Code)
	}
}
