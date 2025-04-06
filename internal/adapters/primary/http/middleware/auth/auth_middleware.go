package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lechitz/AionApi/internal/platform/config"
	"net/http"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/auth/constants"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/ports/output/cache"
	"go.uber.org/zap"
)

type MiddlewareAuth struct {
	TokenService cache.TokenRepositoryPort
	LoggerSugar  *zap.SugaredLogger
}

func NewAuthMiddleware(tokenService cache.TokenRepositoryPort, logger *zap.SugaredLogger) *MiddlewareAuth {
	return &MiddlewareAuth{
		TokenService: tokenService,
		LoggerSugar:  logger,
	}
}

func (a *MiddlewareAuth) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := &domain.ContextControl{
			BaseContext:     r.Context(),
			CancelCauseFunc: nil,
		}

		tokenCookie, err := extractTokenFromCookie(r)
		if err != nil {
			a.LoggerSugar.Warnw(constants.ErrorUnauthorizedAccessMissingToken, constants.Error, err.Error())
			http.Error(w, constants.ErrorUnauthorizedAccessMissingToken, http.StatusUnauthorized)
			return
		}

		parsedToken, err := jwt.Parse(tokenCookie, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Setting.SecretKey), nil
		})
		if err != nil || parsedToken == nil || !parsedToken.Valid {
			a.LoggerSugar.Warnw(constants.ErrorUnauthorizedAccessInvalidToken, constants.Error, err)
			http.Error(w, constants.ErrorUnauthorizedAccessInvalidToken, http.StatusUnauthorized)
			return
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			a.LoggerSugar.Warnw(constants.ErrorUnauthorizedAccessInvalidToken)
			http.Error(w, constants.ErrorUnauthorizedAccessInvalidToken, http.StatusUnauthorized)
			return
		}

		userIDFloat, ok := claims[constants.UserID].(float64)
		if !ok {
			a.LoggerSugar.Warnw(constants.ErrorUnauthorizedAccessInvalidToken)
			http.Error(w, constants.ErrorUnauthorizedAccessInvalidToken, http.StatusUnauthorized)
			return
		}
		userID := uint64(userIDFloat)

		tokenDomain := domain.TokenDomain{UserID: userID, Token: tokenCookie}
		_, err = a.TokenService.Get(*ctx, tokenDomain)
		if err != nil {
			a.LoggerSugar.Warnw(constants.ErrorUnauthorizedAccessInvalidToken, constants.Error, err.Error())
			http.Error(w, constants.ErrorUnauthorizedAccessInvalidToken, http.StatusUnauthorized)
			return
		}

		newCtx := context.WithValue(r.Context(), constants.UserID, tokenDomain.UserID)
		newCtx = context.WithValue(newCtx, constants.Token, tokenCookie)

		a.LoggerSugar.Infow(constants.SuccessTokenValidated, constants.UserID, tokenDomain.UserID)
		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}

func extractTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(constants.AuthToken)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
