// Package handler user implements HTTP controllers for user-related endpoints.
package handler

import (
	"net/http"

	"github.com/lechitz/AionApi/internal/platform/server/http/utils/httpresponse"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"github.com/lechitz/AionApi/internal/user/adapter/primary/http/dto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ListAll returns all users.
//
// @Summary      List all users
// @Description  Returns a list containing all user resources.
// @Tags         Users
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Success      200  {array}   dto.GetUserResponse  "Users fetched"
// @Failure      401  {string}  string               "Unauthorized"
// @Failure      500  {string}  string               "Internal server error"
// @Router       /user/all [get].
func (h *Handler) ListAll(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(TracerUserHandler).
		Start(r.Context(), SpanGetAllUsersHandler)
	defer span.End()

	ip := r.RemoteAddr
	userAgent := r.UserAgent()

	span.AddEvent(EventUserServiceGetAllUsers,
		trace.WithAttributes(
			attribute.String(tracingkeys.RequestIPKey, ip),
			attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
		),
	)

	users, err := h.UserService.ListAll(ctx)
	if err != nil {
		h.Logger.ErrorwCtx(ctx, ErrGetUsers,
			tracingkeys.RequestIPKey, ip,
			tracingkeys.RequestUserAgentKey, userAgent,
			commonkeys.Error, err.Error(),
		)
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, ErrGetUsers, h.Logger)
		return
	}

	span.SetAttributes(
		attribute.Int(commonkeys.UsersCount, len(users)),
		attribute.Int(tracingkeys.HTTPStatusCodeKey, http.StatusOK),
		attribute.String(tracingkeys.RequestIPKey, ip),
		attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
	)
	span.SetStatus(codes.Ok, StatusUsersFetched)

	h.Logger.InfowCtx(ctx, MsgUsersFetched,
		commonkeys.UsersCount, len(users),
		tracingkeys.RequestIPKey, ip,
		tracingkeys.RequestUserAgentKey, userAgent,
	)

	span.AddEvent(EventUsersFetchedSuccess,
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

	httpresponse.WriteSuccess(w, http.StatusOK, res, MsgUsersFetched)
}
