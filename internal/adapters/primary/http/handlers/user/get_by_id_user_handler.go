// Package user implements HTTP handlers for user-related endpoints.
package user

import (
	"net/http"
	"strconv"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/user/constants"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/user/dto"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"github.com/lechitz/AionApi/internal/shared/handlerhelpers"
	"github.com/lechitz/AionApi/internal/shared/httpresponse"
	"github.com/lechitz/AionApi/internal/shared/sharederrors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetUserByID handles GET /user/{user_id}.
func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(constants.TracerUserHandler).
		Start(r.Context(), constants.SpanGetUserByIDHandler)
	defer span.End()

	ip := r.RemoteAddr
	userAgent := r.UserAgent()

	userID, err := parseUserIDParam(r, h.Logger)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, constants.ErrDecodeGetUserByIDRequest)
		span.SetAttributes(
			attribute.Int(tracingkeys.HTTPStatusCodeKey, http.StatusBadRequest),
			attribute.String(tracingkeys.RequestIPKey, ip),
			attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
		)
		h.Logger.ErrorwCtx(ctx, constants.ErrDecodeGetUserByIDRequest,
			commonkeys.Error, err.Error(),
			tracingkeys.RequestIPKey, ip,
			tracingkeys.RequestUserAgentKey, userAgent,
		)
		handlerhelpers.WriteDecodeError(ctx, w, span, err, h.Logger)
		return
	}

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(tracingkeys.RequestIPKey, ip),
		attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
	)

	span.AddEvent(constants.EventUserServiceGetUserByID)

	user, err := h.Service.GetUserByID(ctx, userID)
	if err != nil {
		statusCode := sharederrors.MapErrorToHTTPStatus(err)
		span.RecordError(err)
		span.SetStatus(codes.Error, constants.ErrGetUserByID)
		span.SetAttributes(attribute.Int(tracingkeys.HTTPStatusCodeKey, statusCode))
		h.Logger.ErrorwCtx(ctx, constants.ErrGetUserByID,
			commonkeys.UserID, strconv.FormatUint(userID, 10),
			commonkeys.Error, err.Error(),
			tracingkeys.RequestIPKey, ip,
			tracingkeys.RequestUserAgentKey, userAgent,
		)
		handlerhelpers.WriteDomainError(ctx, w, span, err, constants.ErrGetUserByID, h.Logger)
		return
	}

	span.SetAttributes(
		attribute.String(commonkeys.Username, user.Username),
		attribute.Int(tracingkeys.HTTPStatusCodeKey, http.StatusOK),
	)
	span.SetStatus(codes.Ok, constants.StatusUserFetched)

	h.Logger.InfowCtx(ctx, constants.MsgUserFetched,
		commonkeys.UserID, strconv.FormatUint(user.ID, 10),
		commonkeys.Username, user.Username,
		commonkeys.Email, user.Email,
		tracingkeys.RequestIPKey, ip,
		tracingkeys.RequestUserAgentKey, userAgent,
	)

	span.AddEvent(constants.EventUserFetchedSuccess)

	res := dto.GetUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	httpresponse.WriteSuccess(w, http.StatusOK, res, constants.MsgUserFetched)
}
