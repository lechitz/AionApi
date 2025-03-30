package handlers

import (
	"encoding/json"
	constants "github.com/lechitz/AionApi/internal/adapters/primary/http/constants"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/dto"
	"github.com/lechitz/AionApi/internal/core/domain"
	inputHttp "github.com/lechitz/AionApi/internal/core/ports/input/http"
	"net/http"
	"time"

	"github.com/lechitz/AionApi/pkg/utils"
	"go.uber.org/zap"
)

type Auth struct {
	AuthService inputHttp.IAuthService
	LoggerSugar *zap.SugaredLogger
}

func (a *Auth) LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctxControl := domain.ContextControl{BaseContext: r.Context()}

	var loginReq dto.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		a.logAndRespondError(w, http.StatusBadRequest, constants.ErrorToDecodeLoginRequest, err)
		return
	}

	userDomain := domain.UserDomain{Username: loginReq.Username}

	userDB, token, err := a.AuthService.Login(ctxControl, userDomain, loginReq.Password)
	if err != nil {
		a.logAndRespondError(w, http.StatusInternalServerError, constants.ErrorToLogin, err)
		return
	}

	setAuthCookie(w, token, 0)

	response := dto.LoginUserResponse{Username: userDB.Username}
	a.LoggerSugar.Infow(constants.SuccessToLogin, constants.Username, userDB.Username)

	utils.ResponseReturn(w, http.StatusOK, utils.ObjectResponse(response, constants.SuccessToLogin).Bytes())
}

func (a *Auth) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctxControl := domain.ContextControl{BaseContext: r.Context()}

	cookie, err := r.Cookie(constants.AuthToken)
	if err != nil {
		a.logAndRespondError(w, http.StatusUnauthorized, constants.ErrorToRetrieveToken, err)
		return
	}

	token := cookie.Value

	userID, ok := ctxControl.BaseContext.Value(constants.UserID).(uint64)
	if !ok {
		a.logAndRespondError(w, http.StatusInternalServerError, constants.ErrorToRetrieveUserID, nil)
		return
	}

	if err := a.AuthService.Logout(ctxControl, token); err != nil {
		a.logAndRespondError(w, http.StatusInternalServerError, constants.ErrorToLogout, err)
		return
	}

	clearAuthCookie(w)

	a.LoggerSugar.Infow(constants.SuccessToLogout, constants.UserID, userID)
	utils.ResponseReturn(w, http.StatusOK, utils.ObjectResponse(nil, constants.SuccessToLogout).Bytes())
}

func setAuthCookie(w http.ResponseWriter, token string, maxAge int) {
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

func clearAuthCookie(w http.ResponseWriter) {
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

func (a *Auth) logAndRespondError(w http.ResponseWriter, status int, message string, err error) {
	if err != nil {
		a.LoggerSugar.Errorw(message, constants.Error, err.Error())
	} else {
		a.LoggerSugar.Errorw(message)
	}
	utils.HandleError(w, a.LoggerSugar, status, message, err)
}
