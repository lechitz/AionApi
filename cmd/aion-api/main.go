// Package main provides the main entry point for the application.
package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/lechitz/AionApi/pkg/observability"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/lechitz/AionApi/internal/core/ports/output"

	"github.com/lechitz/AionApi/internal/def"

	"github.com/lechitz/AionApi/cmd/aion-api/constants"
	"github.com/lechitz/AionApi/internal/adapters/primary/graph/graphqlserver"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/httpserver"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/response"
	"github.com/lechitz/AionApi/internal/platform/bootstrap"
	"github.com/lechitz/AionApi/internal/platform/config"
	loggerBuilder "github.com/lechitz/AionApi/pkg/logger"
)

// main initializes and runs the AionAPI application lifecycle.
func main() {
	logger, cleanupLogger := loggerBuilder.NewZapLogger()
	defer cleanupLogger()

	cfg := loadConfig(logger)

	cleanupMetrics := observability.InitOtelMetrics(cfg, logger)
	defer cleanupMetrics()

	cleanupTracer := observability.InitTracer(cfg, logger)
	defer cleanupTracer()

	appDeps, cleanupDeps := initDependencies(cfg, logger)
	defer cleanupDeps()

	httpSrv := createHTTPServer(appDeps, &cfg, logger)
	graphqlSrv := createGraphQLServer(appDeps, cfg, logger)

	handleServers(httpSrv, graphqlSrv, cfg, logger)
}

// loadConfig loads the environment configuration using envconfig, panicking on failure.
func loadConfig(logger output.Logger) config.Config {
	cfgLoader := config.NewLoader()
	cfg, err := cfgLoader.Load(logger)
	if err != nil {
		response.HandleCriticalError(logger, constants.ErrToFailedLoadConfiguration, err)
		panic(err)
	}

	logger.Infof(constants.LoadedConfig, cfg)
	logger.Infow(constants.SuccessToLoadConfiguration)

	return cfg
}

// initDependencies initializes services, repositories, and infrastructure connections.
func initDependencies(cfg config.Config, logger output.Logger) (*bootstrap.AppDependencies, func()) {
	appDeps, cleanup, err := bootstrap.InitializeDependencies(cfg, logger)
	if err != nil {
		response.HandleCriticalError(logger, constants.ErrInitializeDependencies, err)
		panic(err)
	}

	logger.Infow(constants.SuccessToInitializeDependencies)

	return appDeps, cleanup
}

// createHTTPServer builds the HTTP server using configuration and application dependencies.
func createHTTPServer(appDeps *bootstrap.AppDependencies, cfg *config.Config, logger output.Logger) *http.Server {
	httpSrv, err := httpserver.NewHTTPServer(appDeps, cfg)
	if err != nil {
		response.HandleCriticalError(logger, constants.ErrStartHTTPServer, err)
		panic(err)
	}

	logger.Infow(constants.ServerHTTPStarted, def.Port, httpSrv.Addr, def.ContextPath, cfg.ServerHTTP.Context)

	return httpSrv
}

// createGraphQLServer builds the GraphQL server using configuration and application dependencies.
func createGraphQLServer(appDeps *bootstrap.AppDependencies, cfg config.Config, logger output.Logger) *http.Server {
	graphqlSrv, err := graphqlserver.NewGraphqlServer(appDeps, cfg)
	if err != nil {
		logger.Errorw(constants.ErrStartGraphqlServer, def.Error, err)
		panic(err)
	}

	logger.Infow(constants.GraphqlServerStarted, def.Port, cfg.ServerGraphql.Port, def.ContextPath, def.GraphQLPath)

	return graphqlSrv
}

// handleServers orchestrates concurrent HTTP and GraphQL server execution and graceful shutdown.
func handleServers(httpSrv, graphqlSrv *http.Server, cfg config.Config, logger output.Logger) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	errChan := make(chan error, 2)
	wg.Add(2)

	// Start HTTP server
	go func() {
		defer wg.Done()
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf(constants.ErrFailedToStartHTTPServer, err)
		}
	}()

	// Start a GraphQL server
	go func() {
		defer wg.Done()
		if err := graphqlSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf(constants.ErrFailedToStartGraphqlServer, err)
		}
	}()

	// Handle shutdown or error event
	select {
	case err := <-errChan:
		logger.Errorw(constants.MsgUnexpectedServerFailure, def.Error, err.Error())
		response.HandleCriticalError(logger, constants.MsgUnexpectedServerFailure, err)
		stop()
	case <-ctx.Done():
		logger.Infow(constants.MsgShutdownSignalReceived)
	}

	shutdownTimeout := time.Duration(cfg.Application.Timeout)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	_ = httpSrv.Shutdown(shutdownCtx)
	_ = graphqlSrv.Shutdown(shutdownCtx)

	wg.Wait()
}
