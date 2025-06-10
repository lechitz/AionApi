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
	"github.com/lechitz/AionApi/internal/infra/bootstrap"
	"github.com/lechitz/AionApi/internal/infra/config"
	loggerBuilder "github.com/lechitz/AionApi/pkg/logger"
)

func main() {
	loggerInstance, loggerCleanup := loggerBuilder.NewZapLogger()
	defer loggerCleanup()
	logger := loggerAdapter.NewZapLoggerAdapter(loggerInstance)

	logger.Infow(constants.StartingApplication)

	if err := config.Load(logger); err != nil {
		response.HandleCriticalError(logger, constants.ErrToFailedLoadConfiguration, err)
		return
	}
	logger.Infow(constants.SuccessToLoadConfiguration)

	appDeps, cleanup, err := bootstrap.InitializeDependencies(*config.Setting(), logger)
	if err != nil {
		response.HandleCriticalError(logger, constants.ErrInitializeDependencies, err)
		return
	}
	logger.Infow(constants.SuccessToInitializeDependencies)

	newHTTPServer, err := httpserver.NewHTTPServer(appDeps, config.Setting())
	if err != nil {
		response.HandleCriticalError(logger, constants.ErrStartHTTPServer, err)
		return
	}
	logger.Infow(
		constants.ServerHTTPStarted,
		constants.Port, newHTTPServer.Addr,
		constants.ContextPath, config.Setting().ServerHTTP.Context,
	)

	graphqlServer, err := graphqlserver.NewGraphqlServer(appDeps)
	if err != nil {
		logger.Errorw(constants.ErrStartGraphqlServer, constants.Error, err)
		return
	}
	logger.Infow(
		constants.GraphqlServerStarted,
		constants.ContextPath, config.Setting().ServerHTTP.Context,
	)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(2)
	errChan := make(chan error, 2)

	go func() {
		defer wg.Done()
		if err := newHTTPServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("failed to start HTTP server: %w", err)
		}
	}()
	go func() {
		defer wg.Done()
		if err := graphqlServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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

	shutdownTimeout := time.Duration(config.Setting().Application.Timeout)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	_ = newHTTPServer.Shutdown(shutdownCtx)
	_ = graphqlServer.Shutdown(shutdownCtx)

	cleanup()
	wg.Wait()
}
