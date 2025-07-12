// Package graphqlserver implements a new HTTP server configured for handling GraphQL requests.
package graphqlserver

import (
	"net/http"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/generic"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/go-chi/chi/v5"
	"github.com/lechitz/AionApi/internal/adapters/primary/graph"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/authmiddleware"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/recoverymiddleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/lechitz/AionApi/internal/platform/bootstrap"
	"github.com/lechitz/AionApi/internal/platform/config"
)

// NewGraphqlHandler initializes and returns an HTTP handler for GraphQL requests.
func NewGraphqlHandler(deps *bootstrap.AppDependencies, cfg *config.Config) (http.Handler, error) {
	router := chi.NewRouter()

	router.Use(authmiddleware.New(deps.TokenRepository, deps.Logger, deps.TokenClaimsExtractor).Auth)
	genericHandler := generic.New(deps.Logger, cfg.General)
	router.Use(recoverymiddleware.New(genericHandler))

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
