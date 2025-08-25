package graphqlserver

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/go-chi/chi/v5"

	"github.com/lechitz/AionApi/internal/adapters/primary/graph"
	"github.com/lechitz/AionApi/internal/adapters/primary/graph/directives"
	"github.com/lechitz/AionApi/internal/adapters/primary/graph/resolvers"

	generic "github.com/lechitz/AionApi/internal/adapters/primary/http/controllers/generic/handler"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/authmiddleware"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/recoverymiddleware"

	"github.com/lechitz/AionApi/internal/core/ports/input"
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/platform/config"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func NewGraphqlHandler(tokenSvc input.TokenService, categoryService input.CategoryService, logger output.ContextLogger, cfg *config.Config) (http.Handler, error) {

	router := chi.NewRouter()

	router.Use(authmiddleware.New(tokenSvc, logger).Auth)

	genericHandler := generic.New(logger, cfg.General)
	router.Use(recoverymiddleware.New(genericHandler))

	resolver := &resolvers.Resolver{
		CategoryService: categoryService,
		Logger:          logger,
	}

	es := graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver,
		Directives: graph.DirectiveRoot{
			Auth: directives.Auth(),
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
