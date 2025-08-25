// Package handler user controllers provide HTTP controllers for user-related endpoints.
package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/controllers/user/constants"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/controllers/user/dto"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"github.com/lechitz/AionApi/internal/shared/handlerhelpers"
	"github.com/lechitz/AionApi/internal/shared/httpresponse"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// maxBodyBytes defines the maximum allowed size for request body in bytes (1MB)
// to prevent memory exhaustion attacks and ensure reasonable payload sizes
const maxBodyBytes = 1 << 20 // 1MB

// Create handles POST /user/create.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	tr := otel.Tracer(constants.TracerUserHandler)
	ctx, span := tr.Start(r.Context(), constants.SpanCreateUserHandler)
	defer span.End()

	ip := r.RemoteAddr
	userAgent := r.UserAgent()

	span.AddEvent(constants.EventRequestUserAgentKeyAndIP,
		trace.WithAttributes(
			attribute.String(tracingkeys.RequestIPKey, ip),
			attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
		),
	)

	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.SetStatus(codes.Error, constants.EventDecodeRequest)
		h.Logger.ErrorwCtx(ctx, constants.ErrCreateUserValidation)
		handlerhelpers.WriteDecodeError(ctx, w, span, err, h.Logger)
		return
	}

	if err := req.ValidateUser(); err != nil {
		span.SetStatus(codes.Error, constants.ErrCreateUserValidation)
		h.Logger.ErrorwCtx(ctx, constants.ErrCreateUserValidation)
		handlerhelpers.WriteDecodeError(ctx, w, span, err, h.Logger)
		return
	}

	span.SetAttributes(
		attribute.String(commonkeys.Username, req.Username),
		attribute.String(commonkeys.Email, req.Email),
	)

	cmd := req.ToCommand()

	span.AddEvent(constants.EventUserServiceCreateUser)

	userCreated, err := h.UserService.Create(ctx, cmd)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		handlerhelpers.WriteDomainError(ctx, w, span, err, constants.ErrCreateUser, h.Logger)
		return
	}

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userCreated.ID, 10)),
		attribute.Int(tracingkeys.HTTPStatusCodeKey, http.StatusCreated),
	)
	span.SetStatus(codes.Ok, constants.StatusUserCreated)

	h.Logger.InfowCtx(ctx, constants.MsgUserCreated,
		commonkeys.UserID, strconv.FormatUint(userCreated.ID, 10),
		commonkeys.Username, userCreated.Username,
		commonkeys.Email, userCreated.Email,
		tracingkeys.RequestIPKey, ip,
		tracingkeys.RequestUserAgentKey, userAgent,
	)

	span.AddEvent(
		constants.EventUserCreatedSuccess,
		trace.WithAttributes(attribute.String(commonkeys.UserID, strconv.FormatUint(userCreated.ID, 10))),
	)

	response := dto.CreateUserResponse{
		Name:     userCreated.Name,
		Username: userCreated.Username,
		Email:    userCreated.Email,
		ID:       userCreated.ID,
	}

	httpresponse.WriteSuccess(w, http.StatusCreated, response, constants.MsgUserCreated)
}
