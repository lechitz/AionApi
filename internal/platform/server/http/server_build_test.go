package http_test

import (
	"context"
	"net"
	stdhttp "net/http"
	"testing"
	"time"

	"github.com/lechitz/aion-api/internal/platform/config"
	serverhttp "github.com/lechitz/aion-api/internal/platform/server/http"
	"github.com/stretchr/testify/require"
)

func TestFromHTTP(t *testing.T) {
	cfg := &config.Config{
		ServerHTTP: config.ServerHTTP{
			Host:              "127.0.0.1",
			Port:              "8080",
			ReadTimeout:       2 * time.Second,
			WriteTimeout:      3 * time.Second,
			ReadHeaderTimeout: 4 * time.Second,
			IdleTimeout:       5 * time.Second,
			MaxHeaderBytes:    2048,
		},
	}

	handler := stdhttp.NewServeMux()
	params := serverhttp.FromHTTP(cfg, handler)

	require.Equal(t, net.JoinHostPort("127.0.0.1", "8080"), params.Addr)
	require.Equal(t, handler, params.Handler)
	require.Equal(t, 2*time.Second, params.ReadTimeout)
	require.Equal(t, 3*time.Second, params.WriteTimeout)
	require.Equal(t, 4*time.Second, params.ReadHeaderTimeout)
	require.Equal(t, 5*time.Second, params.IdleTimeout)
	require.Equal(t, 2048, params.MaxHeaderBytes)
}

type testCtxKey string

func TestBuild(t *testing.T) {
	baseCtx := context.WithValue(t.Context(), testCtxKey("k"), "v")
	handler := stdhttp.NewServeMux()
	params := serverhttp.Params{
		Addr:              "127.0.0.1:8080",
		Handler:           handler,
		ReadTimeout:       time.Second,
		WriteTimeout:      2 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		IdleTimeout:       4 * time.Second,
		MaxHeaderBytes:    1024,
	}

	srv := serverhttp.Build(baseCtx, params, nil)
	require.Equal(t, params.Addr, srv.Addr)
	require.Equal(t, handler, srv.Handler)
	require.Equal(t, params.ReadTimeout, srv.ReadTimeout)
	require.Equal(t, params.WriteTimeout, srv.WriteTimeout)
	require.Equal(t, params.ReadHeaderTimeout, srv.ReadHeaderTimeout)
	require.Equal(t, params.IdleTimeout, srv.IdleTimeout)
	require.Equal(t, params.MaxHeaderBytes, srv.MaxHeaderBytes)

	ctx := srv.BaseContext(nil)
	require.Equal(t, "v", ctx.Value(testCtxKey("k")))
}
