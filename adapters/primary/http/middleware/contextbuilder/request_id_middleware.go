// Package context provides middleware for injecting request IDs into HTTP request contexts.
package context

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/lechitz/AionApi/adapters/primary/http/middleware/contextbuilder/constants"
)

// InjectRequestIDMiddleware adds a unique request ID to each HTTP request context and response header for traceability.
func InjectRequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New().String()

		ctx := context.WithValue(r.Context(), constants.RequestID, reqID)
		r = r.WithContext(ctx)

		w.Header().Set(constants.XRequestID, reqID)

		next.ServeHTTP(w, r)
	})
}
