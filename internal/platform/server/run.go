// Package server is the server implementation.
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

// RunAll decide (via cfg) quais servidores sobem, registra start/interrupt
// no runtime e bloqueia até sinal/erro, aplicando graceful shutdown.
func RunAll(ctx context.Context, cfg *config.Config, deps *bootstrap.AppDependencies, log logger.ContextLogger) error {
	var group runtime.Group

	// --------------------- HTTP (REST + GraphQL montado no composer) ---------------------
	httpHandler, err := serverHTTP.ComposeHandler(cfg, deps, log)
	if err != nil {
		return fmt.Errorf("compose HTTP handler: %w", err)
	}
	httpParams := serverHTTP.FromHTTP(cfg, httpHandler)
	httpSrv := serverHTTP.Build(ctx, httpParams, log)

	group.Add(
		func() error {
			if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				return err
			}
			return nil
		},
		func(_ error) {
			shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Application.Timeout)
			defer cancel()
			_ = httpSrv.Shutdown(shutdownCtx)
		},
	)

	// --------------------- gRPC (opcional, via config) -----------------------------------
	// if cfg.ServerGRPC.Enabled {
	//	grpcSrv, err := serverGRPC.ComposeServer(cfg, deps, log)
	//	if err != nil {
	//		return fmt.Errorf("compose gRPC server: %w", err)
	//	}
	//
	//	grpcParams := serverGRPC.FromGRPC(cfg) // ex.: Addr
	//	lis, err := net.Listen("tcp", grpcParams.Addr)
	//	if err != nil {
	//		return fmt.Errorf("gRPC listen %s: %w", grpcParams.Addr, err)
	//	}
	//
	//	group.Add(
	//		func() error { return grpcSrv.Serve(lis) },
	//		func(_ error) {
	//			grpcSrv.GracefulStop()
	//			_ = lis.Close()
	//		},
	//	)
	// }

	// --------------------- Supervisor único ---------------------------------------------
	log.Infow("servers starting...",
		commonkeys.APIName, cfg.General.Name,
		commonkeys.AppEnv, cfg.General.Env,
	)
	group.Run(ctx, cfg.Application.Timeout, log)
	return nil
}
