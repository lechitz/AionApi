// Package handler (auth) implements HTTP handlers for authentication endpoints.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/lechitz/AionApi/internal/auth/adapter/primary/http/dto"
	"github.com/lechitz/AionApi/internal/platform/server/http/utils/cookies"
	"github.com/lechitz/AionApi/internal/platform/server/http/utils/httpresponse"
	"github.com/lechitz/AionApi/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Login authenticates a user and issues a session token.
//
// @Summary      Authenticate user (login)
// @Description  Validates credentials and issues a session token (JWT or cookie). On success, returns user-facing info and sets `auth_token` cookie.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        login  body      dto.LoginUserRequest  true  "Login payload"
// @Success      200    {object}  dto.LoginUserResponse  "Login succeeded"
// @Header       200    {string}  Set-Cookie  "auth_token=<opaque or JWT>; Path=/; HttpOnly; Secure (if enabled)"
// @Failure      400    {string}  string  "Invalid request payload or validation error"
// @Failure      401    {string}  string  "Invalid credentials"
// @Failure      500    {string}  string  "Internal server error"
// @Router       /auth/login [post].
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(TracerAuthHandler).
		Start(r.Context(), SpanLoginHandler)
	defer span.End()

	ip := r.RemoteAddr
	userAgent := r.UserAgent()

	span.SetAttributes(
		attribute.String(tracingkeys.RequestIPKey, ip),
		attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
	)

	span.AddEvent(EventDecodeRequest)
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB
	var loginReq dto.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		httpresponse.WriteDecodeErrorSpan(ctx, w, span, err, h.Logger)
		return
	}

	if err := validateLoginRequest(loginReq); err != nil {
		httpresponse.WriteDecodeErrorSpan(ctx, w, span, sharederrors.NewValidationError("credentials", err.Error()), h.Logger)
		return
	}
	span.SetAttributes(attribute.String(commonkeys.Username, loginReq.Username))

	span.AddEvent(EventAuthServiceLogin)
	user, token, err := h.Service.Login(ctx, loginReq.Username, loginReq.Password)
	if err != nil {
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, ErrLogin, h.Logger)
		return
	}
	span.SetAttributes(attribute.String(commonkeys.UserID, strconv.FormatUint(user.ID, 10)))

	cookies.SetAuthCookie(w, token, h.Config.Cookie)

	loginResponse := dto.LoginUserResponse{
		Token: token,
		ID:    user.ID,
		Name:  user.Name,
		Roles: user.Roles,
	}

	span.AddEvent(EventLoginSuccess)
	span.SetStatus(codes.Ok, StatusLoginSuccess)

	h.Logger.InfowCtx(ctx, MsgLoginSuccess,
		commonkeys.UserID, strconv.FormatUint(user.ID, 10),
		commonkeys.Name, loginResponse.Name,
		tracingkeys.RequestIPKey, ip,
		tracingkeys.RequestUserAgentKey, userAgent,
	)

	httpresponse.WriteSuccess(w, http.StatusOK, loginResponse, MsgLoginSuccess)
}

// validateLoginRequest checks the login request payload for required fields
// and enforces minimum length constraints on username and password.
func validateLoginRequest(req dto.LoginUserRequest) error {
	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		return errors.New(ErrRequiredFields)
	}
	if len(strings.TrimSpace(req.Username)) < MinUsernameLength {
		return errors.New(ErrMinUsernameLength)
	}
	if len(req.Password) < MinPasswordLength {
		return errors.New(ErrMinPasswordLength)
	}
	return nil
}
