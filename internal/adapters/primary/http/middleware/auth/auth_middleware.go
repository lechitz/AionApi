package auth

import (
	"context"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/auth/constants"
	"github.com/lechitz/AionApi/internal/core/domain"
	inputHttp "github.com/lechitz/AionApi/internal/core/ports/input/http"
	tokenports "github.com/lechitz/AionApi/internal/core/ports/output/cache"
	"net/http"

	"go.uber.org/zap"
)

type MiddlewareAuth struct {
	AuthService  inputHttp.IAuthService
	TokenService tokenports.TokenRepository
	LoggerSugar  *zap.SugaredLogger
}

func NewAuthMiddleware(
	authService inputHttp.IAuthService,
	tokenService tokenports.TokenRepository,
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
			a.LoggerSugar.Warnw(constants.ErrorUnauthorizedAccessMissingToken, constants.Error, err.Error())
			http.Error(w, constants.ErrorUnauthorizedAccessMissingToken, http.StatusUnauthorized)
			return
		}

		tokenDomain := domain.TokenDomain{
			Token: tokenCookie,
		}

		userID, token, err := a.TokenService.Check(*ctx, tokenDomain.Token)
		if err != nil {
			a.LoggerSugar.Warnw(constants.ErrorUnauthorizedAccessInvalidToken, constants.Error, err.Error())
			http.Error(w, constants.ErrorUnauthorizedAccessInvalidToken, http.StatusUnauthorized)
			return
		}

		if r.Context().Value(constants.UserID) == nil {
			newCtx := context.WithValue(r.Context(), constants.UserID, userID)
			r = r.WithContext(newCtx)
		}

		if r.Context().Value(constants.Token) == nil {
			newCtx := context.WithValue(r.Context(), constants.Token, token)
			r = r.WithContext(newCtx)
		}

		a.LoggerSugar.Infow(constants.SuccessTokenValidated, constants.UserID, userID)

		next.ServeHTTP(w, r)
	})
}

func extractTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(constants.AuthToken)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
