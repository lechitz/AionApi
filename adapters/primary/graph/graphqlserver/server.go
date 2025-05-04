package graphqlserver

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/go-chi/chi/v5"
	"github.com/lechitz/AionApi/adapters/primary/graph"
	"github.com/lechitz/AionApi/adapters/primary/http/middleware/auth"
	"github.com/lechitz/AionApi/adapters/primary/http/middleware/recovery"
	"github.com/lechitz/AionApi/internal/infra/bootstrap"
	"github.com/lechitz/AionApi/internal/infra/config"
	"net/http"
)

func NewGraphqlServer(deps *bootstrap.AppDependencies) (*http.Server, error) {
	router := chi.NewRouter()

	authMiddleware := auth.NewAuthMiddleware(deps.TokenRepository, deps.Logger)
	router.Use(authMiddleware.Auth)
	router.Use(recovery.RecoverMiddleware(deps.Logger))

	schema := graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			CategoryService: deps.CategoryService,
			Logger:          deps.Logger,
		},
	})

	srv := handler.New(schema)
	srv.AddTransport(transport.POST{})

	router.Post("/graphql", srv.ServeHTTP)

	httpServer := &http.Server{
		Addr:    ":" + config.Setting.ServerGraphql.Port,
		Handler: router,
	}

	return httpServer, nil
}
