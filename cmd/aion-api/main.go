package main

import (
	"github.com/lechitz/AionApi/adapters/primary/http/server"
	"github.com/lechitz/AionApi/internal/core/msg"
	"github.com/lechitz/AionApi/internal/platform/bootstrap"
	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/logger"
	"github.com/lechitz/AionApi/pkg/contextkeys"
	"github.com/lechitz/AionApi/pkg/errors"
)

func main() {

	loggerSugar, closeLogger := logger.InitLoggerSugar()
	defer closeLogger()

	loggerSugar.Infow(msg.StartingApplication)

	if err := config.LoadConfig(loggerSugar); err != nil {
		errors.HandleCriticalError(loggerSugar, msg.ErrToFailedLoadConfiguration, err)
		return
	}
	loggerSugar.Infow(msg.SuccessToLoadConfiguration, contextkeys.Setting, config.Setting)

	appDependencies, cleanup, err := bootstrap.InitializeDependencies(loggerSugar, config.Setting)
	if err != nil {
		errors.HandleCriticalError(loggerSugar, msg.ErrInitializeDependencies, err)
		return
	}
	defer cleanup()

	loggerSugar.Infow(msg.SuccessToInitializeDependencies)

	newServer, err := server.NewHTTPServer(appDependencies, loggerSugar, &config.Setting)
	if err != nil {
		loggerSugar.Infow(msg.ErrStartServer)
		errors.HandleCriticalError(loggerSugar, msg.ErrStartServer, nil)
		return
	}

	loggerSugar.Infow(msg.ServerStarted, contextkeys.Port, newServer.Addr, contextkeys.ContextPath, config.Setting.Server.Context)

	if err := newServer.ListenAndServe(); err != nil {
		errors.HandleCriticalError(loggerSugar, msg.ErrStartServer, err)
		loggerSugar.Infow(msg.ErrStartServer, contextkeys.Error, err)
		return
	}
}
