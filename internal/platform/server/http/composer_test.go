package http_test

import (
	"context"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lechitz/aion-api/internal/platform/app"
	"github.com/lechitz/aion-api/internal/platform/config"
	serverhttp "github.com/lechitz/aion-api/internal/platform/server/http"
	"github.com/stretchr/testify/require"
)

type mockServerLogger struct{}

func (mockServerLogger) Infof(string, ...any)                      {}
func (mockServerLogger) Errorf(string, ...any)                     {}
func (mockServerLogger) Debugf(string, ...any)                     {}
func (mockServerLogger) Warnf(string, ...any)                      {}
func (mockServerLogger) Infow(string, ...any)                      {}
func (mockServerLogger) Errorw(string, ...any)                     {}
func (mockServerLogger) Debugw(string, ...any)                     {}
func (mockServerLogger) Warnw(string, ...any)                      {}
func (mockServerLogger) InfowCtx(context.Context, string, ...any)  {}
func (mockServerLogger) ErrorwCtx(context.Context, string, ...any) {}
func (mockServerLogger) WarnwCtx(context.Context, string, ...any)  {}
func (mockServerLogger) DebugwCtx(context.Context, string, ...any) {}

func testHTTPConfig() *config.Config {
	return &config.Config{
		General: config.GeneralConfig{
			Name:    "aion-api",
			Env:     "test",
			Version: "v1",
		},
		Observability: config.ObservabilityConfig{
			OtelServiceName: "aion-api",
		},
		ServerHTTP: config.ServerHTTP{
			Context:          "/aion",
			APIRoot:          "/api/v1",
			SwaggerMountPath: "/swagger",
			DocsAliasPath:    "/docs",
			HealthRoute:      "/health",
		},
		ServerGraphql: config.ServerGraphql{
			Path: "/graphql",
		},
		Application: config.Application{
			Timeout: time.Second,
		},
	}
}

func TestComposeHandler_HealthAndDocsRoutes(t *testing.T) {
	cfg := testHTTPConfig()
	handler, err := serverhttp.ComposeHandler(cfg, &app.Dependencies{}, mockServerLogger{})
	require.NoError(t, err)
	require.NotNil(t, handler)

	req := httptest.NewRequestWithContext(t.Context(), stdhttp.MethodGet, "/aion/health", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, stdhttp.StatusOK, rec.Code)

	req = httptest.NewRequestWithContext(t.Context(), stdhttp.MethodGet, "/aion/docs", nil)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, stdhttp.StatusTemporaryRedirect, rec.Code)
}

func TestComposeHandler_DefaultFallbackRoutes(t *testing.T) {
	cfg := testHTTPConfig()
	cfg.ServerHTTP.Context = ""
	cfg.ServerHTTP.SwaggerMountPath = ""
	cfg.ServerHTTP.DocsAliasPath = ""
	cfg.ServerHTTP.HealthRoute = ""

	handler, err := serverhttp.ComposeHandler(cfg, &app.Dependencies{}, mockServerLogger{})
	require.NoError(t, err)
	require.NotNil(t, handler)

	t.Run("root health", func(t *testing.T) {
		req := httptest.NewRequestWithContext(t.Context(), stdhttp.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		require.Equal(t, stdhttp.StatusOK, rec.Code)
	})

	t.Run("root health with trailing slash", func(t *testing.T) {
		req := httptest.NewRequestWithContext(t.Context(), stdhttp.MethodGet, "/health/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		require.Equal(t, stdhttp.StatusOK, rec.Code)
	})

	t.Run("alt health under api root", func(t *testing.T) {
		req := httptest.NewRequestWithContext(t.Context(), stdhttp.MethodGet, "/api/v1/health", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		require.Equal(t, stdhttp.StatusOK, rec.Code)
	})

	t.Run("docs alias fallback redirect", func(t *testing.T) {
		req := httptest.NewRequestWithContext(t.Context(), stdhttp.MethodGet, "/docs", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		require.Equal(t, stdhttp.StatusTemporaryRedirect, rec.Code)
		require.Equal(t, "/swagger/index.html", rec.Header().Get("Location"))
	})
}
