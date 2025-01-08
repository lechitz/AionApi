package handlers

import (
	"context"
	"encoding/json"
	"github.com/jinzhu/copier"
	"github.com/lechitz/AionApi/adapters/input/constants"
	"github.com/lechitz/AionApi/adapters/input/http/dto"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/pkg/utils"
	"github.com/lechitz/AionApi/ports/input"
	"go.uber.org/zap"
	"net/http"
)

type Auth struct {
	AuthService input.IAuthService
	LoggerSugar *zap.SugaredLogger
}

func (u *Auth) LoginHandler(w http.ResponseWriter, r *http.Request) {

	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	var loginUserRequest dto.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&loginUserRequest); err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusInternalServerError, constants.ErrorToDecodeLoginRequest, err)
		return
	}

	var userDomain domain.UserDomain
	copier.Copy(&userDomain, &loginUserRequest)

	userDB, token, err := u.AuthService.Login(contextControl, userDomain)
	if err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusInternalServerError, constants.ErrorToLogin, err)
		return
	}

	//If you want to use the function: auth.extractTokenFromBearer
	// you don't need to set the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true, // It's not possible to access the cookie via JavaScript
		Secure:   true, // Only send the cookie if the request is being sent over HTTPS
		SameSite: http.SameSiteStrictMode,
	})

	var loginUserResponse dto.LoginUserResponse
	copier.Copy(&loginUserResponse, &userDB)

	response := utils.ObjectResponse(loginUserResponse, constants.SuccessToLogin)
	utils.ResponseReturn(w, http.StatusOK, response.Bytes())
}
