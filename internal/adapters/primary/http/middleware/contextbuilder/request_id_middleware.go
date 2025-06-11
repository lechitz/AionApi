// Package contextbuilder provides middleware for injecting request IDs into HTTP request contexts.
package contextbuilder

import (
	"context"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/contextbuilder/constants"
	"net/http"

	"github.com/google/uuid"
)

// ctxKeyRequestID is a type for storing request IDs in HTTP request contexts.
type ctxKeyRequestID struct{}

// InjectRequestIDMiddleware injects a request ID into the HTTP request context and sets the X-Request-ID header.
// It is used to track requests across multiple services and to correlate logs and metrics.
// It is recommended to use this middleware as early as possible in the middleware chain.
func InjectRequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New().String()
		ctx := context.WithValue(r.Context(), ctxKeyRequestID{}, reqID)
		r = r.WithContext(ctx)
		w.Header().Set(constants.XRequestID, reqID)
		next.ServeHTTP(w, r)
	})
}
