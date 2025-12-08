// Package servicetoken provides S2S authentication middleware.
// See README.md for detailed documentation.
package servicetoken

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
)

const (
	HeaderServiceKey  = "X-Service-Key"
	HeaderServiceUser = "X-Service-User-Id"

	ErrServiceTokenInvalid = "service token invalid"

	LogNoServiceKey  = "servicetoken: no service key configured"
	LogInvalidKey    = "servicetoken: invalid service key"
	LogS2SAuth       = "servicetoken: S2S auth"
	LogInvalidUserID = "servicetoken: invalid user id"

	LogKeyUserID = "user_id"
	LogKeyPath   = "path"
	LogKeyValue  = "value"
)

// New returns a middleware that validates S2S authentication via X-Service-Key header.
func New(cfg *config.Config, log logger.ContextLogger) func(next http.Handler) http.Handler {
	serviceKey := ""
	if cfg != nil {
		serviceKey = strings.TrimSpace(cfg.AionChat.ServiceKey)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			providedKey := r.Header.Get(HeaderServiceKey)

			if providedKey == "" {
				next.ServeHTTP(w, r)
				return
			}

			if serviceKey == "" {
				log.Warnw(LogNoServiceKey)
				http.Error(w, ErrServiceTokenInvalid, http.StatusUnauthorized)
				return
			}

			if providedKey != serviceKey {
				log.Warnw(LogInvalidKey)
				http.Error(w, ErrServiceTokenInvalid, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ctxkeys.ServiceAccount, true)

			if userIDStr := r.Header.Get(HeaderServiceUser); userIDStr != "" {
				if userID, err := strconv.ParseUint(userIDStr, 10, 64); err == nil {
					ctx = context.WithValue(ctx, ctxkeys.UserID, userID)
					log.Infow(LogS2SAuth, LogKeyUserID, userID, LogKeyPath, r.URL.Path)
				} else {
					log.Warnw(LogInvalidUserID, LogKeyValue, userIDStr)
				}
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
