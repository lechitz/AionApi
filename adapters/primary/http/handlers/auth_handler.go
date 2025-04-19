package handlers

import (
	"encoding/json"
	constants2 "github.com/lechitz/AionApi/adapters/primary/http/constants"
	"github.com/lechitz/AionApi/adapters/primary/http/dto"
	"github.com/lechitz/AionApi/adapters/primary/http/middleware/response"
	"net/http"
	"time"

	"github.com/lechitz/AionApi/internal/core/domain"
	inputHttp "github.com/lechitz/AionApi/internal/core/ports/input/http"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
)

type Auth struct {
	AuthService inputHttp.AuthService
	Logger      logger.Logger
}

func NewAuth(authService inputHttp.AuthService, logger logger.Logger) *Auth {
	return &Auth{
		AuthService: authService,
		Logger:      logger,
	}
}

func (a *Auth) LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var loginReq dto.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		a.logAndRespondError(w, http.StatusBadRequest, constants2.ErrorToDecodeLoginRequest, err)
		return
	}

	userDomain := domain.UserDomain{Username: loginReq.Username}

	userDB, token, err := a.AuthService.Login(ctx, userDomain, loginReq.Password)
	if err != nil {
		a.logAndRespondError(w, http.StatusInternalServerError, constants2.ErrorToLogin, err)
		return
	}

	setAuthCookie(w, token, 0)

	response := dto.LoginUserResponse{Username: userDB.Username}

	response.ResponseReturn(w, http.StatusOK, response.ObjectResponse(response, constants2.SuccessToLogin).Bytes())
}

func (a *Auth) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value(constants2.UserID).(uint64)
	if !ok || userID == 0 {
		a.logAndRespondError(w, http.StatusUnauthorized, constants2.ErrorToRetrieveUserID, nil)
		return
	}

	tokenString, ok := ctx.Value(constants2.Token).(string)
	if !ok || tokenString == "" {
		a.logAndRespondError(w, http.StatusUnauthorized, constants2.ErrorToRetrieveToken, nil)
		return
	}

	if err := a.AuthService.Logout(ctx, tokenString); err != nil {
		a.logAndRespondError(w, http.StatusInternalServerError, constants2.ErrorToLogout, err)
		return
	}

	clearAuthCookie(w)

	tokenPreview := ""
	if len(tokenString) >= 10 {
		tokenPreview = tokenString[:10] + "..."
	}

	a.Logger.Infow(
		constants2.SuccessToLogout,
		constants2.UserID, userID,
		constants2.Token, tokenPreview,
	)

	response.ResponseReturn(w, http.StatusOK, response.ObjectResponse(nil, constants2.SuccessToLogout).Bytes())
}

func (a *Auth) logAndRespondError(w http.ResponseWriter, status int, message string, err error) {
	if err != nil {
		a.Logger.Errorw(message, constants2.Error, err.Error())
	} else {
		a.Logger.Errorw(message)
	}
	response.HandleError(w, a.Logger, status, message, err)
}

func setAuthCookie(w http.ResponseWriter, token string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     constants2.AuthToken,
		Value:    token,
		Path:     constants2.Path,
		Domain:   constants2.Domain,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   maxAge,
	})
}

func clearAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     constants2.AuthToken,
		Value:    "",
		Path:     constants2.Path,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}
