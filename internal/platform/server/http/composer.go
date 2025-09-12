// Package http is the HTTP server implementation.
package http

import (
	"fmt"
	"net/http"

	graphql "github.com/lechitz/AionApi/internal/adapter/primary/graphql"

	authhandler "github.com/lechitz/AionApi/internal/auth/adapter/primary/http/handler"
	userhandler "github.com/lechitz/AionApi/internal/user/adapter/primary/http/handler"

	"github.com/lechitz/AionApi/internal/platform/bootstrap"
	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	generic "github.com/lechitz/AionApi/internal/platform/server/http/generic/handler"
	"github.com/lechitz/AionApi/internal/platform/server/http/middleware/recovery"
	"github.com/lechitz/AionApi/internal/platform/server/http/middleware/requestid"
	"github.com/lechitz/AionApi/internal/platform/server/http/ports"
	"github.com/lechitz/AionApi/internal/platform/server/http/router/chi"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// ComposeHandler create a new HTTP server.
func ComposeHandler(cfg *config.Config, deps *bootstrap.AppDependencies, log logger.ContextLogger) (http.Handler, error) {
	r := chi.New() // ports.Router

	// Global middlewares
	genericHandler := generic.New(log, cfg.General)
	r.Use(
		recovery.New(genericHandler), // sempre o mais externo
		requestid.New(),
	)

	// Default handlers
	r.SetNotFound(http.HandlerFunc(genericHandler.NotFoundHandler))
	r.SetMethodNotAllowed(http.HandlerFunc(genericHandler.MethodNotAllowedHandler))
	r.SetError(genericHandler.ErrorHandler)

	apiPrefix := cfg.ServerHTTP.Context
	r.Group(apiPrefix, func(api ports.Router) {
		api.GET("/health", http.HandlerFunc(genericHandler.HealthCheck))

		if deps.AuthService != nil {
			ah := authhandler.New(deps.AuthService, cfg, log)
			authhandler.RegisterHTTP(api, ah)
		}
		if deps.UserService != nil {
			uh := userhandler.New(deps.UserService, cfg, log)
			userhandler.RegisterHTTP(api, uh, deps.AuthService, log)
		}

		gqlHandler, err := graphql.NewGraphqlHandler(deps.AuthService, deps.CategoryService, log, cfg)
		if err != nil {
			log.Errorw("failed to compose GraphQL handler", commonkeys.Error, err)
			return
		}
		api.Mount(cfg.ServerGraphql.Path, gqlHandler)
	})

	h := otelhttp.NewHandler(
		r,
		fmt.Sprintf("%s-HTTP", cfg.Observability.OtelServiceName),
	)
	return h, nil
}
