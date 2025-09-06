package graphqlserver

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/go-chi/chi/v5"
	"github.com/lechitz/AionApi/internal/auth/adapter/primary/http/middleware"
	authInput "github.com/lechitz/AionApi/internal/auth/core/ports/input"
	categoryInput "github.com/lechitz/AionApi/internal/category/core/ports/input"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/internal/platform/server/graph"
	"github.com/lechitz/AionApi/internal/platform/server/graph/resolver"
	generic "github.com/lechitz/AionApi/internal/platform/server/http/generic/handler"
	"github.com/lechitz/AionApi/internal/platform/server/http/middleware/recoverymiddleware"

	"github.com/lechitz/AionApi/internal/platform/config"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func NewGraphqlHandler(authService authInput.AuthService, categoryService categoryInput.CategoryService, logger logger.ContextLogger, cfg *config.Config) (http.Handler, error) {

	router := chi.NewRouter()

	router.Use(middleware.New(authService, logger).Auth)

	genericHandler := generic.New(logger, cfg.General)
	router.Use(recoverymiddleware.New(genericHandler))

	resolver := &resolver.Resolver{
		CategoryService: categoryService,
		Logger:          logger,
	}

	es := graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver,
		Directives: graph.DirectiveRoot{
			Auth: middleware.Auth(),
		},
	})

	srv := handler.New(es)

	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.MultipartForm{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})

	router.Handle(
		cfg.ServerGraphql.Path,
		otelhttp.NewHandler(srv, cfg.Observability.OtelServiceName+"-GraphQL"),
	)
	return router, nil
}
