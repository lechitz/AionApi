// Package main is the entry point for the application.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/lechitz/AionApi/internal/adapter/secondary/contextlogger"
	"github.com/lechitz/AionApi/internal/adapter/secondary/crypto"
	"github.com/lechitz/AionApi/internal/platform/bootstrap"
	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/observability/metric"
	"github.com/lechitz/AionApi/internal/platform/observability/tracer"
	"github.com/lechitz/AionApi/internal/platform/ports/output/keygen"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/internal/platform/server"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
)

// the main is the entry point for the application.
func main() {
	os.Exit(run())
}

// run is the main application logic.
func run() int {
	logs, cleanupLogger := contextlogger.New()
	defer cleanupLogger()

	keyGenerator := crypto.New()

	cfg, err := loadConfig(keyGenerator, logs)
	if err != nil {
		logs.Errorw(ErrLoadConfig, commonkeys.Error, err.Error())
		return 2
	}

	appCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cleanupMetrics := metric.InitOtelMetrics(cfg, logs)
	defer cleanupMetrics()

	cleanupTracer := tracer.InitTracer(cfg, logs)
	defer cleanupTracer()

	appDeps, cleanupDeps, err := initDependencies(appCtx, cfg, logs)
	if err != nil {
		logs.Errorw(ErrInitDeps, commonkeys.Error, err.Error())
		return 3
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Application.Timeout)
		defer cancel()
		cleanupDeps(shutdownCtx)
	}()

	if err := server.RunAll(appCtx, cfg, appDeps, logs); err != nil {
		logs.Errorw(ErrServerRunFailed, commonkeys.Error, err.Error())
		return 1
	}
	return 0
}

// loadConfig loads the application configuration.
func loadConfig(keyGenerator keygen.Generator, logs logger.ContextLogger) (*config.Config, error) {
	cfg, err := config.New(keyGenerator).Load(logs)
	if err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	logs.Infow(
		MsgConfigLoaded,
		commonkeys.APIName, cfg.General.Name,
		commonkeys.AppEnv, cfg.General.Env,
		commonkeys.AppVersion, cfg.General.Version,
	)
	return cfg, nil
}

// initDependencies initializes the application dependencies.
func initDependencies(appCtx context.Context, cfg *config.Config, logs logger.ContextLogger) (*bootstrap.AppDependencies, func(context.Context), error) {
	deps, cleanup, err := bootstrap.InitializeDependencies(appCtx, cfg, logs)
	if err != nil {
		return nil, nil, err
	}
	logs.Infow(MsgDepsInitialized)
	return deps, cleanup, nil
}
