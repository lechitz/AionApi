// Package handler user implements HTTP controllers for user-related endpoints.
package handler

import (
	"net/http"
	"strconv"

	"github.com/lechitz/AionApi/internal/platform/server/http/helpers/httpresponse"
	"github.com/lechitz/AionApi/internal/platform/server/http/helpers/sharederrors"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"github.com/lechitz/AionApi/internal/user/adapter/primary/http/dto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetUserByID handles GET /user/{user_id}.
func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(TracerUserHandler).
		Start(r.Context(), SpanGetUserByIDHandler)
	defer span.End()

	ip := r.RemoteAddr
	userAgent := r.UserAgent()

	userID, err := parseUserIDParam(r, h.Logger)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrDecodeGetUserByIDRequest)
		span.SetAttributes(
			attribute.Int(tracingkeys.HTTPStatusCodeKey, http.StatusBadRequest),
			attribute.String(tracingkeys.RequestIPKey, ip),
			attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
		)
		h.Logger.ErrorwCtx(ctx, ErrDecodeGetUserByIDRequest,
			commonkeys.Error, err.Error(),
			tracingkeys.RequestIPKey, ip,
			tracingkeys.RequestUserAgentKey, userAgent,
		)
		helpers.WriteDecodeError(ctx, w, span, err, h.Logger)
		return
	}

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(tracingkeys.RequestIPKey, ip),
		attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
	)

	span.AddEvent(EventUserServiceGetUserByID)

	user, err := h.UserService.GetByID(ctx, userID)
	if err != nil {
		statusCode := sharederrors.MapErrorToHTTPStatus(err)
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrGetUserByID)
		span.SetAttributes(attribute.Int(tracingkeys.HTTPStatusCodeKey, statusCode))
		h.Logger.ErrorwCtx(ctx, ErrGetUserByID,
			commonkeys.UserID, strconv.FormatUint(userID, 10),
			commonkeys.Error, err.Error(),
			tracingkeys.RequestIPKey, ip,
			tracingkeys.RequestUserAgentKey, userAgent,
		)
		helpers.WriteDomainError(ctx, w, span, err, ErrGetUserByID, h.Logger)
		return
	}

	span.SetAttributes(
		attribute.String(commonkeys.Username, user.Username),
		attribute.Int(tracingkeys.HTTPStatusCodeKey, http.StatusOK),
	)
	span.SetStatus(codes.Ok, StatusUserFetched)

	h.Logger.InfowCtx(ctx, MsgUserFetched,
		commonkeys.UserID, strconv.FormatUint(user.ID, 10),
		commonkeys.Username, user.Username,
		commonkeys.Email, user.Email,
		tracingkeys.RequestIPKey, ip,
		tracingkeys.RequestUserAgentKey, userAgent,
	)

	span.AddEvent(EventUserFetchedSuccess)

	res := dto.GetUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	httpresponse.WriteSuccess(w, http.StatusOK, res, MsgUserFetched)
}
