package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/lechitz/AionApi/internal/platform/contextlogger"

	"github.com/lechitz/AionApi/cmd/aion-api/constants"
	"github.com/lechitz/AionApi/internal/adapters/primary/graph/graphqlserver"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/httpserver"
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/platform/bootstrap"
	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/observability"
	"github.com/lechitz/AionApi/internal/shared/commonkeys"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	logger, cleanupLogger := contextlogger.New()
	defer cleanupLogger()

	cfg := loadConfig(logger)

	appCtx, stopApp := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopApp()

	cleanupMetrics := observability.InitOtelMetrics(cfg, logger)
	defer cleanupMetrics()

	cleanupTracer := observability.InitTracer(cfg, logger)
	defer cleanupTracer()

	appDeps, cleanupDeps := initDependencies(appCtx, cfg, logger)
	defer cleanupDeps()

	servers := buildAllServers(appCtx, appDeps, cfg, logger)
	handleServers(appCtx, servers, cfg, logger, stopApp)
}

func loadConfig(logger output.ContextLogger) *config.Config {
	configLoader := config.New()
	cfg, err := configLoader.Load(logger)
	if err != nil {
		logger.Errorw(
			constants.ErrToFailedLoadConfiguration,
			commonkeys.Error, err.Error(),
		)

		os.Exit(1)
	}

	if err := cfg.Validate(); err != nil {
		logger.Errorw(
			constants.ErrInvalidConfiguration,
			commonkeys.Error, err.Error(),
		)

		os.Exit(1)
	}

	logger.Infow(
		constants.SuccessToLoadConfiguration,
		commonkeys.APIName, cfg.General.Name,
		commonkeys.AppEnv, cfg.General.Env,
		commonkeys.AppVersion, cfg.General.Version,
	)

	return cfg
}

func initDependencies(ctx context.Context, cfg *config.Config, logger output.ContextLogger) (*bootstrap.AppDependencies, func()) {
	deps, cleanup, err := bootstrap.InitializeDependencies(ctx, cfg, logger)
	if err != nil {
		logger.Errorw(
			constants.ErrInitializeDependencies,
			commonkeys.Error, err.Error(),
		)

		os.Exit(1)
	}

	logger.Infow(constants.SuccessToInitializeDependencies)
	return deps, cleanup
}

func setupHTTPHandler(deps *bootstrap.AppDependencies, cfg *config.Config, logger output.ContextLogger) http.Handler {
	router, err := httpserver.ComposeRouter(deps, cfg)
	if err != nil {
		logger.Errorw(
			constants.ErrStartHTTPServer,
			commonkeys.Error, err.Error(),
		)

		os.Exit(1)
	}

	return otelhttp.NewHandler(router, fmt.Sprintf("%s-REST", cfg.Observability.OtelServiceName))
}

func setupGraphQLHandler(deps *bootstrap.AppDependencies, cfg *config.Config, logger output.ContextLogger) http.Handler {
	handler, err := graphqlserver.NewGraphqlHandler(deps, cfg)
	if err != nil {
		logger.Errorw(
			constants.ErrStartGraphqlServer,
			commonkeys.Error, err.Error())

		os.Exit(1)
	}

	return otelhttp.NewHandler(handler, fmt.Sprintf("%s-GraphQL", cfg.Observability.OtelServiceName))
}

func buildServer(ctx context.Context, name, host, port string, handler http.Handler, readTimeout, writeTimeout time.Duration, logger output.ContextLogger) *http.Server {
	addr := net.JoinHostPort(host, port)

	srv := &http.Server{
		Addr:         addr,
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		Handler:      handler,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	logger.Infow(fmt.Sprintf(
		constants.ServerStartFmt, name),
		commonkeys.ServerHTTPName, name,
		commonkeys.ServerHTTPAddr, addr,
	)

	return srv
}

func buildAllServers(ctx context.Context, deps *bootstrap.AppDependencies, cfg *config.Config, logger output.ContextLogger) []*http.Server {
	return []*http.Server{
		buildServer(
			ctx,
			cfg.ServerHTTP.Name,
			cfg.ServerHTTP.Host,
			cfg.ServerHTTP.Port,
			setupHTTPHandler(deps, cfg, logger),
			cfg.ServerHTTP.ReadTimeout,
			cfg.ServerHTTP.WriteTimeout,
			logger,
		),
		buildServer(
			ctx,
			cfg.ServerGraphql.Name,
			cfg.ServerGraphql.Host,
			cfg.ServerGraphql.Port,
			setupGraphQLHandler(deps, cfg, logger),
			cfg.ServerGraphql.ReadTimeout,
			cfg.ServerGraphql.WriteTimeout,
			logger,
		),
	}
}

func handleServers(ctx context.Context, servers []*http.Server, cfg *config.Config, logger output.ContextLogger, stop context.CancelFunc) {
	var wg sync.WaitGroup
	errChan := make(chan error, len(servers))

	for _, srv := range servers {
		wg.Add(1)
		go func(s *http.Server) {
			defer wg.Done()
			if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errChan <- fmt.Errorf(constants.ServerFailureFmt, s.Addr, err)
			}
		}(srv)
	}

	select {
	case err := <-errChan:
		logger.Errorw(
			constants.MsgUnexpectedServerFailure,
			commonkeys.Error, err.Error(),
		)
		stop()

	case <-ctx.Done():
		logger.Infow(constants.MsgShutdownSignalReceived)
	}

	shutdownServers(servers, cfg.Application.Timeout, logger)
	wg.Wait()
}

func shutdownServers(servers []*http.Server, timeout time.Duration, logger output.ContextLogger) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for _, srv := range servers {
		if err := srv.Shutdown(ctx); err != nil {
			logger.Errorw(
				constants.ShutdownFailureFmt, srv.Addr, err,
				commonkeys.ServerHTTPAddr, srv.Addr,
				commonkeys.Error, err.Error(),
			)
		}
	}
}
