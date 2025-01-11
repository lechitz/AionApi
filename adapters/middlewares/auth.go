package middlewares

import (
	"context"
	"github.com/lechitz/AionApi/ports/output"
	"go.uber.org/zap"
	"net/http"
)

type AuthMiddleware struct {
	TokenStore output.IAuthRepository
	Logger     *zap.SugaredLogger
}

func NewAuthMiddleware(tokenStore output.IAuthRepository, logger *zap.SugaredLogger) *AuthMiddleware {
	return &AuthMiddleware{
		TokenStore: tokenStore,
		Logger:     logger,
	}
}

func (a *AuthMiddleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenCookie, err := extractTokenFromCookie(r)
		if err != nil {
			a.Logger.Warnw("Unauthorized access: missing token", "error", err.Error())
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, userIDFloat, err := a.TokenStore.ValidateToken(r.Context(), tokenCookie)
		if err != nil {
			a.Logger.Warnw("Unauthorized access: invalid token", "error", err.Error())
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		a.Logger.Infow("Token validated successfully, adding userID to context", "userID", userIDFloat)

		ctx := context.WithValue(r.Context(), "id", userIDFloat)
		ctx = context.WithValue(ctx, "auth_token", token)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
