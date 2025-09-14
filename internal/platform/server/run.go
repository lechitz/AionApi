// Package server wires up and runs the application servers.
package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/lechitz/AionApi/internal/platform/bootstrap"
	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	serverHTTP "github.com/lechitz/AionApi/internal/platform/server/http"
	"github.com/lechitz/AionApi/internal/platform/server/runtime"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
)

// Log / error message constants.
const (
	// LogServersStarting is emitted right before servers start.
	LogServersStarting = "servers starting..."

	// ErrComposeHTTPHandler is used when composing the HTTP handler fails.
	ErrComposeHTTPHandler = "compose HTTP handler"
)

// RunAll composes and starts the HTTP server, registers start/interrupt hooks
// in the runtime group, and blocks until an error or shutdown signal occurs.
// A graceful shutdown is applied using the configured timeout.
func RunAll(ctx context.Context, cfg *config.Config, deps *bootstrap.AppDependencies, log logger.ContextLogger) error {
	var group runtime.Group

	// --- HTTP (REST + GraphQL) ---
	httpHandler, err := serverHTTP.ComposeHandler(cfg, deps, log)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrComposeHTTPHandler, err)
	}
	httpParams := serverHTTP.FromHTTP(cfg, httpHandler)
	httpSrv := serverHTTP.Build(ctx, httpParams, log)

	group.Add(
		func() error {
			// Start HTTP server and propagate non-graceful errors.
			if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				return err
			}
			return nil
		},
		func(_ error) {
			// Gracefully shutdown HTTP server.
			shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Application.Timeout)
			defer cancel()
			_ = httpSrv.Shutdown(shutdownCtx)
		},
	)

	// --- gRPC (optional) ---
	// gRPC is not enabled in this project. If needed later:
	// - Create a composer under internal/platform/server/grpc (mirroring HTTP).
	// - Build a net listener and the gRPC server.
	// - Register start/stop with group.Add(...) similar to the HTTP block above.

	log.Infow(
		LogServersStarting,
		commonkeys.APIName, cfg.General.Name,
		commonkeys.AppEnv, cfg.General.Env,
	)

	group.Run(ctx, cfg.Application.Timeout, log)
	return nil
}
