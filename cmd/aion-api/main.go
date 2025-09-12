// Package main
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

func main() {
	logs, cleanupLogger := contextlogger.New()
	defer cleanupLogger()

	keyGenerator := crypto.New()
	cfg := loadConfig(keyGenerator, logs)

	appCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cleanupMetrics := metric.InitOtelMetrics(cfg, logs)
	defer cleanupMetrics()

	cleanupTracer := tracer.InitTracer(cfg, logs)
	defer cleanupTracer()

	appDeps, cleanupDeps := initDependencies(appCtx, cfg, logs)
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Application.Timeout)
		defer cancel()
		cleanupDeps(shutdownCtx)
	}()

	if err := server.RunAll(appCtx, cfg, appDeps, logs); err != nil {
		logs.Errorw(ErrServerRunFailed, commonkeys.Error, err.Error())
		os.Exit(1)
	}
}

func loadConfig(keyGenerator keygen.Generator, logs logger.ContextLogger) *config.Config {
	cfg, err := config.New(keyGenerator).Load(logs)
	if err != nil {
		logs.Errorw(ErrLoadConfig, commonkeys.Error, err.Error())
		os.Exit(1)
	}

	if err := cfg.Validate(); err != nil {
		logs.Errorw(ErrInvalidConfig, commonkeys.Error, err.Error())
		os.Exit(1)
	}
	logs.Infow(MsgConfigLoaded, commonkeys.APIName, cfg.General.Name, commonkeys.AppEnv, cfg.General.Env, commonkeys.AppVersion, cfg.General.Version)
	return cfg
}

func initDependencies(appCtx context.Context, cfg *config.Config, logs logger.ContextLogger) (*bootstrap.AppDependencies, func(context.Context)) {
	deps, cleanup, err := bootstrap.InitializeDependencies(appCtx, cfg, logs)
	if err != nil {
		logs.Errorw(ErrInitDeps, commonkeys.Error, err.Error())
		os.Exit(1)
	}
	logs.Infow(MsgDepsInitialized)
	return deps, cleanup
}
