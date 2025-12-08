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
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
)

const (
	// HeaderServiceKey is the HTTP header name used to pass the service API key.
	HeaderServiceKey = "X-Service-Key" // #nosec G101: false positive — header name, not a credential

	// HeaderServiceUser is the HTTP header name used to pass an optional service user id.
	HeaderServiceUser = "X-Service-User-Id"

	// ErrServiceTokenInvalid is the error message returned when the service token is invalid.
	ErrServiceTokenInvalid = "service token invalid"

	// LogNoServiceKey is the log message when no service key is configured.
	LogNoServiceKey = "servicetoken: no service key configured"

	// LogInvalidKey is the log message when an invalid service key is provided.
	LogInvalidKey = "servicetoken: invalid service key"

	// LogS2SAuth is the log message when S2S authentication is successful.
	LogS2SAuth = "servicetoken: S2S auth"

	// LogInvalidUserID is the log message when an invalid user ID is provided.
	LogInvalidUserID = "servicetoken: invalid user id"

	// LogKeyValue is the log key for generic value.
	LogKeyValue = "value"
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
					log.Infow(LogS2SAuth, commonkeys.UserID, userID, commonkeys.URLPath, r.URL.Path)
				} else {
					log.Warnw(LogInvalidUserID, LogKeyValue, userIDStr)
				}
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
