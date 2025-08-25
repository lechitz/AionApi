// Package httpserver registers routes.
package httpserver

import (
	handlerAuth "github.com/lechitz/AionApi/internal/adapters/primary/http/controllers/auth/handler"
	handlerGeneric "github.com/lechitz/AionApi/internal/adapters/primary/http/controllers/generic/handler"
	handlerUser "github.com/lechitz/AionApi/internal/adapters/primary/http/controllers/user/handler"
)

// AddHealthCheckRoutes registers a health check endpoint.
func (r *RouteComposer) AddHealthCheckRoutes(gh *handlerGeneric.Handler) func(HTTPRouter) {
	return func(healthGroup HTTPRouter) {
		healthGroup.Route("/health-check", func(healthProtected HTTPRouter) {
			healthProtected.Get("/", gh.HealthCheck)
		})
	}
}

// AddAuthRoutes registers authentication-related endpoints.
func (r *RouteComposer) AddAuthRoutes(ah *handlerAuth.Handler) func(HTTPRouter) {
	return func(rt HTTPRouter) {
		rt.Route("/auth", func(authGroup HTTPRouter) {
			authGroup.Post("/login", ah.Login)

			authGroup.Route("/", func(authProtected HTTPRouter) {
				authProtected.Use(r.authMiddleware.Auth)
				authProtected.Post("/logout", ah.Logout)
			})
		})
	}
}

// AddUserRoutes registers user-related endpoints (CRUD + account management).
func (r *RouteComposer) AddUserRoutes(uh *handlerUser.Handler) func(HTTPRouter) {
	return func(rt HTTPRouter) {
		rt.Route("/user", func(userGroup HTTPRouter) {
			userGroup.Post("/create", uh.Create)

			userGroup.Route("/", func(userProtected HTTPRouter) {
				userProtected.Use(r.authMiddleware.Auth)
				userProtected.Get("/all", uh.ListAll)
				userProtected.Get("/{user_id}", uh.GetUserByID)
				userProtected.Put("/", uh.UpdateUser)
				userProtected.Put("/password", uh.UpdateUserPassword)
				userProtected.Delete("/", uh.SoftDeleteUser)
			})
		})
	}
}
