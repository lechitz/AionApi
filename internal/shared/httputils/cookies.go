// Package httputils provides utilities for working with HTTP cookies.
package httputils

import (
	"net/http"
	"time"

	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/shared/commonkeys"
)

// SetAuthCookie sets a secure HTTP-only authentication cookie with the given token and expiration configuration.
func SetAuthCookie(w http.ResponseWriter, token string, cfg config.CookieConfig) {
	http.SetCookie(w, &http.Cookie{
		Name:     commonkeys.AuthToken,
		Value:    token,
		Path:     cfg.Path,
		Domain:   cfg.Domain,
		HttpOnly: true,
		Secure:   true,
		SameSite: mapSameSite(cfg.SameSite),
		MaxAge:   cfg.MaxAge,
	})
}

// ClearAuthCookie invalidates the authentication cookie by setting its value to empty and expiration to a past timestamp.
func ClearAuthCookie(w http.ResponseWriter, cfg config.CookieConfig) {
	http.SetCookie(w, &http.Cookie{
		Name:     commonkeys.AuthToken,
		Value:    "",
		Path:     cfg.Path,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

func mapSameSite(sameSite string) http.SameSite {
	switch sameSite {
	case "Strict":
		return http.SameSiteStrictMode
	case "Lax":
		return http.SameSiteLaxMode
	case "None":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteDefaultMode
	}
}
