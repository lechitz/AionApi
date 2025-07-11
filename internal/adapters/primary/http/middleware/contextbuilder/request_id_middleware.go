// Package contextbuilder provides middleware for injecting request IDs into HTTP request contexts.
package contextbuilder

import (
	"context"
	"net/http"

	"github.com/lechitz/AionApi/internal/shared/ctxkeys"

	"github.com/lechitz/AionApi/internal/shared/commonkeys"

	"github.com/google/uuid"
)

// InjectRequestIDMiddleware injects a request ID into the HTTP request context and sets the X-Request-ID header.
func InjectRequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get(commonkeys.XRequestID)
		if reqID == "" {
			reqID = uuid.New().String()
		}

		ctx := context.WithValue(r.Context(), ctxkeys.CtxKeyRequestID, reqID)

		r = r.WithContext(ctx)

		w.Header().Set(commonkeys.XRequestID, reqID)

		next.ServeHTTP(w, r)
	})
}
