package handlers

import (
	"encoding/json"
	"github.com/lechitz/AionApi/core/domain"
	inputHttp "github.com/lechitz/AionApi/ports/input/http"
	"net/http"
	"time"

	"github.com/lechitz/AionApi/adapters/input/http/dto"
	msg "github.com/lechitz/AionApi/adapters/input/http/handlers/messages"
	"github.com/lechitz/AionApi/pkg/contextkeys"
	"github.com/lechitz/AionApi/pkg/utils"
	"go.uber.org/zap"
)

type Auth struct {
	AuthService inputHttp.IAuthService
	LoggerSugar *zap.SugaredLogger
}

func (a *Auth) LoginHandler(w http.ResponseWriter, r *http.Request) {

	contextControl := domain.ContextControl{
		BaseContext: r.Context(),
	}

	var loginUserRequest dto.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&loginUserRequest); err != nil {
		utils.HandleError(w, a.LoggerSugar, http.StatusBadRequest, msg.ErrorToDecodeLoginRequest, err)
		a.LoggerSugar.Errorw(msg.ErrorToDecodeLoginRequest, contextkeys.Error, err.Error())
		return
	}

	userDomain := domain.UserDomain{
		Username: loginUserRequest.Username,
		Password: loginUserRequest.Password,
	}

	userDB, token, err := a.AuthService.Login(contextControl, userDomain, loginUserRequest.Password)
	if err != nil {
		utils.HandleError(w, a.LoggerSugar, http.StatusInternalServerError, msg.ErrorToLogin, err)
		a.LoggerSugar.Errorw(msg.ErrorToLogin, contextkeys.Error, err.Error())
		return
	}

	setAuthCookie(w, token, 0)

	loginUserResponse := dto.LoginUserResponse{
		Username: userDB.Username,
	}

	response := utils.ObjectResponse(loginUserResponse, msg.SuccessToLogin)
	a.LoggerSugar.Infow(msg.SuccessToLogin, contextkeys.Username, loginUserResponse.Username)
	utils.ResponseReturn(w, http.StatusOK, response.Bytes())
}

func (a *Auth) LogoutHandler(w http.ResponseWriter, r *http.Request) {

	contextControl := domain.ContextControl{
		BaseContext: r.Context(),
	}

	tokenCookie, err := r.Cookie(contextkeys.AuthToken)
	if err != nil {
		utils.HandleError(w, a.LoggerSugar, http.StatusUnauthorized, msg.ErrorToRetrieveToken, err)
		a.LoggerSugar.Errorw(msg.ErrorToRetrieveToken, contextkeys.Error, err.Error())
		return
	}

	tokenValue := tokenCookie.Value

	userID, ok := contextControl.BaseContext.Value(contextkeys.UserID).(uint64)
	if !ok {
		utils.HandleError(w, a.LoggerSugar, http.StatusInternalServerError, msg.ErrorToRetrieveUserID, err)
		a.LoggerSugar.Errorw(msg.ErrorToRetrieveUserID, contextkeys.Error, msg.ErrorUserIDIsNil)
		return
	}

	if err := a.AuthService.Logout(contextControl, tokenValue); err != nil {
		utils.HandleError(w, a.LoggerSugar, http.StatusInternalServerError, msg.ErrorToLogout, err)
		a.LoggerSugar.Errorw(msg.ErrorToLogout, contextkeys.Error, err.Error())
		return
	}

	clearAuthCookie(w)

	response := utils.ObjectResponse(nil, msg.SuccessToLogout)
	a.LoggerSugar.Infow(msg.SuccessToLogout, contextkeys.UserID, userID)
	utils.ResponseReturn(w, http.StatusOK, response.Bytes())
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
