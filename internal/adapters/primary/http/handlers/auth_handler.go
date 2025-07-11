package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/lechitz/AionApi/internal/platform/config"

	"github.com/lechitz/AionApi/internal/shared/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/ctxkeys"
	"github.com/lechitz/AionApi/internal/shared/httputils"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/ports/input"
	"github.com/lechitz/AionApi/internal/core/ports/output"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/constants"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/dto"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/response"
)

// TODO: Separar os Handlers e o construtor.. achar um nome bom pro arquivo de NewUser !

// TODO: Melhorar as msgs de Span, Ajustar os erros e logs para ficarem completos !

// Auth provides authentication handlers for login and logout functionalities.
// Combines AuthService for logic and Logger for logging operations.
type Auth struct {
	AuthService input.AuthService
	Logger      output.ContextLogger
	Config      *config.Config
}

// NewAuth initializes and returns a new Auth instance with AuthService and Logger dependencies.
func NewAuth(authService input.AuthService, cfg *config.Config, logger output.ContextLogger) *Auth {
	return &Auth{
		AuthService: authService,
		Config:      cfg,
		Logger:      logger,
	}
}

// LoginHandler handles the user login request, validates the credentials, and returns an authentication token.
func (a *Auth) LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("AionApi/AuthHandler").Start(r.Context(), "LoginHandler")
	defer span.End()

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

	httputils.SetAuthCookie(w, token, a.Config.Cookie)

	loginUserResponse := dto.LoginUserResponse{Username: userDB.Username}
	span.SetAttributes(attribute.String(commonkeys.Username, userDB.Username))

	body := response.ObjectResponse(loginUserResponse, constants.SuccessLogin, a.Logger)
	response.Return(w, http.StatusOK, body.Bytes(), a.Logger)
}

// LogoutHandler processes user logout requests by invalidating tokens, clearing cookies, logging the event, and returning a success response.
func (a *Auth) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("AionApi/AuthHandler").Start(r.Context(), "LogoutHandler")
	defer span.End()

	userID, ok := ctx.Value(ctxkeys.UserID).(uint64)
	if !ok || userID == 0 {
		a.logAndRespondError(w, http.StatusUnauthorized, constants.ErrorToRetrieveUserID, nil)
		return
	}

	tokenVal := ctx.Value(ctxkeys.Token)
	tokenString, ok := tokenVal.(string)
	if !ok || tokenString == "" {
		a.logAndRespondError(w, http.StatusUnauthorized, constants.ErrorToRetrieveToken, nil)
		return
	}

	if err := a.AuthService.Logout(ctx, tokenString); err != nil {
		a.logAndRespondError(w, http.StatusInternalServerError, constants.ErrorToLogout, err)
		return
	}

	httputils.ClearAuthCookie(w, a.Config.Cookie)

	// TODO:passar pra uma outra função.
	tokenPreview := ""
	if len(tokenString) >= 10 {
		tokenPreview = tokenString[:10] + "..."
	}

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.TokenPreview, tokenPreview),
	)

	a.Logger.Infow(constants.SuccessLogout, commonkeys.UserID, strconv.FormatUint(userID, 10), commonkeys.Token, tokenPreview)

	body := response.ObjectResponse(nil, constants.SuccessLogout, a.Logger)
	response.Return(w, http.StatusOK, body.Bytes(), a.Logger)
}

// TODO: deveria melhorar algo nessa parte ?

// logAndRespondError logs an error message and sends an appropriate HTTP response with the specified status, message, and error details.
func (a *Auth) logAndRespondError(w http.ResponseWriter, status int, message string, err error) {
	if err != nil {
		a.Logger.Errorw(message, commonkeys.Error, err.Error())
	} else {
		a.Logger.Errorw(message)
	}
	response.HandleError(w, a.Logger, status, message, err)
}
