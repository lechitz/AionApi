// Package httpserver provides functionality for configuring and managing HTTP routers.
package httpserver

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/auth"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/generic"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/user"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/authmiddleware"
	"github.com/lechitz/AionApi/internal/core/ports/input"
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
)

// BuildRouterRoutes sets up API routes, integrates middlewares.
func BuildRouterRoutes(
	logger output.ContextLogger,
	userService input.UserService,
	authService input.AuthService,
	tokenRepository output.TokenStore,
	contextPath string,
	adapter input.HTTPRouter,
	tokenClaimsExtractor output.TokenClaimsExtractor,
	cfg *config.Config,
	genericHandler *generic.Handler,
) (input.HTTPRouter, error) {
	userHandler := user.New(userService, cfg, logger)
	authHandler := auth.New(authService, cfg, logger)

	newAuthMiddleware := authmiddleware.New(tokenRepository, logger, tokenClaimsExtractor)

	r := &RouteComposer{
		BasePath:       contextPath,
		Router:         adapter,
		logger:         logger,
		authMiddleware: newAuthMiddleware,
	}

	adapter.SetNotFoundHandler(genericHandler.NotFoundHandler)
	adapter.SetMethodNotAllowedHandler(genericHandler.MethodNotAllowedHandler)
	adapter.SetErrorHandler(genericHandler.ErrorHandler)

	adapter.Route(contextPath, func(rt input.HTTPRouter) {
		rt.Group(r.AddHealthCheckRoutes(genericHandler))
		rt.Group(r.AddUserRoutes(userHandler))
		rt.Group(r.AddAuthRoutes(authHandler))
	})

	return adapter, nil
}

// injectRequestIDMiddleware injects a request ID into the HTTP request context and sets the X-Request-ID header.
func injectRequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get(commonkeys.XRequestID)
		if reqID == "" {
			reqID = uuid.New().String()
		}

		ctx := context.WithValue(r.Context(), ctxkeys.RequestID, reqID)

		r = r.WithContext(ctx)

		w.Header().Set(commonkeys.XRequestID, reqID)

		next.ServeHTTP(w, r)
	})
}
