// Package handler (auth) implements HTTP handlers for authentication endpoints.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/controllers/auth/constants"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/controllers/auth/dto"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"github.com/lechitz/AionApi/internal/shared/handlerhelpers"
	"github.com/lechitz/AionApi/internal/shared/httpresponse"
	"github.com/lechitz/AionApi/internal/shared/httputils"
	"github.com/lechitz/AionApi/internal/shared/sharederrors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Login handles the user login request, validates the credentials, and returns an authentication token.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(constants.TracerAuthHandler).
		Start(r.Context(), constants.SpanLoginHandler)
	defer span.End()

	ip := r.RemoteAddr
	userAgent := r.UserAgent()

	span.SetAttributes(
		attribute.String(tracingkeys.RequestIPKey, ip),
		attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
	)

	span.AddEvent(constants.EventDecodeRequest)
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB
	var loginReq dto.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		handlerhelpers.WriteDecodeError(ctx, w, span, err, h.Logger)
		return
	}

	if err := validateLoginRequest(loginReq); err != nil {
		handlerhelpers.WriteDecodeError(ctx, w, span, sharederrors.NewValidationError("credentials", err.Error()), h.Logger)
		return
	}
	span.SetAttributes(attribute.String(commonkeys.Username, loginReq.Username))

	span.AddEvent(constants.EventAuthServiceLogin)
	user, token, err := h.Service.Login(ctx, loginReq.Username, loginReq.Password)
	if err != nil {
		handlerhelpers.WriteDomainError(ctx, w, span, err, constants.ErrLogin, h.Logger)
		return
	}
	span.SetAttributes(attribute.String(commonkeys.UserID, strconv.FormatUint(user.ID, 10)))

	httputils.SetAuthCookie(w, token, h.Config.Cookie)

	loginResponse := dto.LoginUserResponse{
		Name: user.Name,
	}

	span.AddEvent(constants.EventLoginSuccess)
	span.SetStatus(codes.Ok, constants.StatusLoginSuccess)

	h.Logger.InfowCtx(ctx, constants.MsgLoginSuccess,
		commonkeys.UserID, strconv.FormatUint(user.ID, 10),
		commonkeys.Name, loginResponse.Name,
		tracingkeys.RequestIPKey, ip,
		tracingkeys.RequestUserAgentKey, userAgent,
	)

	httpresponse.WriteSuccess(w, http.StatusOK, loginResponse, constants.MsgLoginSuccess)
}

func validateLoginRequest(req dto.LoginUserRequest) error {
	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		return errors.New(constants.ErrRequiredFields)
	}
	if len(strings.TrimSpace(req.Username)) < constants.MinUsernameLength {
		return errors.New(constants.ErrMinUsernameLength)
	}
	if len(req.Password) < constants.MinPasswordLength {
		return errors.New(constants.ErrMinPasswordLength)
	}
	return nil
}
