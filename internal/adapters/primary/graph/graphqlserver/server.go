// Package graphqlserver implements a new HTTP server configured for handling GraphQL requests.
package graphqlserver

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/go-chi/chi/v5"
	"github.com/lechitz/AionApi/internal/adapters/primary/graph"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/auth"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/recovery"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/lechitz/AionApi/internal/platform/bootstrap"
	"github.com/lechitz/AionApi/internal/platform/config"
)

// NewGraphqlHandler initializes and returns an HTTP handler for GraphQL requests.
func NewGraphqlHandler(deps *bootstrap.AppDependencies, cfg *config.Config) (http.Handler, error) {
	router := chi.NewRouter()

	router.Use(auth.NewAuthMiddleware(deps.TokenRepository, deps.Logger, deps.TokenClaimsExtractor).Auth)
	router.Use(recovery.RecoverMiddleware(deps.Logger))

	resolver := &graph.Resolver{
		CategoryService: deps.CategoryService,
		Logger:          deps.Logger,
	}

	srv := handler.New(
		graph.NewExecutableSchema(graph.Config{Resolvers: resolver}),
	)
	srv.AddTransport(transport.POST{})

	router.Handle(cfg.ServerGraphql.Path, otelhttp.NewHandler(srv, cfg.Observability.OtelServiceName+"-GraphQL"))

	return router, nil
}
