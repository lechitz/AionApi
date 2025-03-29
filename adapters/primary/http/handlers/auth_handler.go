package handlers

import (
	"encoding/json"
	inputHttp "github.com/lechitz/AionApi/core/ports/input/http"
	"net/http"
	"time"

	"github.com/lechitz/AionApi/adapters/primary/http/dto"
	msg "github.com/lechitz/AionApi/adapters/primary/http/handlers/messages"
	"github.com/lechitz/AionApi/core/domain"
	"github.com/lechitz/AionApi/pkg/contextkeys"
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
		a.logAndRespondError(w, http.StatusBadRequest, msg.ErrorToDecodeLoginRequest, err)
		return
	}

	userDomain := domain.UserDomain{Username: loginReq.Username}

	userDB, token, err := a.AuthService.Login(ctxControl, userDomain, loginReq.Password)
	if err != nil {
		a.logAndRespondError(w, http.StatusInternalServerError, msg.ErrorToLogin, err)
		return
	}

	setAuthCookie(w, token, 0)

	response := dto.LoginUserResponse{Username: userDB.Username}
	a.LoggerSugar.Infow(msg.SuccessToLogin, contextkeys.Username, userDB.Username)

	utils.ResponseReturn(w, http.StatusOK, utils.ObjectResponse(response, msg.SuccessToLogin).Bytes())
}

func (a *Auth) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctxControl := domain.ContextControl{BaseContext: r.Context()}

	cookie, err := r.Cookie(contextkeys.AuthToken)
	if err != nil {
		a.logAndRespondError(w, http.StatusUnauthorized, msg.ErrorToRetrieveToken, err)
		return
	}

	token := cookie.Value

	userID, ok := ctxControl.BaseContext.Value(contextkeys.UserID).(uint64)
	if !ok {
		a.logAndRespondError(w, http.StatusInternalServerError, msg.ErrorToRetrieveUserID, nil)
		return
	}

	if err := a.AuthService.Logout(ctxControl, token); err != nil {
		a.logAndRespondError(w, http.StatusInternalServerError, msg.ErrorToLogout, err)
		return
	}

	clearAuthCookie(w)

	a.LoggerSugar.Infow(msg.SuccessToLogout, contextkeys.UserID, userID)
	utils.ResponseReturn(w, http.StatusOK, utils.ObjectResponse(nil, msg.SuccessToLogout).Bytes())
}

func setAuthCookie(w http.ResponseWriter, token string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     contextkeys.AuthToken,
		Value:    token,
		Path:     contextkeys.Path,
		Domain:   contextkeys.Domain,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   maxAge,
	})
}

func clearAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     contextkeys.AuthToken,
		Value:    "",
		Path:     contextkeys.Path,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

func (a *Auth) logAndRespondError(w http.ResponseWriter, status int, message string, err error) {
	if err != nil {
		a.LoggerSugar.Errorw(message, contextkeys.Error, err.Error())
	} else {
		a.LoggerSugar.Errorw(message)
	}
	utils.HandleError(w, a.LoggerSugar, status, message, err)
}
