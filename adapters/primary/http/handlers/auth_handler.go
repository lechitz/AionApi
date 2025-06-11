package handlers

import (
	"encoding/json"
	"github.com/lechitz/AionApi/internal/shared/contextutil"
	"net/http"
	"time"

	"github.com/lechitz/AionApi/adapters/primary/http/constants"
	"github.com/lechitz/AionApi/adapters/primary/http/dto"
	"github.com/lechitz/AionApi/adapters/primary/http/middleware/response"

	"github.com/lechitz/AionApi/internal/core/domain"
	inputHttp "github.com/lechitz/AionApi/internal/core/ports/input/http"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
)

// Auth provides authentication handlers for login and logout functionalities.
// Combines AuthService for logic and Logger for logging operations.
type Auth struct {
	AuthService inputHttp.AuthService
	Logger      logger.Logger
}

// NewAuth initializes and returns a new Auth instance with AuthService and Logger dependencies.
func NewAuth(authService inputHttp.AuthService, logger logger.Logger) *Auth {
	return &Auth{
		AuthService: authService,
		Logger:      logger,
	}
}

// LoginHandler handles the user login request, validates the credentials, and returns an authentication token.
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

	loginUserResponse := dto.LoginUserResponse{Username: userDB.Username}

	body := response.ObjectResponse(loginUserResponse, constants.SuccessLogin, a.Logger)
	response.Return(w, http.StatusOK, body.Bytes(), a.Logger)
}

// LogoutHandler processes user logout requests by invalidating tokens, clearing cookies, logging the event, and returning a success response.
func (a *Auth) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := contextutil.GetUserID(ctx)
	if !ok || userID == 0 {
		a.logAndRespondError(w, http.StatusUnauthorized, constants.ErrorToRetrieveUserID, nil)
		return
	}

	tokenString, ok := contextutil.GetToken(ctx)
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
		constants.SuccessLogout,
		constants.UserID, userID,
		constants.Token, tokenPreview,
	)

	body := response.ObjectResponse(nil, constants.SuccessLogout, a.Logger)
	response.Return(w, http.StatusOK, body.Bytes(), a.Logger)
}

// logAndRespondError logs an error message and sends an appropriate HTTP response with the specified status, message, and error details.
func (a *Auth) logAndRespondError(w http.ResponseWriter, status int, message string, err error) {
	if err != nil {
		a.Logger.Errorw(message, constants.Error, err.Error())
	} else {
		a.Logger.Errorw(message)
	}
	response.HandleError(w, a.Logger, status, message, err)
}

// setAuthCookie sets a secure HTTP-only authentication cookie with the given token and expiration configuration.
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

// clearAuthCookie invalidates the authentication cookie by setting its value to empty and expiration to a past timestamp.
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
