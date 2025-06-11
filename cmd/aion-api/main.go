// Package main provides the main entry point for the application.
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/lechitz/AionApi/adapters/primary/graph/graphqlserver"
	"github.com/lechitz/AionApi/adapters/primary/http/httpserver"
	"github.com/lechitz/AionApi/adapters/primary/http/middleware/response"
	loggerAdapter "github.com/lechitz/AionApi/adapters/secondary/logger"
	"github.com/lechitz/AionApi/cmd/aion-api/constants"
	loggerPort "github.com/lechitz/AionApi/internal/core/ports/output/logger"
	"github.com/lechitz/AionApi/internal/infra/bootstrap"
	"github.com/lechitz/AionApi/internal/infra/config"
	loggerBuilder "github.com/lechitz/AionApi/pkg/logger"
)

// main initializes and runs the AionAPI application lifecycle.
func main() {
	logger, cleanup := loggerBuilder.NewZapLogger()
	defer cleanup()

	appLogger := loggerAdapter.NewZapLoggerAdapter(logger)

	cfg := loadConfig(appLogger)
	appLogger.Infof("loaded config: %+v", cfg)

	appDeps, cleanupDeps := initDependencies(cfg, appLogger)
	defer cleanupDeps()

	httpSrv := createHTTPServer(appDeps, &cfg, appLogger)
	graphqlSrv := createGraphQLServer(appDeps, cfg, appLogger)

	handleServers(httpSrv, graphqlSrv, cfg, appLogger)
}

// loadConfig loads the environment configuration using envconfig, panicking on failure.
func loadConfig(logger loggerPort.Logger) config.Config {
	cfgLoader := config.NewLoader()
	cfg, err := cfgLoader.Load(logger)
	if err != nil {
		response.HandleCriticalError(logger, constants.ErrToFailedLoadConfiguration, err)
		panic(err)
	}
	logger.Infow(constants.SuccessToLoadConfiguration)
	return cfg
}

// initDependencies initializes services, repositories, and infrastructure connections.
func initDependencies(cfg config.Config, logger loggerPort.Logger) (*bootstrap.AppDependencies, func()) {
	appDeps, cleanup, err := bootstrap.InitializeDependencies(cfg, logger)
	if err != nil {
		response.HandleCriticalError(logger, constants.ErrInitializeDependencies, err)
		panic(err)
	}
	logger.Infow(constants.SuccessToInitializeDependencies)
	return appDeps, cleanup
}

// createHTTPServer builds the HTTP server using configuration and application dependencies.
func createHTTPServer(appDeps *bootstrap.AppDependencies, cfg *config.Config, logger loggerPort.Logger) *http.Server {
	httpSrv, err := httpserver.NewHTTPServer(appDeps, cfg)
	if err != nil {
		response.HandleCriticalError(logger, constants.ErrStartHTTPServer, err)
		panic(err)
	}
	logger.Infow(
		constants.ServerHTTPStarted,
		constants.Port, httpSrv.Addr,
		constants.ContextPath, cfg.ServerHTTP.Context,
	)
	return httpSrv
}

// createGraphQLServer builds the GraphQL server using configuration and application dependencies.
func createGraphQLServer(appDeps *bootstrap.AppDependencies, cfg config.Config, logger loggerPort.Logger) *http.Server {
	graphqlSrv, err := graphqlserver.NewGraphqlServer(appDeps, cfg)
	if err != nil {
		logger.Errorw(constants.ErrStartGraphqlServer, constants.Error, err)
		panic(err)
	}
	logger.Infow(
		constants.GraphqlServerStarted,
		constants.Port, cfg.ServerGraphql.Port,
		constants.ContextPath, "/graphql",
	)

	return graphqlSrv
}

// handleServers orchestrates concurrent HTTP and GraphQL server execution and graceful shutdown.
func handleServers(httpSrv, graphqlSrv *http.Server, cfg config.Config, logger loggerPort.Logger) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	errChan := make(chan error, 2)
	wg.Add(2)

	// Start HTTP server
	go func() {
		defer wg.Done()
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("failed to start HTTP server: %w", err)
		}
	}()

	// Start a GraphQL server
	go func() {
		defer wg.Done()
		if err := graphqlSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("failed to start GraphQL server: %w", err)
		}
	}()

	// Handle shutdown or error event
	select {
	case err := <-errChan:
		logger.Errorw("server error", "error", err.Error())
		response.HandleCriticalError(logger, constants.ErrStartHTTPServer, err)
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
