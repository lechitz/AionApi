package handler

import (
	"context"
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

// ChatText processes a chat message from the user and returns the AI response.
//
// @Summary      Send chat message
// @Description  Sends a message to the AI assistant and receives a response. Requires authentication.
// @Tags         ChatText
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string              true  "Bearer token"
// @Param        chat           body      dto.ChatRequest     true  "ChatText message"
// @Success      200            {object}  dto.ChatResponse    "ChatText response"
// @Failure      400            {string}  string              "Invalid request payload or validation error"
// @Failure      401            {string}  string              "Unauthorized - missing or invalid token"
// @Failure      500            {string}  string              "Internal server error"
// @Failure      503            {string}  string              "Service unavailable - AI service is down"
// @Router       /chat/text [post]
// @Security     BearerAuth.
func (h *Handler) ChatText(w http.ResponseWriter, r *http.Request) {
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
		span.SetStatus(codes.Error, ErrInvalidUserIDType)
		h.Logger.ErrorwCtx(ctx, LogInvalidUserIDType, LogKeyValue, userIDValue)
		httpresponse.WriteDecodeErrorSpan(ctx, w, span,
			sharederrors.NewAuthenticationError(ErrInvalidUserID), h.Logger)
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
			sharederrors.NewValidationError(FormFieldMessage, err.Error()), h.Logger)
		return
	}

	span.SetAttributes(
		attribute.Int(AttrMessageLength, len(chatReq.Message)),
	)
	h.logUIActionMetadata(ctx, userID, chatReq.Context)

	h.Logger.InfowCtx(ctx, MsgChatRequestStart, commonkeys.UserID, strconv.FormatUint(userID, 10), AttrMessageLength, len(chatReq.Message))

	span.AddEvent(EventCallService)
	result, err := h.Service.ProcessMessage(ctx, userID, chatReq.Message, chatReq.Context)
	if err != nil {
		if isClientCancelledError(err) {
			span.AddEvent(EventChatCancelled)
			span.SetStatus(codes.Ok, StatusChatCancelled)
			h.Logger.InfowCtx(ctx, MsgChatCancelled,
				commonkeys.UserID, strconv.FormatUint(userID, 10),
			)
			httpresponse.WriteSuccess(w, HTTPStatusClientClosedRequest, nil, MsgChatCancelledResponse)
			return
		}
		span.AddEvent(EventChatError)
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, ErrChat, h.Logger)
		return
	}

	response := dto.ChatResponse{
		Response: result.Response,
		UI:       result.UI,
		Sources:  convertToMapSlice(result.Sources),
	}

	if result.TokensUsed > 0 {
		response.Usage = &dto.TokenUsage{
			TotalTokens: result.TokensUsed,
		}
	}

	span.SetAttributes(
		attribute.Int(AttrTokensUsed, result.TokensUsed),
		attribute.Int(AttrResponseLength, len(result.Response)),
		attribute.Int(AttrSourcesCount, len(result.Sources)),
	)
	span.AddEvent(EventChatSuccess)
	span.SetStatus(codes.Ok, StatusChatSuccess)

	h.Logger.InfowCtx(ctx, MsgChatSuccess,
		commonkeys.UserID, strconv.FormatUint(userID, 10),
		AttrTokensUsed, result.TokensUsed,
		AttrResponseLength, len(result.Response),
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

func isClientCancelledError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, context.Canceled) {
		return true
	}

	text := strings.ToLower(err.Error())
	return strings.Contains(text, ErrorTextContextCanceled) ||
		strings.Contains(text, ErrorTextRequestCanceled) ||
		strings.Contains(text, ErrorTextOperationCanceled)
}

func (h *Handler) logUIActionMetadata(
	ctx context.Context,
	userID uint64,
	requestContext map[string]interface{},
) {
	if requestContext == nil {
		return
	}
	rawAction, ok := requestContext[ContextKeyUIAction].(map[string]interface{})
	if !ok || rawAction == nil {
		return
	}

	actionType, _ := rawAction[ContextKeyUIActionType].(string)
	draftID, _ := rawAction[ContextKeyDraftID].(string)
	consentRequired := false
	consentConfirmed := false
	consentPolicyVersion := ""
	if rawConsent, ok := rawAction[ContextKeyConsent].(map[string]interface{}); ok && rawConsent != nil {
		if value, ok := rawConsent["required"].(bool); ok {
			consentRequired = value
		}
		if value, ok := rawConsent["confirmed"].(bool); ok {
			consentConfirmed = value
		}
		if value, ok := rawConsent["policy_version"].(string); ok {
			consentPolicyVersion = value
		}
	}
	h.Logger.InfowCtx(
		ctx,
		MsgChatRequestIncludesUIAction,
		commonkeys.UserID, strconv.FormatUint(userID, 10),
		LogKeyUIActionType, actionType,
		LogKeyDraftID, draftID,
		LogKeyConsentRequired, consentRequired,
		LogKeyConsentConfirmed, consentConfirmed,
		LogKeyConsentPolicyVersion, consentPolicyVersion,
	)
}
