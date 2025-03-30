package main

import (
	"github.com/lechitz/AionApi/internal/adapters/primary/http/server"
	"github.com/lechitz/AionApi/internal/platform/bootstrap"
	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/logger"
	"github.com/lechitz/AionApi/pkg/utils"
)

func main() {

	loggerSugar, closeLogger := logger.InitLoggerSugar()
	defer closeLogger()

	loggerSugar.Infow(StartingApplication)

	if err := config.LoadConfig(loggerSugar); err != nil {
		utils.HandleCriticalError(loggerSugar, ErrToFailedLoadConfiguration, err)
		return
	}
	loggerSugar.Infow(SuccessToLoadConfiguration, Settings, config.Setting)

	appDependencies, cleanup, err := bootstrap.InitializeDependencies(loggerSugar, config.Setting)
	if err != nil {
		utils.HandleCriticalError(loggerSugar, ErrInitializeDependencies, err)
		return
	}
	defer cleanup()

	loggerSugar.Infow(SuccessToInitializeDependencies)

	newServer, err := server.NewHTTPServer(appDependencies, loggerSugar, &config.Setting)
	if err != nil {
		loggerSugar.Infow(ErrStartServer)
		utils.HandleCriticalError(loggerSugar, ErrStartServer, nil)
		return
	}

	loggerSugar.Infow(ServerStarted, Port, newServer.Addr, ContextPath, config.Setting.Server.Context)

	if err := newServer.ListenAndServe(); err != nil {
		utils.HandleCriticalError(loggerSugar, ErrStartServer, err)
		loggerSugar.Infow(ErrStartServer, Error, err)
		return
	}
}
