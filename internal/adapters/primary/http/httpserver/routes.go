// Package httpserver provides functionality for configuring and managing HTTP routers.
package httpserver

import (
	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/auth"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/generic"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/user"
	"github.com/lechitz/AionApi/internal/core/ports/input"
)

// AddHealthCheckRoutes registers a health check endpoint.
// Used for readiness/liveness probes or basic service health verification.
func (r *RouteComposer) AddHealthCheckRoutes(gh *generic.Handler) func(input.HTTPRouter) {
	return func(healthGroup input.HTTPRouter) {
		healthGroup.Route("/health-check", func(healthProtected input.HTTPRouter) {
			healthProtected.Get("/", gh.HealthCheck)
		})
	}
}

// AddAuthRoutes registers authentication-related HTTP endpoints.
// This uses a subgroup to ensure middleware is attached BEFORE the protected route.
func (r *RouteComposer) AddAuthRoutes(ah *auth.Handler) func(input.HTTPRouter) {
	return func(rt input.HTTPRouter) {
		rt.Route("/auth", func(authGroup input.HTTPRouter) {
			authGroup.Post("/login", ah.Login)

			authGroup.Route("/", func(authProtected input.HTTPRouter) {
				authProtected.Use(r.authMiddleware.Auth)
				authProtected.Post("/logout", ah.Logout)
			})
		})
	}
}

// AddUserRoutes registers all user-related HTTP routes (CRUD and account management).
// The inner group is protected with authentication middleware.
func (r *RouteComposer) AddUserRoutes(uh *user.Handler) func(input.HTTPRouter) {
	return func(rt input.HTTPRouter) {
		rt.Route("/user", func(userGroup input.HTTPRouter) {
			userGroup.Post("/create", uh.Create)

			userGroup.Route("/", func(userProtected input.HTTPRouter) {
				userProtected.Use(r.authMiddleware.Auth)
				userProtected.Get("/all", uh.GetAllUsers)
				userProtected.Get("/{user_id}", uh.GetUserByID)
				userProtected.Put("/", uh.UpdateUser)
				userProtected.Put("/password", uh.UpdateUserPassword)
				userProtected.Delete("/", uh.SoftDeleteUser)
			})
		})
	}
}
