package http

import (
	"net/http"

	"github.com/lechitz/AionApi/internal/platform/server/http/ports"
)

// internal/user/adapter/primary/http/register.go
func RegisterHTTP(r ports.Router, h *Handler, d Deps) {
	r.POST("/users/register", http.HandlerFunc(h.Register)) // público
	if d.Auth != nil {
		r.Group("", func(pr ports.Router) {
			pr.Use(d.Auth)
			pr.GET("/users", http.HandlerFunc(h.List))
			pr.PUT("/users/{id}", http.HandlerFunc(h.Update))
		})
	}
}
