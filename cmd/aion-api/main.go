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

func main() {
	logger := setupLogger()
	cfg := loadConfig(logger)
	appDeps, cleanup := initDependencies(cfg, logger)
	defer cleanup()

	httpSrv := createHTTPServer(appDeps, &cfg, logger)
	graphqlSrv := createGraphQLServer(appDeps, cfg, logger)

	handleServers(httpSrv, graphqlSrv, cfg, logger)
}

func setupLogger() loggerPort.Logger {
	loggerInstance, cleanup := loggerBuilder.NewZapLogger()
	defer cleanup()
	return loggerAdapter.NewZapLoggerAdapter(loggerInstance)
}

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

func initDependencies(cfg config.Config, logger loggerPort.Logger) (*bootstrap.AppDependencies, func()) {
	appDeps, cleanup, err := bootstrap.InitializeDependencies(cfg, logger)
	if err != nil {
		response.HandleCriticalError(logger, constants.ErrInitializeDependencies, err)
		panic(err)
	}
	logger.Infow(constants.SuccessToInitializeDependencies)
	return appDeps, cleanup
}

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

func createGraphQLServer(appDeps *bootstrap.AppDependencies, cfg config.Config, logger loggerPort.Logger) *http.Server {
	graphqlSrv, err := graphqlserver.NewGraphqlServer(appDeps, cfg)
	if err != nil {
		logger.Errorw(constants.ErrStartGraphqlServer, constants.Error, err)
		panic(err)
	}
	logger.Infow(constants.GraphqlServerStarted, constants.ContextPath, cfg.ServerHTTP.Context)
	return graphqlSrv
}

func handleServers(httpSrv, graphqlSrv *http.Server, cfg config.Config, logger loggerPort.Logger) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	errChan := make(chan error, 2)
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("failed to start HTTP server: %w", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := graphqlSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("failed to start GraphQL server: %w", err)
		}
	}()

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
