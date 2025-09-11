// package handler
package handler

import (
	"net/http"

	"github.com/lechitz/AionApi/internal/auth/adapter/primary/http/middleware"
	"github.com/lechitz/AionApi/internal/platform/server/http/ports"
)

// RegisterHTTP registra as rotas REST do contexto de autenticação.
func RegisterHTTP(r ports.Router, h *Handler) {
	r.Group("/v1/auth", func(ar ports.Router) {
		// Público
		ar.POST("/login", http.HandlerFunc(h.Login))

		// Protegido por Auth
		authmw := middleware.New(h.Service, h.Logger) // h já tem deps
		ar.GroupWith(authmw.Auth, func(pr ports.Router) {
			pr.POST("/logout", http.HandlerFunc(h.Logout))
		})
	})
}
