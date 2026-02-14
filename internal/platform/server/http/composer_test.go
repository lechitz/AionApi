package http_test

import (
	"context"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lechitz/AionApi/internal/platform/app"
	"github.com/lechitz/AionApi/internal/platform/config"
	serverhttp "github.com/lechitz/AionApi/internal/platform/server/http"
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
			Name:    "AionApi",
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

	req := httptest.NewRequest(stdhttp.MethodGet, "/aion/health", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, stdhttp.StatusOK, rec.Code)

	req = httptest.NewRequest(stdhttp.MethodGet, "/aion/docs", nil)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, stdhttp.StatusTemporaryRedirect, rec.Code)
}
