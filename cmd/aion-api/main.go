package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/lechitz/AionApi/internal/adapters/secondary/security/keygen"
	"github.com/lechitz/AionApi/internal/platform/server"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"

	"github.com/lechitz/AionApi/internal/platform/contextlogger"

	"github.com/lechitz/AionApi/cmd/aion-api/constants"
	"github.com/lechitz/AionApi/internal/adapters/primary/graph/graphqlserver"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/httpserver"
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/platform/bootstrap"
	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/observability"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	logger, cleanupLogger := contextlogger.New()
	defer cleanupLogger()

	keyGenerator := keygen.New()

	cfg := loadConfig(keyGenerator, logger)

	appCtx, stopApp := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopApp()

	cleanupMetrics := observability.InitOtelMetrics(cfg, logger)
	defer cleanupMetrics()

	cleanupTracer := observability.InitTracer(cfg, logger)
	defer cleanupTracer()

	appDeps, cleanupDeps := initDependencies(appCtx, cfg, logger)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	defer cleanupDeps(shutdownCtx)

	servers := buildAllServers(appCtx, cfg, appDeps, logger)
	handleServers(appCtx, cfg, servers, logger, stopApp)
}

func loadConfig(keyGenerator output.KeyGenerator, logger output.ContextLogger) *config.Config {
	configLoader := config.New(keyGenerator)

	cfg, err := configLoader.Load(logger)
	if err != nil {
		logger.Errorw(constants.ErrToFailedLoadConfiguration, commonkeys.Error, err.Error())
		os.Exit(1)
	}
	if err := cfg.Validate(); err != nil {
		logger.Errorw(constants.ErrInvalidConfiguration, commonkeys.Error, err.Error())
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

func initDependencies(appCtx context.Context, cfg *config.Config, logger output.ContextLogger) (*bootstrap.AppDependencies, func(context.Context)) {
	deps, cleanup, err := bootstrap.InitializeDependencies(appCtx, cfg, logger)
	if err != nil {
		logger.Errorw(constants.ErrInitializeDependencies, commonkeys.Error, err.Error())
		os.Exit(1)
	}
	logger.Infow(constants.SuccessToInitializeDependencies)
	return deps, cleanup
}

func setupHTTPHandler(cfg *config.Config, deps *bootstrap.AppDependencies, logger output.ContextLogger) (http.Handler, error) {
	router, err := httpserver.ComposeRouter(cfg, deps.TokenService, deps.UserService, deps.AuthService, logger)
	if err != nil {
		logger.Errorw(constants.ErrStartHTTPServer, commonkeys.Error, err.Error())
		return nil, err
	}
	return otelhttp.NewHandler(router, fmt.Sprintf("%s-REST", cfg.Observability.OtelServiceName)), nil
}

func setupGraphQLHandler(cfg *config.Config, deps *bootstrap.AppDependencies, logger output.ContextLogger) (http.Handler, error) {
	handler, err := graphqlserver.NewGraphqlHandler(deps.TokenService, deps.CategoryService, logger, cfg)
	if err != nil {
		logger.Errorw(constants.ErrStartGraphqlServer, commonkeys.Error, err.Error())
		return nil, err
	}
	return otelhttp.NewHandler(handler, fmt.Sprintf("%s-GraphQL", cfg.Observability.OtelServiceName)), nil
}

func buildAllServers(appCtx context.Context, cfg *config.Config, appDeps *bootstrap.AppDependencies, logger output.ContextLogger) []*http.Server {
	httpHandler, err := setupHTTPHandler(cfg, appDeps, logger)
	if err != nil {
		logger.Errorw(constants.ErrStartHTTPServer, commonkeys.Error, err.Error())
		os.Exit(1)
	}
	httpParams := server.FromHTTP(cfg, httpHandler)

	gqlHandler, err := setupGraphQLHandler(cfg, appDeps, logger)
	if err != nil {
		logger.Errorw(constants.ErrStartGraphqlServer, commonkeys.Error, err.Error())
		os.Exit(1)
	}
	gqlParams := server.FromGraphQL(cfg, gqlHandler)

	return []*http.Server{
		server.Build(appCtx, httpParams, logger),
		server.Build(appCtx, gqlParams, logger),
	}
}

func handleServers(ctx context.Context, cfg *config.Config, servers []*http.Server, logger output.ContextLogger, stop context.CancelFunc) {
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
		logger.Errorw(constants.MsgUnexpectedServerFailure, commonkeys.Error, err.Error())
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
			logger.Errorw(constants.ShutdownFailureFmt, commonkeys.ServerAddr, srv.Addr, commonkeys.Error, err)
		}
	}
}
