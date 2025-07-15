// Package auth implements HTTP handlers for authentication endpoints.
package auth

import (
	"encoding/json"
	"net/http"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/auth/constants"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/auth/dto"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"github.com/lechitz/AionApi/internal/shared/handlerhelpers"
	"github.com/lechitz/AionApi/internal/shared/httpresponse"
	"github.com/lechitz/AionApi/internal/shared/httputils"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Login handles the user login request, validates the credentials, and returns an authentication token.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(constants.TracerAuthHandler).
		Start(r.Context(), constants.SpanLoginHandler)
	defer span.End()

	ip := r.RemoteAddr
	userAgent := r.UserAgent()

	span.AddEvent(constants.EventDecodeRequest,
		trace.WithAttributes(
			attribute.String(tracingkeys.RequestIPKey, ip),
			attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
		),
	)

	var loginReq dto.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		handlerhelpers.WriteDecodeError(ctx, w, span, err, h.Logger)
		return
	}

	//TODO: adicionar validação basica da quantidade de caractere em user e senha.

	span.SetAttributes(
		attribute.String(commonkeys.Username, loginReq.Username),
		attribute.String(tracingkeys.RequestIPKey, ip),
		attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
	)

	span.AddEvent(constants.EventAuthServiceLogin)

	user, token, err := h.Service.Login(ctx, loginReq.Username, loginReq.Password)
	if err != nil {
		handlerhelpers.WriteDomainError(ctx, w, span, err, constants.ErrLogin, h.Logger)
		return
	}

	httputils.SetAuthCookie(w, token, h.Config.Cookie)

	loginResponse := dto.LoginUserResponse{Name: user.Name}
	span.SetAttributes(
		attribute.String(commonkeys.Name, loginResponse.Name),
		attribute.String(tracingkeys.RequestIPKey, ip),
		attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
	)
	span.SetStatus(codes.Ok, constants.StatusLoginSuccess)

	span.AddEvent(constants.EventLoginSuccess)

	h.Logger.InfowCtx(ctx, constants.MsgLoginSuccess,
		commonkeys.Name, loginResponse.Name,
		tracingkeys.RequestIPKey, ip,
		tracingkeys.RequestUserAgentKey, userAgent,
	)

	httpresponse.WriteSuccess(w, http.StatusOK, loginResponse, constants.MsgLoginSuccess)
}
