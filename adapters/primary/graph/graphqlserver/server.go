// Package graphqlserver implements a new HTTP server configured for handling GraphQL requests.
package graphqlserver

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/go-chi/chi/v5"

	"github.com/lechitz/AionApi/adapters/primary/graph"
	"github.com/lechitz/AionApi/adapters/primary/http/middleware/auth"
	"github.com/lechitz/AionApi/adapters/primary/http/middleware/recovery"
	"github.com/lechitz/AionApi/internal/infra/bootstrap"
	"github.com/lechitz/AionApi/internal/infra/config"
)

// NewGraphqlServer initializes and returns a new HTTP server
// configured to handle GraphQL requests using Chi router.
// Parameters:
//   - deps: AppDependencies container with services and repos
//   - cfg: Runtime configuration (includes HTTP port and secret key)
//
// Returns:
//   - *http.Server: GraphQL server instance ready to be started
//   - error: In case of setup failure
func NewGraphqlServer(deps *bootstrap.AppDependencies, cfg config.Config) (*http.Server, error) {
	router := chi.NewRouter()

	router.Use(auth.NewAuthMiddleware(
		deps.TokenRepository,
		deps.Logger,
		cfg.Secret.Key,
	).Auth)

	router.Use(recovery.RecoverMiddleware(deps.Logger))

	resolver := &graph.Resolver{
		CategoryService: deps.CategoryService,
		Logger:          deps.Logger,
	}

	srv := handler.New(
		graph.NewExecutableSchema(graph.Config{Resolvers: resolver}),
	)
	srv.AddTransport(transport.POST{})

	httpSrv := &http.Server{
		Addr:              ":" + cfg.ServerGraphql.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return httpSrv, nil
}
