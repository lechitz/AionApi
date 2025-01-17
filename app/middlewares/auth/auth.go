package auth

import (
	"context"
	"net/http"

	"github.com/lechitz/AionApi/core/domain/entities"
	"github.com/lechitz/AionApi/core/service"
	"github.com/lechitz/AionApi/pkg/contextkeys"
	"go.uber.org/zap"
)

type MiddlewareAuth struct {
	AuthService  *service.AuthService
	TokenService *service.TokenService
	LoggerSugar  *zap.SugaredLogger
}

func NewAuthMiddleware(authService *service.AuthService, tokenService *service.TokenService, loggerSugar *zap.SugaredLogger) *MiddlewareAuth {
	return &MiddlewareAuth{
		AuthService:  authService,
		TokenService: tokenService,
		LoggerSugar:  loggerSugar,
	}
}

func (a *MiddlewareAuth) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := &entities.ContextControl{
			BaseContext:     r.Context(),
			CancelCauseFunc: nil,
		}

		tokenCookie, err := extractTokenFromCookie(r)
		if err != nil {
			a.LoggerSugar.Warnw(ErrorUnauthorizedAccessMissingToken, contextkeys.Error, err.Error())
			http.Error(w, ErrorUnauthorizedAccessMissingToken, http.StatusUnauthorized)
			return
		}

		tokenDomain := entities.TokenDomain{
			Token: tokenCookie,
		}

		userID, token, err := a.TokenService.CheckToken(*ctx, tokenDomain.Token)
		if err != nil {
			a.LoggerSugar.Warnw(ErrorUnauthorizedAccessInvalidToken, contextkeys.Error, err.Error())
			http.Error(w, ErrorUnauthorizedAccessInvalidToken, http.StatusUnauthorized)
			return
		}

		if r.Context().Value(contextkeys.UserID) == nil {
			newCtx := context.WithValue(r.Context(), contextkeys.UserID, userID)
			r = r.WithContext(newCtx)
		}

		if r.Context().Value(contextkeys.Token) == nil {
			newCtx := context.WithValue(r.Context(), contextkeys.Token, token)
			r = r.WithContext(newCtx)
		}

		a.LoggerSugar.Infow(SuccessTokenValidated, contextkeys.UserID, userID)

		next.ServeHTTP(w, r)
	})
}

func (a *MiddlewareAuth) AuthRevoke(next http.Handler) http.Handler {

	//TODO: Implement AuthRevoke middleware

	return nil
}

func extractTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(contextkeys.AuthToken)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
