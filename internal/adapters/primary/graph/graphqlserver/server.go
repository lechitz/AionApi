// Package graphqlserver implements a new HTTP server configured for handling GraphQL requests.
package graphqlserver

import (
	"context"
	"net/http"
	"time"

	"github.com/lechitz/AionApi/internal/adapters/primary/graph"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/auth"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/recovery"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/lechitz/AionApi/internal/platform/bootstrap"
	"github.com/lechitz/AionApi/internal/platform/config"
)

// NewGraphqlServer initializes and returns a new HTTP server
// configured to handle GraphQL requests using Chi router.
func NewGraphqlServer(deps *bootstrap.AppDependencies, cfg config.Config) (*http.Server, error) {
	router := chi.NewRouter()

	router.Use(auth.NewAuthMiddleware(deps.TokenRepository, deps.Logger, cfg.Secret.Key).Auth)

	router.Use(recovery.RecoverMiddleware(deps.Logger))

	resolver := &graph.Resolver{
		CategoryService: deps.CategoryService,
		Logger:          deps.Logger,
	}

	srv := handler.New(
		graph.NewExecutableSchema(graph.Config{Resolvers: resolver}),
	)
	srv.AddTransport(transport.POST{})

	// Middleware to propagate context like userID from HTTP middleware to gqlgen resolvers
	srv.AroundOperations(
		func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
			return next(ctx)
		},
	)

	router.Handle("/graphql", otelhttp.NewHandler(srv, "AionApi-GraphQL"))

	httpSrv := &http.Server{
		Addr:              ":" + cfg.ServerGraphql.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return httpSrv, nil
}
