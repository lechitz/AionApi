package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/lechitz/AionApi/internal/shared/commonkeys"

	"github.com/lechitz/AionApi/internal/platform/observability"

	"github.com/lechitz/AionApi/internal/core/ports/output"

	"github.com/lechitz/AionApi/cmd/aion-api/constants"
	"github.com/lechitz/AionApi/internal/adapters/primary/graph/graphqlserver"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/httpserver"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/response"
	"github.com/lechitz/AionApi/internal/platform/bootstrap"
	"github.com/lechitz/AionApi/internal/platform/config"
	loggerBuilder "github.com/lechitz/AionApi/pkg/logger"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// TODO:mover para as Configs.
type ServerConfig struct {
	Handler  http.Handler
	Name     string
	Addr     string
	Timeouts Timeouts
}

// TODO:mover para as Configs.
type Timeouts struct {
	Read  time.Duration
	Write time.Duration
}

func main() {
	logger, cleanupLogger := loggerBuilder.NewZapLogger()
	defer cleanupLogger()

	cfg := loadConfig(logger)

	appCtx, stopApp := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopApp()

	cleanupMetrics := observability.InitOtelMetrics(cfg, logger)
	defer cleanupMetrics()

	cleanupTracer := observability.InitTracer(cfg, logger)
	defer cleanupTracer()

	appDeps, cleanupDeps := initDependencies(appCtx, cfg, logger)
	defer cleanupDeps()

	servers := buildAllServers(appCtx, appDeps, cfg, logger)

	handleServers(appCtx, servers, cfg, logger, stopApp)
}

// loadConfig loads the application configuration from the given path.
func loadConfig(logger output.Logger) *config.Config {
	cfgLoader := config.NewLoader()
	cfg, err := cfgLoader.Load(logger)
	if err != nil {
		response.HandleCriticalError(logger, constants.ErrToFailedLoadConfiguration, err) // TODO: AJUSTAR ERRO.
		panic(err)
	}

	if err := cfg.Validate(); err != nil {
		response.HandleCriticalError(logger, constants.ErrInvalidConfiguration, err) // TODO: AJUSTAR ERRO.
		panic(err)
	}

	logger.Infof(constants.LoadedConfig, cfg)
	logger.Infow(constants.SuccessToLoadConfiguration)

	return cfg
}

// initDependencies initializes the application dependencies using the given configuration.
func initDependencies(appCtx context.Context, cfg *config.Config, logger output.Logger) (*bootstrap.AppDependencies, func()) {
	appDeps, cleanupResources, err := bootstrap.InitializeDependencies(appCtx, cfg, logger)
	if err != nil {
		response.HandleCriticalError(logger, constants.ErrInitializeDependencies, err) // TODO: AJUSTAR ERRO.
		panic(err)
	}

	logger.Infow(constants.SuccessToInitializeDependencies)

	return appDeps, cleanupResources
}

// setupHTTPHandler configures the HTTP router with middlewares, instrumentation and handlers.
func setupHTTPHandler(appDeps *bootstrap.AppDependencies, cfg *config.Config, logger output.Logger) http.Handler {
	router, err := httpserver.ComposeRouter(appDeps, cfg.ServerHTTP.Context)
	if err != nil {
		response.HandleCriticalError(logger, constants.ErrStartHTTPServer, err) // TODO: AJUSTAR ERRO.
		panic(err)
	}

	return otelhttp.NewHandler(router, cfg.Observability.OtelServiceName+"-REST")
}

// setupGraphQLHandler configures the GraphQL router with middlewares and handlers.
func setupGraphQLHandler(appDeps *bootstrap.AppDependencies, cfg *config.Config, logger output.Logger) http.Handler {
	handlerGraphQL, err := graphqlserver.NewGraphqlHandler(appDeps, cfg)
	if err != nil {
		response.HandleCriticalError(logger, constants.ErrStartGraphqlServer, err) // TODO: AJUSTAR ERRO.
		panic(err)
	}

	return handlerGraphQL
}

// buildServer builds a new HTTP server with the given configuration.
func buildServer(appCtx context.Context, serverConfig ServerConfig, logger output.Logger) *http.Server {
	srv := &http.Server{
		Addr:         serverConfig.Addr,
		BaseContext:  func(_ net.Listener) context.Context { return appCtx },
		Handler:      serverConfig.Handler,
		ReadTimeout:  serverConfig.Timeouts.Read,
		WriteTimeout: serverConfig.Timeouts.Write,
	}
	logger.Infow(fmt.Sprintf("%s server started", serverConfig.Name), commonkeys.Port, serverConfig.Addr)
	return srv
}

// buildAllServers builds all HTTP servers using the given configuration.
func buildAllServers(appCtx context.Context, appDeps *bootstrap.AppDependencies, cfg *config.Config, logger output.Logger) []*http.Server {
	serversConfig := []ServerConfig{
		{
			Name:    "HTTP", // TODO: Adicionar a Common.
			Addr:    fmt.Sprintf(":%s", cfg.ServerHTTP.Port),
			Handler: setupHTTPHandler(appDeps, cfg, logger),
			Timeouts: Timeouts{
				Read:  cfg.ServerHTTP.ReadTimeout,
				Write: cfg.ServerHTTP.WriteTimeout,
			},
		},
		{
			Name:    "GraphQL", // TODO: Adicionar a Common.
			Addr:    fmt.Sprintf(":%s", cfg.ServerGraphql.Port),
			Handler: setupGraphQLHandler(appDeps, cfg, logger),
			Timeouts: Timeouts{
				Read:  cfg.ServerGraphql.ReadTimeout,
				Write: cfg.ServerGraphql.WriteTimeout,
			},
		},
	}

	var servers []*http.Server
	for _, sc := range serversConfig {
		srv := buildServer(appCtx, sc, logger)
		servers = append(servers, srv)
	}

	return servers
}

// handleServers starts the given HTTP servers and handles.
func handleServers(appCtx context.Context, servers []*http.Server, cfg *config.Config, logger output.Logger, stop context.CancelFunc) {
	var wg sync.WaitGroup
	errChan := make(chan error, len(servers))

	for _, srv := range servers {
		wg.Add(1)
		go func(s *http.Server) {
			defer wg.Done()
			if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errChan <- fmt.Errorf("failed to start server on %s: %w", s.Addr, err) // TODO: AJUSTAR ERRO.
			}
		}(srv)
	}

	select {
	case err := <-errChan:
		logger.Errorw(constants.MsgUnexpectedServerFailure, commonkeys.Error, err.Error())
		response.HandleCriticalError(logger, constants.MsgUnexpectedServerFailure, err) // TODO: AJUSTAR ERRO.
		stop()
	case <-appCtx.Done():
		logger.Infow(constants.MsgShutdownSignalReceived)
	}

	shutdownTimeout := cfg.Application.Timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	for _, srv := range servers {
		_ = srv.Shutdown(shutdownCtx)
	}

	wg.Wait()
}
