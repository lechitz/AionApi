package main

import (
	"github.com/lechitz/AionApi/adapters/primary/http/middleware/response"
	"github.com/lechitz/AionApi/adapters/primary/http/server"
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

	appDependencies, cleanup, err := bootstrap.InitializeDependencies(config.Setting, logger)
	if err != nil {
		response.HandleCriticalError(logger, constants.ErrInitializeDependencies, err)
		return
	}
	defer cleanup()

	logger.Infow(constants.SuccessToInitializeDependencies)

	newServer, err := server.NewHTTPServer(appDependencies, logger, &config.Setting)
	if err != nil {
		response.HandleCriticalError(logger, constants.ErrStartServer, err)
		return
	}

	logger.Infow(constants.ServerStarted, constants.Port, newServer.Addr, constants.ContextPath, config.Setting.Server.Context)

	if err := newServer.ListenAndServe(); err != nil {
		logger.Errorw(constants.ErrStartServer, constants.Error, err)
		response.HandleCriticalError(logger, constants.ErrStartServer, err)
	}
}
