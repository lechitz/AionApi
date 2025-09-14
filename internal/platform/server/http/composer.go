// Package http provides the HTTP server composition for all adapters and routes.
package http

import (
	"fmt"
	"net/http"

	"github.com/lechitz/AionApi/internal/adapter/primary/graphql"

	authhandler "github.com/lechitz/AionApi/internal/auth/adapter/primary/http/handler"
	userhandler "github.com/lechitz/AionApi/internal/user/adapter/primary/http/handler"

	"github.com/lechitz/AionApi/internal/platform/bootstrap"
	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	genericHandler "github.com/lechitz/AionApi/internal/platform/server/http/generic/handler"
	"github.com/lechitz/AionApi/internal/platform/server/http/middleware/recovery"
	"github.com/lechitz/AionApi/internal/platform/server/http/middleware/requestid"
	"github.com/lechitz/AionApi/internal/platform/server/http/ports"
	"github.com/lechitz/AionApi/internal/platform/server/http/router/chi"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Log/route/name constants used in this file.
const (
	// RouteHealth is the health-check endpoint path.
	RouteHealth = "/health"

	// LogErrComposeGraphQL is the log message when composing the GraphQL handler fails.
	LogErrComposeGraphQL = "failed to compose GraphQL handler"

	// OTelHTTPHandlerNameFormat formats the OpenTelemetry HTTP handler name with the service name.
	OTelHTTPHandlerNameFormat = "%s-HTTP"
)

// ComposeHandler assembles the HTTP handler with platform middlewares, domain routes, and GraphQL.
// It returns an otel-instrumented http.Handler ready to be served by an HTTP server.
func ComposeHandler(cfg *config.Config, deps *bootstrap.AppDependencies, log logger.ContextLogger) (http.Handler, error) {
	r := chi.New() // ports.Router

	// Global middlewares (recovery should be the outermost).
	gh := genericHandler.New(log, cfg.General)
	r.Use(
		recovery.New(gh),
		requestid.New(),
	)

	// Default handlers.
	r.SetNotFound(http.HandlerFunc(gh.NotFoundHandler))
	r.SetMethodNotAllowed(http.HandlerFunc(gh.MethodNotAllowedHandler))
	r.SetError(gh.ErrorHandler)

	apiPrefix := cfg.ServerHTTP.Context
	r.Group(apiPrefix, func(api ports.Router) {
		// Health check.
		api.GET(RouteHealth, http.HandlerFunc(gh.HealthCheck))

		// Auth REST endpoints.
		if deps.AuthService != nil {
			ah := authhandler.New(deps.AuthService, cfg, log)
			authhandler.RegisterHTTP(api, ah)
		}

		// User REST endpoints.
		if deps.UserService != nil {
			uh := userhandler.New(deps.UserService, cfg, log)
			userhandler.RegisterHTTP(api, uh, deps.AuthService, log)
		}

		// GraphQL endpoint.
		gqlHandler, err := graphql.NewGraphqlHandler(deps.AuthService, deps.CategoryService, log, cfg)
		if err != nil {
			log.Errorw(LogErrComposeGraphQL, commonkeys.Error, err)
			return
		}
		api.Mount(cfg.ServerGraphql.Path, gqlHandler)
	})

	// Wrap router with OpenTelemetry instrumentation.
	h := otelhttp.NewHandler(
		r,
		fmt.Sprintf(OTelHTTPHandlerNameFormat, cfg.Observability.OtelServiceName),
	)
	return h, nil
}
