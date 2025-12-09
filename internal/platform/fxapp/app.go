package fxapp

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/fx"

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

// ProvideLogger builds the structured logger and registers its cleanup on shutdown.
func ProvideLogger(lc fx.Lifecycle) logger.ContextLogger {
	logs, cleanup := contextlogger.New()
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			cleanup()
			return nil
		},
	})
	return logs
}

// ProvideKeyGenerator exposes the crypto-based key generator.
func ProvideKeyGenerator() keygen.Generator {
	return crypto.New()
}

// ProvideConfig loads and validates configuration using the provided key generator and logger.
func ProvideConfig(logs logger.ContextLogger, keyGen keygen.Generator) (*config.Config, error) {
	cfg, err := config.New(keyGen).Load(logs)
	if err != nil {
		return nil, err
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	logs.Infow(
		"configuration loaded",
		commonkeys.APIName, cfg.General.Name,
		commonkeys.AppEnv, cfg.General.Env,
		commonkeys.AppVersion, cfg.General.Version,
	)
	return cfg, nil
}

// InitObservability wires OTEL metrics and traces, cleaning up on shutdown.
func InitObservability(lc fx.Lifecycle, cfg *config.Config, logs logger.ContextLogger) {
	cleanupMetrics := metric.InitOtelMetrics(cfg, logs)
	cleanupTracer := tracer.InitTracer(cfg, logs)

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			cleanupMetrics()
			cleanupTracer()
			return nil
		},
	})
}

// ProvideDependencies composes repositories/usecases and ensures they are cleaned up gracefully.
func ProvideDependencies(lc fx.Lifecycle, cfg *config.Config, logs logger.ContextLogger) (*bootstrap.AppDependencies, error) {
	deps, cleanup, err := bootstrap.InitializeDependencies(context.Background(), cfg, logs)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			shutdownCtx, cancel := context.WithTimeout(ctx, cfg.Application.Timeout)
			defer cancel()
			cleanup(shutdownCtx)
			return nil
		},
	})

	logs.Infow("dependencies initialized")
	return deps, nil
}

// RunServers starts all servers (REST/GraphQL) and stops them on shutdown.
func RunServers(lc fx.Lifecycle, cfg *config.Config, deps *bootstrap.AppDependencies, logs logger.ContextLogger) {
	var (
		cancel context.CancelFunc
		done   chan error
	)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			runCtx, c := context.WithCancel(ctx)
			cancel = c
			done = make(chan error, 1)

			go func() {
				done <- server.RunAll(runCtx, cfg, deps, logs)
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if cancel != nil {
				cancel()
			}

			if done == nil {
				return nil
			}

			select {
			case err := <-done:
				if err != nil && !errors.Is(err, context.Canceled) {
					return fmt.Errorf("server run failed: %w", err)
				}
			case <-ctx.Done():
				return ctx.Err()
			}
			return nil
		},
	})
}
