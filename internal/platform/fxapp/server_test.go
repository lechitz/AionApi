//nolint:testpackage // tests require package-private lifecycle wiring helpers.
package fxapp

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

type fakeLifecycle struct {
	hooks []fx.Hook
}

func (f *fakeLifecycle) Append(h fx.Hook) {
	f.hooks = append(f.hooks, h)
}

func testFxConfig() *config.Config {
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
			Host:              "127.0.0.1",
			Port:              "5001",
			Context:           "/aion",
			APIRoot:           "/api/v1",
			SwaggerMountPath:  "/swagger",
			DocsAliasPath:     "/docs",
			HealthRoute:       "/health",
			ReadTimeout:       2 * time.Second,
			WriteTimeout:      2 * time.Second,
			ReadHeaderTimeout: time.Second,
			IdleTimeout:       time.Second,
			MaxHeaderBytes:    1024,
		},
		ServerGraphql: config.ServerGraphql{
			Path: "/graphql",
		},
		Application: config.Application{
			Timeout: time.Second,
		},
	}
}

func TestProvideHTTPHandlerAndServer(t *testing.T) {
	cfg := testFxConfig()
	log := noopLoggerFx{}
	deps := &AppDependencies{}

	handler, err := ProvideHTTPHandler(cfg, deps, log)
	require.NoError(t, err)
	require.NotNil(t, handler)

	req := httptest.NewRequest(http.MethodGet, "/aion/health", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	srv := ProvideHTTPServer(cfg, handler, log)
	require.NotNil(t, srv)
	require.Equal(t, "127.0.0.1:5001", srv.Addr)
}

func TestRunHTTPServerRegistersLifecycleHooks(t *testing.T) {
	cfg := testFxConfig()
	log := noopLoggerFx{}
	srv := &http.Server{
		Addr:              "127.0.0.1:0",
		Handler:           http.NewServeMux(),
		ReadHeaderTimeout: time.Second,
	}

	lc := &fakeLifecycle{}
	RunHTTPServer(lc, srv, cfg, log)
	require.Len(t, lc.hooks, 1)

	require.NoError(t, lc.hooks[0].OnStart(t.Context()))
	require.NoError(t, lc.hooks[0].OnStop(t.Context()))
}
