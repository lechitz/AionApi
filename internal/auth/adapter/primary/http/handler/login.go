// Package handler (auth) implements HTTP handlers for authentication endpoints.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/lechitz/AionApi/internal/auth/adapter/primary/http/dto"
	"github.com/lechitz/AionApi/internal/platform/server/http/helpers/httpresponse"
	"github.com/lechitz/AionApi/internal/platform/server/http/helpers/sharederrors"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"github.com/lechitz/AionApi/internal/shared/httputils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Login handles the user login request, validates the credentials, and returns an authentication token.
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
		helpers.WriteDecodeError(ctx, w, span, err, h.Logger)
		return
	}

	if err := validateLoginRequest(loginReq); err != nil {
		helpers.WriteDecodeError(ctx, w, span, sharederrors.NewValidationError("credentials", err.Error()), h.Logger)
		return
	}
	span.SetAttributes(attribute.String(commonkeys.Username, loginReq.Username))

	span.AddEvent(EventAuthServiceLogin)
	user, token, err := h.Service.Login(ctx, loginReq.Username, loginReq.Password)
	if err != nil {
		helpers.WriteDomainError(ctx, w, span, err, ErrLogin, h.Logger)
		return
	}
	span.SetAttributes(attribute.String(commonkeys.UserID, strconv.FormatUint(user.ID, 10)))

	cookies.SetAuthCookie(w, token, h.Config.Cookie)

	loginResponse := dto.LoginUserResponse{
		Name: user.Name,
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
