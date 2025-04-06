package main

import (
	"github.com/lechitz/AionApi/cmd/aion-api/constants"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/server"
	"github.com/lechitz/AionApi/internal/platform/bootstrap"
	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/logger"
	"github.com/lechitz/AionApi/pkg/utils"
)

func main() {

	loggerSugar, closeLogger := logger.InitLoggerSugar()
	defer closeLogger()

	loggerSugar.Infow(constants.StartingApplication)

	if err := config.LoadConfig(loggerSugar); err != nil {
		utils.HandleCriticalError(loggerSugar, constants.ErrToFailedLoadConfiguration, err)
		return
	}

	loggerSugar.Infow(constants.SuccessToLoadConfiguration, constants.Settings, config.Setting)

	appDependencies, cleanup, err := bootstrap.InitializeDependencies(loggerSugar, config.Setting)
	if err != nil {
		utils.HandleCriticalError(loggerSugar, constants.ErrInitializeDependencies, err)
		return
	}
	defer cleanup()

	loggerSugar.Infow(constants.SuccessToInitializeDependencies)

	newServer, err := server.NewHTTPServer(appDependencies, loggerSugar, &config.Setting)
	if err != nil {
		loggerSugar.Infow(constants.ErrStartServer)
		utils.HandleCriticalError(loggerSugar, constants.ErrStartServer, err)
		return
	}

	loggerSugar.Infow(constants.ServerStarted, constants.Port, newServer.Addr, constants.ContextPath, config.Setting.Server.Context)

	if err := newServer.ListenAndServe(); err != nil {
		utils.HandleCriticalError(loggerSugar, constants.ErrStartServer, err)
		loggerSugar.Infow(constants.ErrStartServer, constants.Error, err)
		return
	}
}
