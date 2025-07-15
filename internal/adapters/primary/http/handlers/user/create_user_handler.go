// Package user handlers provides HTTP handlers for user-related endpoints.
package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/user/constants"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/user/dto"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"github.com/lechitz/AionApi/internal/shared/handlerhelpers"
	"github.com/lechitz/AionApi/internal/shared/httpresponse"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Create handles POST /user/create.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(constants.TracerUserHandler).
		Start(r.Context(), constants.SpanCreateUserHandler)
	defer span.End()

	ip := r.RemoteAddr
	userAgent := r.UserAgent()

	span.AddEvent(constants.EventDecodeRequest,
		trace.WithAttributes(
			attribute.String(tracingkeys.RequestIPKey, ip),
			attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
		),
	)

	var userDTO dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		handlerhelpers.WriteDecodeError(ctx, w, span, err, h.Logger)
		return
	}

	//TODO: passar validação pro DTO.
	err := handlerhelpers.CheckRequiredFields(map[string]string{
		commonkeys.Name:     userDTO.Name,
		commonkeys.Username: userDTO.Username,
		commonkeys.Email:    userDTO.Email,
		commonkeys.Password: userDTO.Password,
	})
	if err != nil {
		h.Logger.ErrorwCtx(ctx, constants.ErrCreateUserValidation,
			commonkeys.Error, err.Error(),
			tracingkeys.RequestIPKey, ip,
			tracingkeys.RequestUserAgentKey, userAgent,
		)
		handlerhelpers.WriteDecodeError(ctx, w, span, err, h.Logger)
		return
	}

	span.SetAttributes(
		attribute.String(commonkeys.Username, userDTO.Username),
		attribute.String(commonkeys.Email, userDTO.Email),
		attribute.String(tracingkeys.RequestIPKey, ip),
		attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
	)

	userDomain := domain.User{
		Name:     userDTO.Name,
		Username: userDTO.Username,
		Email:    userDTO.Email,
		Password: userDTO.Password,
	}

	span.AddEvent(constants.EventUserServiceCreateUser)
	newUser, err := h.Service.CreateUser(ctx, userDomain)
	if err != nil {
		handlerhelpers.WriteDomainError(ctx, w, span, err, constants.ErrCreateUser, h.Logger)
		return
	}

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(newUser.ID, 10)),
		attribute.Int(tracingkeys.HTTPStatusCodeKey, http.StatusCreated),
	)
	span.SetStatus(codes.Ok, constants.StatusUserCreated)

	h.Logger.InfowCtx(ctx, constants.MsgUserCreated,
		commonkeys.UserID, strconv.FormatUint(newUser.ID, 10),
		commonkeys.Username, newUser.Username,
		commonkeys.Email, newUser.Email,
		tracingkeys.RequestIPKey, ip,
		tracingkeys.RequestUserAgentKey, userAgent,
	)

	span.AddEvent(
		constants.EventUserCreatedSuccess,
		trace.WithAttributes(attribute.String(commonkeys.UserID, strconv.FormatUint(newUser.ID, 10))),
	)

	response := dto.CreateUserResponse{
		Name:     newUser.Name,
		Username: newUser.Username,
		Email:    newUser.Email,
		ID:       newUser.ID,
	}

	httpresponse.WriteSuccess(w, http.StatusCreated, response, constants.MsgUserCreated)
}
