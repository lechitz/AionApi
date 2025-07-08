// Package httputils provides utilities for working with HTTP cookies.
package httputils

import (
	"net/http"
	"time"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/constants"
)

// SetAuthCookie sets a secure HTTP-only authentication cookie with the given token and expiration configuration.
func SetAuthCookie(w http.ResponseWriter, token string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     constants.AuthToken,
		Value:    token,
		Path:     constants.Path,
		Domain:   constants.Domain,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   maxAge,
	})
}

// ClearAuthCookie invalidates the authentication cookie by setting its value to empty and expiration to a past timestamp.
func ClearAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     constants.AuthToken,
		Value:    "",
		Path:     constants.Path,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}
