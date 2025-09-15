// Package handler user implements HTTP controllers for user-related endpoints.
package handler

import (
	"net/http"
	"strconv"

	"github.com/lechitz/AionApi/internal/platform/server/http/utils/httpresponse"
	"github.com/lechitz/AionApi/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"github.com/lechitz/AionApi/internal/user/adapter/primary/http/dto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetUserByID retrieves a user by its ID.
//
// @Summary      Get user by ID
// @Description  Returns a single user resource by its unique identifier.
// @Tags         Users
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        user_id  path      integer  true  "User ID (uint64)"
// @Success      200      {object}  dto.GetUserResponse  "User fetched"
// @Failure      400      {string}  string                "Invalid user_id"
// @Failure      404      {string}  string                "User not found"
// @Failure      500      {string}  string                "Internal server error"
// @Router       /user/{user_id} [get].
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
		httpresponse.WriteDecodeErrorSpan(ctx, w, span, err, h.Logger)
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
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, ErrGetUserByID, h.Logger)
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
