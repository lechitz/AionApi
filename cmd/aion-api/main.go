package main

import (
	"github.com/lechitz/AionApi/adapters/input/http/server"
	"github.com/lechitz/AionApi/core/msg"
	"github.com/lechitz/AionApi/internal/platform/bootstrap"
	config2 "github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/logger"
	"github.com/lechitz/AionApi/pkg/contextkeys"
	"github.com/lechitz/AionApi/pkg/errors"
)

func main() {

	loggerSugar, closeLogger := logger.InitLoggerSugar()
	defer closeLogger()

	loggerSugar.Infow(msg.StartingApplication)

	if err := config2.LoadConfig(loggerSugar); err != nil {
		errors.HandleCriticalError(loggerSugar, msg.ErrToFailedLoadConfiguration, err)
		return
	}
	loggerSugar.Infow(msg.SuccessToLoadConfiguration, contextkeys.Setting, config2.Setting)

	appDependencies, cleanup, err := bootstrap.InitializeDependencies(loggerSugar, config2.Setting)
	if err != nil {
		errors.HandleCriticalError(loggerSugar, msg.ErrInitializeDependencies, err)
		return
	}
	defer cleanup()

	loggerSugar.Infow(msg.SuccessToInitializeDependencies)

	newServer, err := server.NewHTTPServer(appDependencies, loggerSugar, &config2.Setting)
	if err != nil {
		loggerSugar.Infow(msg.ErrStartServer)
		errors.HandleCriticalError(loggerSugar, msg.ErrStartServer, nil)
		return
	}

	loggerSugar.Infow(msg.ServerStarted, contextkeys.Port, newServer.Addr, contextkeys.ContextPath, config2.Setting.Server.Context)

	if err := newServer.ListenAndServe(); err != nil {
		errors.HandleCriticalError(loggerSugar, msg.ErrStartServer, err)
		loggerSugar.Infow(msg.ErrStartServer, contextkeys.Error, err)
		return
	}
}
