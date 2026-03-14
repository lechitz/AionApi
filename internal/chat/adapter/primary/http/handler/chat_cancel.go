package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/lechitz/AionApi/internal/chat/adapter/primary/http/dto"
	"github.com/lechitz/AionApi/internal/platform/server/http/utils/httpresponse"
	"github.com/lechitz/AionApi/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type cancelPayload struct {
	UserID uint64 `json:"user_id"`
}

// ChatCancel cancels active AI processing for the authenticated user.
//
// @Summary      Cancel active chat processing
// @Description  Cancels ongoing AI processing for the current authenticated user.
// @Tags         ChatText
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                  true  "Bearer token"
// @Success      200            {object}  dto.ChatCancelResponse  "Cancellation status"
// @Failure      401            {string}  string                  "Unauthorized - missing or invalid token"
// @Failure      503            {string}  string                  "Service unavailable - AI service is down"
// @Router       /chat/cancel [post]
// @Security     BearerAuth.
func (h *Handler) ChatCancel(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(TracerChatHandler).Start(r.Context(), "chat.handler.cancel")
	defer span.End()

	userIDValue := ctx.Value(ctxkeys.UserID)
	if userIDValue == nil {
		span.SetStatus(codes.Error, ErrUserIDNotFound)
		httpresponse.WriteDecodeErrorSpan(ctx, w, span,
			sharederrors.NewAuthenticationError(ErrUserIDNotFound), h.Logger)
		return
	}

	userID, ok := userIDValue.(uint64)
	if !ok {
		span.SetStatus(codes.Error, ErrInvalidUserIDType)
		httpresponse.WriteDecodeErrorSpan(ctx, w, span,
			sharederrors.NewAuthenticationError(ErrInvalidUserID), h.Logger)
		return
	}

	span.SetAttributes(attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)))

	payloadBytes, err := json.Marshal(cancelPayload{UserID: userID})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "cancel_payload_marshal_failed")
		httpresponse.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to prepare cancel request",
		})
		return
	}

	url := h.Config.AionChat.BaseURL + "/internal/cancel"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "cancel_request_create_failed")
		httpresponse.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{
			"error": "Failed to cancel chat processing",
		})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: h.Config.AionChat.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "cancel_request_failed")
		h.Logger.ErrorwCtx(ctx, "Failed to call Aion-Chat cancel endpoint", commonkeys.Error, err.Error(), commonkeys.URL, url)
		httpresponse.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{
			"error": "Failed to cancel chat processing",
		})
		return
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		span.SetAttributes(attribute.Int("aion_chat.cancel_status_code", resp.StatusCode))
		span.SetStatus(codes.Error, "cancel_non_200")
		httpresponse.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{
			"error": "Failed to cancel chat processing",
		})
		return
	}

	span.SetStatus(codes.Ok, "cancelled")
	httpresponse.WriteSuccess(w, http.StatusOK, dto.ChatCancelResponse{
		Cancelled: true,
		Message:   "Cancel request sent",
	}, "Chat cancel request sent")
}
