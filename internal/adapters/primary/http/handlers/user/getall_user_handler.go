// Package user implements HTTP handlers for user-related endpoints.
package user

import (
	"net/http"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/user/constants"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/user/dto"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"github.com/lechitz/AionApi/internal/shared/handlerhelpers"
	"github.com/lechitz/AionApi/internal/shared/httpresponse"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// GetAllUsers handles GET /user/all.
func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(constants.TracerUserHandler).
		Start(r.Context(), constants.SpanGetAllUsersHandler)
	defer span.End()

	ip := r.RemoteAddr
	userAgent := r.UserAgent()

	span.AddEvent(constants.EventUserServiceGetAllUsers,
		trace.WithAttributes(
			attribute.String(tracingkeys.RequestIPKey, ip),
			attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
		),
	)

	users, err := h.Service.GetAllUsers(ctx)
	if err != nil {
		h.Logger.ErrorwCtx(ctx, constants.ErrGetUsers,
			tracingkeys.RequestIPKey, ip,
			tracingkeys.RequestUserAgentKey, userAgent,
			commonkeys.Error, err.Error(),
		)
		handlerhelpers.WriteDomainError(ctx, w, span, err, constants.ErrGetUsers, h.Logger)
		return
	}

	span.SetAttributes(
		attribute.Int(commonkeys.UsersCount, len(users)),
		attribute.Int(tracingkeys.HTTPStatusCodeKey, http.StatusOK),
		attribute.String(tracingkeys.RequestIPKey, ip),
		attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
	)
	span.SetStatus(codes.Ok, constants.StatusUsersFetched)

	h.Logger.InfowCtx(ctx, constants.MsgUsersFetched,
		commonkeys.UsersCount, len(users),
		tracingkeys.RequestIPKey, ip,
		tracingkeys.RequestUserAgentKey, userAgent,
	)

	span.AddEvent(constants.EventUsersFetchedSuccess,
		trace.WithAttributes(attribute.Int(commonkeys.UsersCount, len(users))),
	)

	res := make([]dto.GetUserResponse, 0, len(users))
	for _, user := range users {
		res = append(res, dto.GetUserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		})
	}

	httpresponse.WriteSuccess(w, http.StatusOK, res, constants.MsgUsersFetched)
}
