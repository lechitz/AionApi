package handlers

import (
	"context"
	"encoding/json"
	"github.com/jinzhu/copier"
	"github.com/lechitz/AionApi/adapters/middlewares"
	"github.com/lechitz/AionApi/adapters/security"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/pkg/utils"
	"net/http"
)

const (
	ErrorToDecodeLoginRequest = "error to decode login request"
	ErrorToGetUserByUsername  = "error to get user by username"
	ErrorToVerifyPassword     = "error to verify password"
	ErrorToCreateToken        = "error to create token"

	SuccessToLogin = "success to login"
)

func (lg *Login) Login(w http.ResponseWriter, r *http.Request) {

	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	var loginRequest LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		utils.HandleError(w, lg.LoggerSugar, http.StatusInternalServerError, ErrorToDecodeLoginRequest, err)
		return
	}

	var loginDomain domain.LoginDomain
	copier.Copy(&loginDomain, &loginRequest)

	userDB, err := lg.LoginService.GetUserByUsername(contextControl, loginDomain)
	if err != nil {
		utils.HandleError(w, lg.LoggerSugar, http.StatusInternalServerError, ErrorToGetUserByUsername, err)
		return
	}

	if err = middlewares.VerifyPassword(userDB.Password, loginDomain.Password); err != nil {
		utils.HandleError(w, lg.LoggerSugar, http.StatusUnauthorized, ErrorToVerifyPassword, err)
		return
	}

	token, err := security.CreateToken(userDB.ID)
	if err != nil {
		utils.HandleError(w, lg.LoggerSugar, http.StatusInternalServerError, ErrorToCreateToken, err)
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

	var loginResponse LoginResponse
	copier.Copy(&loginResponse, &loginDomain)

	loginResponse.Token = token

	response := utils.ObjectResponse(loginResponse, SuccessToLogin)
	utils.ResponseReturn(w, http.StatusOK, response.Bytes())
}
