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

// InfraModule bundles core infrastructure providers (logger, config, tracer/metrics, cache, database, http client).
//
//nolint:gochecknoglobals // Fx modules are intended as package-level options.
var InfraModule = fx.Options(
	fx.Provide(
		ProvideLogger,
		ProvideKeyGenerator,
		ProvideConfig,
		ProvideCache,
		ProvideDatabase,
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

// CacheOut groups all cache instances for dependency injection.
type CacheOut struct {
	fx.Out

	AuthCache     cache.Cache `name:"authCache"`
	CategoryCache cache.Cache `name:"categoryCache"`
	TagCache      cache.Cache `name:"tagCache"`
	RecordCache   cache.Cache `name:"recordCache"`
	UserCache     cache.Cache `name:"userCache"`
}

// ProvideCache initializes all cache instances for bounded contexts.
// Uses Redis as the underlying implementation with isolated databases per context.
func ProvideCache(lc fx.Lifecycle, cfg *config.Config, log logger.ContextLogger) (CacheOut, error) {
	ctx := context.Background()

	// Track all created caches for cleanup on error
	var createdCaches []cache.Cache
	cleanup := func() {
		for _, c := range createdCaches {
			_ = c.Close()
		}
	}

	// Helper to create a cache instance
	createCache := func(contextName string, dbNum int) (cache.Cache, error) {
		cacheClient, err := cacheRedis.NewConnection(ctx, cfg.Cache, dbNum, log)
		if err != nil {
			log.Errorw("failed to create cache", commonkeys.Error, err, "context", contextName, "db", dbNum)
			cleanup() // Cleanup any previously created caches
			return nil, err
		}
		log.Infow("cache initialized", "context", contextName, "db", dbNum)
		createdCaches = append(createdCaches, cacheClient)
		return cacheClient, nil
	}

	authCache, err := createCache("auth", cfg.Cache.AuthDB)
	if err != nil {
		return CacheOut{}, err
	}

	categoryCache, err := createCache("category", cfg.Cache.CategoryDB)
	if err != nil {
		return CacheOut{}, err
	}

	tagCache, err := createCache("tag", cfg.Cache.TagDB)
	if err != nil {
		return CacheOut{}, err
	}

	recordCache, err := createCache("record", cfg.Cache.RecordDB)
	if err != nil {
		return CacheOut{}, err
	}

	userCache, err := createCache("user", cfg.Cache.UserDB)
	if err != nil {
		return CacheOut{}, err
	}

	// Register shutdown hook to close all caches
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			cleanup()
			log.Infow("all cache connections closed")
			return nil
		},
	})

	return CacheOut{
		AuthCache:     authCache,
		CategoryCache: categoryCache,
		TagCache:      tagCache,
		RecordCache:   recordCache,
		UserCache:     userCache,
	}, nil
}

// ProvideDatabase initializes the database connection and closes it on shutdown.
// Returns db.DB interface following Dependency Inversion Principle.
// Uses PostgreSQL as the underlying implementation.
func ProvideDatabase(lc fx.Lifecycle, cfg *config.Config, log logger.ContextLogger) (db.DB, error) {
	conn, err := postgres.NewConnection(context.Background(), cfg.DB, log)
	if err != nil {
		log.Errorw("failed to connect to database", commonkeys.Error, err)
		return nil, err
	}

	// Wrap concrete DB with adapter to return interface
	dbAdapter := postgres.NewDBAdapter(conn)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
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

	log.Infow("database connection initialized", "type", "postgresql")
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
