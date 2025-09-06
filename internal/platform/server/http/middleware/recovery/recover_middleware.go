// Package recoverymiddleware provides HTTP middleware that recovers from panics.
package recoverymiddleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/lechitz/AionApi/internal/adapter/server/http/generic/handler"
)

// New is a middleware that recovers from panics, logs the error, and returns an internal server error response.
func New(recoveryHandler *handler.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					errorID := uuid.New().String()

					recoveryHandler.RecoveryHandler(w, r, rec, errorID)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
