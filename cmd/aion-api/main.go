package main

import (
	"context"
	"errors"
	"github.com/lechitz/AionApi/adapters/primary/http/middleware/response"
	"github.com/lechitz/AionApi/adapters/primary/http/server"
	loggerAdapter "github.com/lechitz/AionApi/adapters/secondary/logger"
	"github.com/lechitz/AionApi/cmd/aion-api/constants"
	"github.com/lechitz/AionApi/internal/infra/bootstrap"
	"github.com/lechitz/AionApi/internal/infra/config"
	loggerBuilder "github.com/lechitz/AionApi/pkg/logger"
	"net/http"
	"os/signal"
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

	newServer, err := server.NewHTTPServer(appDependencies, logger, &config.Setting)
	if err != nil {
		response.HandleCriticalError(logger, constants.ErrStartServer, err)
		return
	}

	logger.Infow(constants.ServerStarted,
		constants.Port, newServer.Addr,
		constants.ContextPath, config.Setting.Server.Context,
	)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := newServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			response.HandleCriticalError(logger, constants.ErrStartServer, err)
		}
	}()

	<-ctx.Done()
	logger.Infow(constants.MsgShutdownSignalReceived)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := newServer.Shutdown(shutdownCtx); err != nil {
		logger.Errorw(constants.ErrGracefulShutdown, constants.Error, err.Error())
	} else {
		logger.Infow(constants.MsgGracefulShutdownSuccess)
	}

	cleanup()
}
