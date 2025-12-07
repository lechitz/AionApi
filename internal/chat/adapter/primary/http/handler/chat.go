package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/lechitz/AionApi/internal/chat/adapter/primary/http/dto"
	"github.com/lechitz/AionApi/internal/platform/server/http/utils/httpresponse"
	"github.com/lechitz/AionApi/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Chat processes a chat message from the user and returns the AI response.
//
// @Summary      Send chat message
// @Description  Sends a message to the AI assistant and receives a response. Requires authentication.
// @Tags         Chat
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string              true  "Bearer token"
// @Param        chat           body      dto.ChatRequest     true  "Chat message"
// @Success      200            {object}  dto.ChatResponse    "Chat response"
// @Failure      400            {string}  string              "Invalid request payload or validation error"
// @Failure      401            {string}  string              "Unauthorized - missing or invalid token"
// @Failure      500            {string}  string              "Internal server error"
// @Failure      503            {string}  string              "Service unavailable - AI service is down"
// @Router       /chat [post]
// @Security     BearerAuth.
func (h *Handler) Chat(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(TracerChatHandler).
		Start(r.Context(), SpanChatHandler)
	defer span.End()

	userIDValue := ctx.Value(ctxkeys.UserID)
	if userIDValue == nil {
		span.SetStatus(codes.Error, ErrUserIDNotFound)
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound)
		httpresponse.WriteDecodeErrorSpan(ctx, w, span,
			sharederrors.NewAuthenticationError(ErrUserIDNotFound), h.Logger)
		return
	}

	userID, ok := userIDValue.(uint64)
	if !ok {
		span.SetStatus(codes.Error, "invalid user ID type")
		h.Logger.ErrorwCtx(ctx, "Invalid user ID type in context", "value", userIDValue)
		httpresponse.WriteDecodeErrorSpan(ctx, w, span,
			sharederrors.NewAuthenticationError("invalid user ID"), h.Logger)
		return
	}

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(tracingkeys.RequestIPKey, r.RemoteAddr),
	)

	span.AddEvent(EventDecodeRequest)
	r.Body = http.MaxBytesReader(w, r.Body, 1<<21) // 2MB
	var chatReq dto.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&chatReq); err != nil {
		httpresponse.WriteDecodeErrorSpan(ctx, w, span, err, h.Logger)
		return
	}

	span.AddEvent(EventValidateRequest)
	if err := validateChatRequest(chatReq); err != nil {
		httpresponse.WriteDecodeErrorSpan(ctx, w, span,
			sharederrors.NewValidationError("message", err.Error()), h.Logger)
		return
	}

	span.SetAttributes(
		attribute.Int("message_length", len(chatReq.Message)),
	)

	h.Logger.InfowCtx(ctx, MsgChatRequestStart, commonkeys.UserID, strconv.FormatUint(userID, 10), "message_length", len(chatReq.Message))

	span.AddEvent(EventCallService)
	result, err := h.Service.ProcessMessage(ctx, userID, chatReq.Message)
	if err != nil {
		span.AddEvent(EventChatError)
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, ErrChat, h.Logger)
		return
	}

	response := dto.ChatResponse{
		Response: result.Response,
		Sources:  convertToMapSlice(result.Sources),
	}

	if result.TokensUsed > 0 {
		response.Usage = &dto.TokenUsage{
			TotalTokens: result.TokensUsed,
		}
	}

	span.SetAttributes(
		attribute.Int("tokens_used", result.TokensUsed),
		attribute.Int("response_length", len(result.Response)),
		attribute.Int("sources_count", len(result.Sources)),
	)
	span.AddEvent(EventChatSuccess)
	span.SetStatus(codes.Ok, StatusChatSuccess)

	h.Logger.InfowCtx(ctx, MsgChatSuccess,
		commonkeys.UserID, strconv.FormatUint(userID, 10),
		"tokens_used", result.TokensUsed,
		"response_length", len(result.Response),
	)

	httpresponse.WriteSuccess(w, http.StatusOK, response, MsgChatSuccess)
}

// validateChatRequest validates the chat request payload.
func validateChatRequest(req dto.ChatRequest) error {
	msg := strings.TrimSpace(req.Message)

	if msg == "" {
		return errors.New(ErrRequiredMessage)
	}

	if len(msg) < MinMessageLength {
		return errors.New(ErrMessageTooShort)
	}

	if len(msg) > MaxMessageLength {
		return errors.New(ErrMessageTooLong)
	}

	return nil
}

// convertToMapSlice converts []interface{} to []map[string]interface{}.
func convertToMapSlice(sources []interface{}) []map[string]interface{} {
	if sources == nil {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(sources))
	for _, source := range sources {
		if m, ok := source.(map[string]interface{}); ok {
			result = append(result, m)
		}
	}

	return result
}
