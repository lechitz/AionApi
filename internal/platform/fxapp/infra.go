// Package fxapp wires the application using Uber Fx modules.
package fxapp

import (
	"context"
	"time"

	cacheRedis "github.com/lechitz/AionApi/internal/adapter/secondary/cache/redis"
	"github.com/lechitz/AionApi/internal/adapter/secondary/contextlogger"
	"github.com/lechitz/AionApi/internal/adapter/secondary/crypto"
	"github.com/lechitz/AionApi/internal/adapter/secondary/db/postgres"
	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/httpclient"
	"github.com/lechitz/AionApi/internal/platform/observability/metric"
	"github.com/lechitz/AionApi/internal/platform/observability/tracer"
	"github.com/lechitz/AionApi/internal/platform/ports/output/cache"
	"github.com/lechitz/AionApi/internal/platform/ports/output/db"
	httpclientPort "github.com/lechitz/AionApi/internal/platform/ports/output/httpclient"
	"github.com/lechitz/AionApi/internal/platform/ports/output/keygen"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.uber.org/fx"
)

// InfraModule bundles core infrastructure providers (logger, config, tracer/metrics, redis, postgres).
//
//nolint:gochecknoglobals // Fx modules are intended as package-level options.
var InfraModule = fx.Options(
	fx.Provide(
		ProvideLogger,
		ProvideKeyGenerator,
		ProvideConfig,
		ProvideRedisClient,
		ProvidePostgres,
		ProvideHTTPClient,
	),
	fx.Invoke(InitObservability),
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

// ProvideRedisClient initializes the Redis client and closes it on shutdown.
func ProvideRedisClient(lc fx.Lifecycle, cfg *config.Config, log logger.ContextLogger) (cache.Cache, error) {
	client, err := cacheRedis.NewConnection(context.Background(), cfg.Cache, log)
	if err != nil {
		log.Errorw("failed to connect to Redis", commonkeys.Error, err)
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return client.Close()
		},
	})
	return client, nil
}

// ProvidePostgres initializes the DB connection and closes it on shutdown.
// Returns db.DB interface (not concrete *gorm.DB) following Dependency Inversion Principle.
func ProvidePostgres(lc fx.Lifecycle, cfg *config.Config, log logger.ContextLogger) (db.DB, error) {
	conn, err := postgres.NewConnection(context.Background(), cfg.DB, log)
	if err != nil {
		log.Errorw("failed to connect to postgres", commonkeys.Error, err)
		return nil, err
	}

	// Wrap gorm.DB with adapter to return interface
	dbAdapter := postgres.NewDBAdapter(conn)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			// Respect shutdown timeout if provided
			timeout := cfg.DB.ConnMaxLifetime
			if timeout <= 0 {
				timeout = 5 * time.Second
			}
			shutdownCtx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			postgres.Close(conn, log)
			<-shutdownCtx.Done()
			return nil
		},
	})

	return dbAdapter, nil
}

// ProvideHTTPClient creates an instrumented outbound HTTP client and returns it as the HTTPClient interface.
// Uses timeout from AionChat config if available, otherwise defaults to 15s.
// The returned client is instrumented with OTEL for automatic tracing and context propagation.
func ProvideHTTPClient(cfg *config.Config) httpclientPort.HTTPClient {
	timeout := 15 * time.Second
	if cfg.AionChat.Timeout > 0 {
		timeout = cfg.AionChat.Timeout
	}

	opts := httpclient.Options{
		Timeout: timeout,
	}

	// NewInstrumentedClient returns *http.Client with OTEL instrumentation
	instrumentedHTTPClient := httpclient.NewInstrumentedClient(opts)

	// Wrap with adapter that implements HTTPClient interface
	return httpclient.NewClient(instrumentedHTTPClient)
}
