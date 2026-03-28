package handler

import (
	"net/http"
	"strconv"

	"github.com/lechitz/aion-api/internal/platform/server/http/utils/httpresponse"
	"github.com/lechitz/aion-api/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"github.com/lechitz/aion-api/internal/shared/constants/ctxkeys"
	"github.com/lechitz/aion-api/internal/shared/constants/tracingkeys"
	"github.com/lechitz/aion-api/internal/user/adapter/primary/http/dto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// DeleteAvatar removes the authenticated user's current avatar reference.
//
// @Summary      Delete current user avatar
// @Description  Removes the authenticated user's avatar and returns the updated profile payload.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Success      200  {object}  dto.UpdateUserResponse  "Avatar removed"
// @Failure      401  {string}  string                  "Unauthorized or missing user context"
// @Failure      404  {string}  string                  "User not found"
// @Failure      500  {string}  string                  "Internal server error"
// @Router       /user/avatar [delete].
func (h *Handler) DeleteAvatar(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(TracerUserHandler).
		Start(r.Context(), SpanDeleteAvatarHandler)
	defer span.End()

	ip := r.RemoteAddr
	userAgent := r.UserAgent()

	userID, ok := ctx.Value(ctxkeys.UserID).(uint64)
	if !ok || userID == 0 {
		err := sharederrors.ErrMissingUserID()
		httpresponse.WriteAuthErrorSpan(ctx, w, span, err, h.Logger)
		return
	}

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(tracingkeys.RequestIPKey, ip),
		attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
	)

	span.AddEvent(EventUserServiceDeleteAvatar,
		trace.WithAttributes(
			attribute.String(tracingkeys.RequestIPKey, ip),
			attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
		),
	)

	userUpdated, err := h.UserService.RemoveAvatar(ctx, userID)
	if err != nil {
		h.Logger.ErrorwCtx(ctx, ErrDeleteAvatar,
			tracingkeys.RequestIPKey, ip,
			tracingkeys.RequestUserAgentKey, userAgent,
			commonkeys.UserID, strconv.FormatUint(userID, 10),
			commonkeys.Error, err.Error(),
		)
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, ErrDeleteAvatar, h.Logger)
		return
	}

	span.SetAttributes(attribute.Int(tracingkeys.HTTPStatusCodeKey, http.StatusOK))
	span.SetStatus(codes.Ok, StatusUserAvatarDeleted)
	span.AddEvent(EventUserAvatarDeletedSuccess,
		trace.WithAttributes(attribute.String(commonkeys.UserID, strconv.FormatUint(userUpdated.ID, 10))),
	)

	h.Logger.InfowCtx(ctx, MsgUserAvatarDeleted,
		commonkeys.UserID, strconv.FormatUint(userUpdated.ID, 10),
		commonkeys.Username, userUpdated.Username,
		commonkeys.Email, userUpdated.Email,
		tracingkeys.RequestIPKey, ip,
		tracingkeys.RequestUserAgentKey, userAgent,
	)

	httpresponse.WriteSuccess(w, http.StatusOK, dto.UpdateUserResponse{
		ID:                  userUpdated.ID,
		Name:                &userUpdated.Name,
		Username:            &userUpdated.Username,
		Email:               &userUpdated.Email,
		UpdatedAt:           userUpdated.UpdatedAt,
		Locale:              userUpdated.Locale,
		Timezone:            userUpdated.Timezone,
		Location:            userUpdated.Location,
		Bio:                 userUpdated.Bio,
		AvatarURL:           userUpdated.AvatarURL,
		OnboardingCompleted: userUpdated.OnboardingCompleted,
	}, MsgUserAvatarDeleted)
}
