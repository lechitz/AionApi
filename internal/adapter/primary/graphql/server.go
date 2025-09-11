package graphql

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/go-chi/chi/v5"
	authdir "github.com/lechitz/AionApi/internal/auth/adapter/primary/graphql/directives"
	authmw "github.com/lechitz/AionApi/internal/auth/adapter/primary/http/middleware"
	authInput "github.com/lechitz/AionApi/internal/auth/core/ports/input"
	categoryInput "github.com/lechitz/AionApi/internal/category/core/ports/input"
	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	generic "github.com/lechitz/AionApi/internal/platform/server/http/generic/handler"
	"github.com/lechitz/AionApi/internal/platform/server/http/middleware/recovery"
)

func NewGraphqlHandler(
	authService authInput.AuthService,
	categoryService categoryInput.CategoryService,
	log logger.ContextLogger,
	cfg *config.Config,
) (http.Handler, error) {
	r := chi.NewRouter()

	// Auth opcional (quando existir)
	if authService != nil {
		r.Use(authmw.New(authService, log).Auth)
	}

	// Recovery específico do subrouter GraphQL
	r.Use(recovery.New(generic.New(log, cfg.General)))

	// Usa o Resolver local (mesmo pacote graphql)
	resolvers := &Resolver{
		CategoryService: categoryService,
		Logger:          log,
	}

	es := NewExecutableSchema(Config{
		Resolvers: resolvers,
		Directives: DirectiveRoot{
			Auth: authdir.Auth(),
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

	// IMPORTANTE: monta em "/" — o composer montará o subrouter em cfg.ServerGraphql.Path
	r.Handle("/", srv)
	return r, nil
}
