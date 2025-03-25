package auth

import (
	"context"
	"github.com/lechitz/AionApi/core/domain"
	inputHttp "github.com/lechitz/AionApi/ports/input/http"
	outputHttp "github.com/lechitz/AionApi/ports/output/security"
	"net/http"

	"github.com/lechitz/AionApi/pkg/contextkeys"
	"go.uber.org/zap"
)

type MiddlewareAuth struct {
	AuthService  inputHttp.IAuthService
	TokenService outputHttp.ITokenService
	LoggerSugar  *zap.SugaredLogger
}

func NewAuthMiddleware(
	authService inputHttp.IAuthService,
	tokenService outputHttp.ITokenService,
	logger *zap.SugaredLogger,
) *MiddlewareAuth {
	return &MiddlewareAuth{
		AuthService:  authService,
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
			a.LoggerSugar.Warnw(ErrorUnauthorizedAccessMissingToken, contextkeys.Error, err.Error())
			http.Error(w, ErrorUnauthorizedAccessMissingToken, http.StatusUnauthorized)
			return
		}

		tokenDomain := domain.TokenDomain{
			Token: tokenCookie,
		}

		userID, token, err := a.TokenService.Check(*ctx, tokenDomain.Token)
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

func extractTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(contextkeys.AuthToken)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
