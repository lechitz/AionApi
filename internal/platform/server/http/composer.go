// Package http provides the HTTP server composition for all adapters and routes.
package http

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/lechitz/AionApi/internal/adapter/primary/graphql"

	authhandler "github.com/lechitz/AionApi/internal/auth/adapter/primary/http/handler"
	userhandler "github.com/lechitz/AionApi/internal/user/adapter/primary/http/handler"

	"github.com/lechitz/AionApi/internal/platform/bootstrap"
	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	genericHandler "github.com/lechitz/AionApi/internal/platform/server/http/generic/handler"
	"github.com/lechitz/AionApi/internal/platform/server/http/middleware/cors"
	"github.com/lechitz/AionApi/internal/platform/server/http/middleware/recovery"
	"github.com/lechitz/AionApi/internal/platform/server/http/middleware/requestid"
	"github.com/lechitz/AionApi/internal/platform/server/http/ports"
	"github.com/lechitz/AionApi/internal/platform/server/http/router/chi"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"

	httpSwagger "github.com/swaggo/http-swagger" // use Handler(...opts) which returns http.Handler
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// ComposeHandler assembles the HTTP handler with platform middlewares, domain routes, Swagger UI and GraphQL.
func ComposeHandler(cfg *config.Config, deps *bootstrap.AppDependencies, log logger.ContextLogger) (http.Handler, error) {
	r := chi.New() // ports.Router

	// Global middlewares
	gh := genericHandler.New(log, cfg.General)
	r.Use(
		recovery.New(gh),
		requestid.New(),
		cors.New(),
	)

	// Default handlers
	r.SetNotFound(http.HandlerFunc(gh.NotFoundHandler))
	r.SetMethodNotAllowed(http.HandlerFunc(gh.MethodNotAllowedHandler))
	r.SetError(gh.ErrorHandler)

	// Resolve context and mount points from config with safe fallbacks
	apiContext := cfg.ServerHTTP.Context
	if apiContext == "" {
		apiContext = "/"
	}

	swaggerMount := cfg.ServerHTTP.SwaggerMountPath
	if swaggerMount == "" {
		swaggerMount = DefaultSwaggerMountPath
	}

	docsAlias := cfg.ServerHTTP.DocsAliasPath
	if docsAlias == "" {
		docsAlias = DefaultDocsAliasPath
	}

	routeHealth := cfg.ServerHTTP.HealthRoute
	if routeHealth == "" {
		routeHealth = DefaultRouteHealth
	}

	r.Group(apiContext, func(api ports.Router) {
		// Swagger UI + OpenAPI JSON mounted under the API context
		// httpSwagger.Handler returns an http.Handler suitable for Mount.
		swaggerDocURL := path.Clean(apiContext + "/" +
			strings.TrimPrefix(swaggerMount, "/") + "/" + DefaultSwaggerDocFile)

		api.Mount(swaggerMount, httpSwagger.Handler(
			httpSwagger.URL(swaggerDocURL), // UI loads doc.json from this URL
			httpSwagger.DeepLinking(true),
			httpSwagger.DocExpansion("none"),
			httpSwagger.DomID("swagger-ui"),
		))

		// Friendly alias: {apiContext}{docsAlias} -> {apiContext}{swaggerMount}/{swaggerIndex}
		api.GET(docsAlias, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, path.Join(apiContext, swaggerMount, DefaultSwaggerIndexFile), http.StatusTemporaryRedirect)
		}))

		// Group API Root (ex.: /api/v1)
		api.Group(cfg.ServerHTTP.APIRoot, func(v1 ports.Router) {
			// REST endpoints
			if deps.AuthService != nil {
				ah := authhandler.New(deps.AuthService, cfg, log)
				authhandler.RegisterHTTP(v1, ah)
			}

			if deps.UserService != nil {
				uh := userhandler.New(deps.UserService, cfg, log)
				userhandler.RegisterHTTP(v1, uh, deps.AuthService, log)
			}

			// GraphQL endpoint
			gqlHandler, err := graphql.NewGraphqlHandler(deps.AuthService, deps.CategoryService, deps.TagService, log, cfg)
			if err != nil {
				log.Errorw(LogErrComposeGraphQL, commonkeys.Error, err)
				return
			}

			v1.Mount(cfg.ServerGraphql.Path, gqlHandler)
		})
	})

	// OpenTelemetry HTTP wrapper: instrument the main router but expose health route uninstrumented
	instrumented := otelhttp.NewHandler(
		r,
		fmt.Sprintf(OTelHTTPHandlerNameFormat, cfg.Observability.OtelServiceName),
	)

	mux := http.NewServeMux()
	p := path.Clean(apiContext + "/" + strings.TrimPrefix(routeHealth, "/"))
	mux.Handle(p, http.HandlerFunc(gh.HealthCheck))
	mux.Handle(p+"/", http.HandlerFunc(gh.HealthCheck))

	// Backwards compatibility: also expose health under {apiContext}{APIRoot}{routeHealth} (e.g., /aion/api/v1/health)
	altHealth := path.Clean(apiContext + "/" + strings.TrimPrefix(cfg.ServerHTTP.APIRoot, "/") + "/" + strings.TrimPrefix(routeHealth, "/"))
	mux.Handle(altHealth, http.HandlerFunc(gh.HealthCheck))
	mux.Handle(altHealth+"/", http.HandlerFunc(gh.HealthCheck))

	mux.Handle("/", instrumented)

	return mux, nil
}
