package main

import (
	"context"
	"errors"
	"github.com/lechitz/AionApi/adapters/primary/graph/graphqlserver"
	"github.com/lechitz/AionApi/adapters/primary/http/httpserver"
	"github.com/lechitz/AionApi/adapters/primary/http/middleware/response"
	loggerAdapter "github.com/lechitz/AionApi/adapters/secondary/logger"
	"github.com/lechitz/AionApi/cmd/aion-api/constants"
	"github.com/lechitz/AionApi/internal/infra/bootstrap"
	"github.com/lechitz/AionApi/internal/infra/config"
	loggerBuilder "github.com/lechitz/AionApi/pkg/logger"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"
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

	appDependencies, cleanup, err := bootstrap.InitializeDependencies(config.Setting, logger)
	if err != nil {
		response.HandleCriticalError(logger, constants.ErrInitializeDependencies, err)
		return
	}
	logger.Infow(constants.SuccessToInitializeDependencies)

	newHTTPServer, err := httpserver.NewHTTPServer(appDependencies, &config.Setting)
	if err != nil {
		response.HandleCriticalError(logger, constants.ErrStartHTTPServer, err)
		return
	}
	logger.Infow(constants.ServerHTTPStarted, constants.Port, newHTTPServer.Addr, constants.ContextPath, config.Setting.ServerHTTP.Context)

	newGraphqlServerHandler, err := graphqlserver.NewGraphqlServer(appDependencies)
	if err != nil {
		response.HandleCriticalError(logger, constants.ErrStartGraphqlServer, err)
		return
	}
	logger.Infow(constants.GraphqlServerStarted, constants.ContextPath, config.Setting.ServerHTTP.Context)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := newHTTPServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			response.HandleCriticalError(logger, constants.ErrStartHTTPServer, err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := http.ListenAndServe(":"+config.Setting.ServerGraphql.Port, newGraphqlServerHandler); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorw(constants.ErrStartGraphqlServer, constants.Error, err.Error())
			response.HandleCriticalError(logger, constants.ErrStartGraphqlServer, err)
		}
	}()

	<-ctx.Done()
	logger.Infow(constants.MsgShutdownSignalReceived)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := newHTTPServer.Shutdown(shutdownCtx); err != nil {
		logger.Errorw(constants.ErrGracefulShutdown, constants.Error, err.Error())
	} else {
		logger.Infow(constants.MsgGracefulShutdownSuccess)
	}

	cleanup()
	wg.Wait()
}
