// Package fxapp wires the application using Uber Fx modules.
package fxapp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/lechitz/aion-api/internal/platform/config"
	"github.com/lechitz/aion-api/internal/platform/ports/output/logger"
	serverHTTP "github.com/lechitz/aion-api/internal/platform/server/http"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"go.uber.org/fx"
)

// ServerModule wires HTTP handler + server lifecycle.
//
//nolint:gochecknoglobals // Fx modules are intended as package-level options.
var ServerModule = fx.Options(
	fx.Provide(
		ProvideHTTPHandler,
		ProvideHTTPServer,
	),
	fx.Invoke(RunHTTPServer),
)

// ProvideHTTPHandler composes REST/GraphQL handler from config + deps.
func ProvideHTTPHandler(cfg *config.Config, deps *AppDependencies, log logger.ContextLogger) (http.Handler, error) {
	return serverHTTP.ComposeHandler(cfg, deps, log)
}

// ProvideHTTPServer builds the HTTP server instance from handler/config.
func ProvideHTTPServer(cfg *config.Config, handler http.Handler, log logger.ContextLogger) *http.Server {
	params := serverHTTP.FromHTTP(cfg, handler)
	return serverHTTP.Build(context.Background(), params, log)
}

// RunHTTPServer starts the HTTP server and shuts it down gracefully via lifecycle hooks.
func RunHTTPServer(lc fx.Lifecycle, srv *http.Server, cfg *config.Config, log logger.ContextLogger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			_ = ctx
			go func() {
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Errorw(logMsgServerError, commonkeys.Error, err)
				}
			}()
			log.Infow(logMsgServerStarting, commonkeys.APIName, cfg.General.Name, commonkeys.AppEnv, cfg.General.Env)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			shutdownCtx, cancel := context.WithTimeout(ctx, cfg.Application.Timeout)
			defer cancel()
			if err := srv.Shutdown(shutdownCtx); err != nil {
				return fmt.Errorf(errMsgHTTPShutdown, err)
			}
			return nil
		},
	})
}
