// Package httpserver provides functionality for configuring and managing HTTP routers.
package httpserver

import (
	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/auth"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/generic"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/user"
	"github.com/lechitz/AionApi/internal/core/ports/output"
)

// AddHealthCheckRoutes registers a health check endpoint.
// Used for readiness/liveness probes or basic service health verification.
func (r *RouteComposer) AddHealthCheckRoutes(gh *generic.Handler) func(output.Router) {
	return func(healthGroup output.Router) {
		healthGroup.Route("/health-check", func(healthProtected output.Router) {
			healthProtected.Get("/", gh.HealthCheck)
		})
	}
}

// AddAuthRoutes registers authentication-related HTTP endpoints.
// This uses a subgroup to ensure middleware is attached BEFORE the protected route.
func (r *RouteComposer) AddAuthRoutes(ah *auth.Handler) func(output.Router) {
	return func(rt output.Router) {
		rt.Route("/auth", func(authGroup output.Router) {
			authGroup.Post("/login", ah.Login)

			authGroup.Route("/", func(authProtected output.Router) {
				authProtected.Use(r.authMiddleware.Auth)
				authProtected.Post("/logout", ah.Logout)
			})
		})
	}
}

// AddUserRoutes registers all user-related HTTP routes (CRUD and account management).
// The inner group is protected with authentication middleware.
func (r *RouteComposer) AddUserRoutes(uh *user.Handler) func(output.Router) {
	return func(rt output.Router) {
		rt.Route("/user", func(userGroup output.Router) {
			userGroup.Post("/create", uh.Create)

			userGroup.Route("/", func(userProtected output.Router) {
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
