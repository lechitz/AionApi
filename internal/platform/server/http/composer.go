// Package http
package http

import (
	"fmt"
	"net/http"

	// GraphQL adapter centralizado
	graphql "github.com/lechitz/AionApi/internal/adapter/primary/graphql"

	// Context adapters (HTTP)
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

// ComposeHandler cria o router HTTP de plataforma, aplica middlewares globais,
// registra handlers genéricos, monta módulos REST e expõe GraphQL em cfg.ServerGraphql.Path.
func ComposeHandler(cfg *config.Config, deps *bootstrap.AppDependencies, log logger.ContextLogger) (http.Handler, error) {
	r := chi.New() // ports.Router

	// Middlewares globais
	genericHandler := generic.New(log, cfg.General)
	r.Use(
		recovery.New(genericHandler), // sempre o mais externo
		requestid.New(),
	)

	// Defaults globais
	r.SetNotFound(http.HandlerFunc(genericHandler.NotFoundHandler))
	r.SetMethodNotAllowed(http.HandlerFunc(genericHandler.MethodNotAllowedHandler))
	r.SetError(genericHandler.ErrorHandler)

	// Prefixo de contexto (ex.: "/aion-api")
	apiPrefix := cfg.ServerHTTP.Context
	r.Group(apiPrefix, func(api ports.Router) {
		// Healthcheck
		api.GET("/health", http.HandlerFunc(genericHandler.HealthCheck))

		// Módulos REST (cada contexto registra suas rotas e a própria policy pública/protegida)
		if deps.AuthService != nil {
			ah := authhandler.New(deps.AuthService, cfg, log)
			authhandler.RegisterHTTP(api, ah)
		}
		if deps.UserService != nil {
			uh := userhandler.New(deps.UserService, cfg, log)
			userhandler.RegisterHTTPPublic(api, uh)
			if deps.AuthService != nil {
				userhandler.RegisterHTTPProtected(api, uh, deps.AuthService, log)
			}
		}

		// GraphQL (adapter central)
		gqlHandler, err := graphql.NewGraphqlHandler(deps.AuthService, deps.CategoryService, log, cfg)
		if err != nil {
			log.Errorw("failed to compose GraphQL handler", commonkeys.Error, err)
			return
		}
		api.Mount(cfg.ServerGraphql.Path, gqlHandler)
	})

	// Instrumentação de borda (OTel) — único wrapper no retorno
	h := otelhttp.NewHandler(
		r,
		fmt.Sprintf("%s-HTTP", cfg.Observability.OtelServiceName),
	)
	return h, nil
}
