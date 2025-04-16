package context

import (
	"context"
	"github.com/google/uuid"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/contextbuilder/constants"
	"net/http"
)

func InjectRequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New().String()

		ctx := context.WithValue(r.Context(), constants.RequestID, reqID)
		r = r.WithContext(ctx)

		w.Header().Set(constants.XRequestID, reqID)

		next.ServeHTTP(w, r)
	})
}
