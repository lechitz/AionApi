package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/constants"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/dto"
	"github.com/lechitz/AionApi/internal/core/domain"
	inputHttp "github.com/lechitz/AionApi/internal/core/ports/input/http"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	"github.com/lechitz/AionApi/pkg/utils"
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
		a.logAndRespondError(w, http.StatusBadRequest, constants.ErrorToDecodeLoginRequest, err)
		return
	}

	userDomain := domain.UserDomain{Username: loginReq.Username}

	userDB, token, err := a.AuthService.Login(ctx, userDomain, loginReq.Password)
	if err != nil {
		a.logAndRespondError(w, http.StatusInternalServerError, constants.ErrorToLogin, err)
		return
	}

	setAuthCookie(w, token, 0)

	response := dto.LoginUserResponse{Username: userDB.Username}

	utils.ResponseReturn(w, http.StatusOK, utils.ObjectResponse(response, constants.SuccessToLogin).Bytes())
}

func (a *Auth) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value(constants.UserID).(uint64)
	if !ok || userID == 0 {
		a.logAndRespondError(w, http.StatusUnauthorized, constants.ErrorToRetrieveUserID, nil)
		return
	}

	tokenString, ok := ctx.Value(constants.Token).(string)
	if !ok || tokenString == "" {
		a.logAndRespondError(w, http.StatusUnauthorized, constants.ErrorToRetrieveToken, nil)
		return
	}

	if err := a.AuthService.Logout(ctx, tokenString); err != nil {
		a.logAndRespondError(w, http.StatusInternalServerError, constants.ErrorToLogout, err)
		return
	}

	clearAuthCookie(w)

	tokenPreview := ""
	if len(tokenString) >= 10 {
		tokenPreview = tokenString[:10] + "..."
	}

	a.Logger.Infow(
		constants.SuccessToLogout,
		constants.UserID, userID,
		constants.Token, tokenPreview,
	)

	utils.ResponseReturn(w, http.StatusOK, utils.ObjectResponse(nil, constants.SuccessToLogout).Bytes())
}

func (a *Auth) logAndRespondError(w http.ResponseWriter, status int, message string, err error) {
	if err != nil {
		a.Logger.Errorw(message, constants.Error, err.Error())
	} else {
		a.Logger.Errorw(message)
	}
	utils.HandleError(w, a.Logger, status, message, err)
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
