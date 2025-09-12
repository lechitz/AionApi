package handler

import (
	"net/http"

	authMiddleware "github.com/lechitz/AionApi/internal/auth/adapter/primary/http/middleware"
	authinput "github.com/lechitz/AionApi/internal/auth/core/ports/input"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/internal/platform/server/http/ports"
)

// RegisterHTTP registra rotas públicas e protegidas sob um único prefixo.
// Ex.: /aion-api/v1/users/...
func RegisterHTTP(r ports.Router, h *Handler, authService authinput.AuthService, lg logger.ContextLogger) {
	r.Group("/user", func(ur ports.Router) {
		// --- Públicas ---
		ur.POST("/create", http.HandlerFunc(h.Create)) // se seu produto permite cadastro aberto

		// --- Protegidas (somente se houver AuthService) ---
		if authService != nil {
			mw := authMiddleware.New(authService, lg)
			ur.GroupWith(mw.Auth, func(pr ports.Router) {
				// Leitura protegida
				pr.GET("/all", http.HandlerFunc(h.ListAll))           // tipicamente admin
				pr.GET("/{user_id}", http.HandlerFunc(h.GetUserByID)) // self ou admin

				// Escritas protegidas
				pr.PUT("/", http.HandlerFunc(h.UpdateUser))
				pr.PUT("/password", http.HandlerFunc(h.UpdateUserPassword))
				pr.DELETE("/", http.HandlerFunc(h.SoftDeleteUser))
			})
		}
	})
}
