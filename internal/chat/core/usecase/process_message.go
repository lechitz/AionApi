// Package usecase implements the chat use cases (business logic).
package usecase

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/lechitz/AionApi/internal/chat/adapter/primary/http/dto"
	"github.com/lechitz/AionApi/internal/chat/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ProcessMessage processes a chat message by forwarding it to the Aion-Chat service.
func (s *ChatService) ProcessMessage(ctx context.Context, userID uint64, message string) (*domain.ChatResult, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanProcessMessage)
	defer span.End()

	span.SetAttributes(
		attribute.String(AttrUserID, strconv.FormatUint(userID, 10)),
		attribute.Int(AttrMessageLength, len(message)),
	)

	s.logger.InfowCtx(ctx, LogProcessingChatMessage, LogKeyUserID, userID, LogKeyMessageLength, len(message))

	// Fetch recent conversation history from CACHE (Redis) for fast access
	// Retrieve last 6 messages (3 exchanges) for context
	// Cache is much faster than database and perfect for volatile conversation context
	const historyLimit = 6
	conversationHistory := []dto.ConversationMessage{}

	history, err := s.chatHistoryCache.GetLatest(ctx, userID, historyLimit)
	if err != nil {
		// Don't fail on cache retrieval error, just log and continue without history
		s.logger.WarnwCtx(ctx, "Failed to retrieve conversation history from cache (non-blocking)",
			LogKeyError, err.Error(),
			LogKeyUserID, userID,
		)
	} else if len(history) > 0 {
		// Convert history to conversation messages format (reverse order - oldest first)
		// History comes DESC (newest first), we need ASC (oldest first) for LLM context
		for i := len(history) - 1; i >= 0; i-- {
			h := history[i]
			// Add user message
			conversationHistory = append(conversationHistory, dto.ConversationMessage{
				Role:    "user",
				Content: h.Message,
			})
			// Add assistant response
			conversationHistory = append(conversationHistory, dto.ConversationMessage{
				Role:    "assistant",
				Content: h.Response,
			})
		}

		s.logger.InfowCtx(ctx, "Including conversation history from cache",
			LogKeyUserID, userID,
			"history_messages", len(conversationHistory),
		)
	}

	req := &dto.InternalChatRequest{
		UserID:              userID,
		Message:             message,
		ConversationHistory: conversationHistory,
		Context: map[string]interface{}{
			ContextKeyTimezone: DefaultTimezone, // TODO: Get from user settings
		},
	}

	resp, err := s.aionChatClient.SendMessage(ctx, req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, StatusFailedToCallAionChat)
		s.logger.ErrorwCtx(ctx, LogFailedToCallAionChat,
			LogKeyError, err.Error(),
			LogKeyUserID, userID,
		)
		return nil, err
	}

	result := &domain.ChatResult{
		Response:      resp.Response,
		Sources:       convertSources(resp.Sources),
		TokensUsed:    resp.TokensUsed,
		FunctionCalls: extractFunctionNames(resp.FunctionCalls),
	}

	span.SetAttributes(
		attribute.Int(AttrTokensUsed, resp.TokensUsed),
		attribute.Int(AttrFunctionCallsCount, len(resp.FunctionCalls)),
	)
	span.SetStatus(codes.Ok, StatusMessageProcessedSuccessfully)

	s.logger.InfowCtx(ctx, LogChatMessageProcessedSuccessfully,
		LogKeyUserID, userID,
		LogKeyTokensUsed, resp.TokensUsed,
		LogKeyResponseLength, len(resp.Response),
	)

	// Save chat history asynchronously (don't block on errors)
	go func() {
		functionCallsMap := make(map[string]string)
		for _, call := range resp.FunctionCalls {
			// Convert Args map to JSON string for storage
			argsJSON := ""
			if len(call.Args) > 0 {
				if jsonBytes, err := json.Marshal(call.Args); err == nil {
					argsJSON = string(jsonBytes)
				}
			}
			functionCallsMap[call.Name] = argsJSON
		}

		if err := s.SaveChatHistory(context.Background(), userID, message, resp.Response, resp.TokensUsed, functionCallsMap); err != nil {
			s.logger.ErrorwCtx(context.Background(), "Failed to save chat history (non-blocking)",
				LogKeyError, err.Error(),
				LogKeyUserID, userID,
			)
		}
	}()

	return result, nil
}

// convertSources converts the sources from the internal response to domain format.
func convertSources(sources []map[string]interface{}) []interface{} {
	if sources == nil {
		return nil
	}
	result := make([]interface{}, len(sources))
	for i, source := range sources {
		result[i] = source
	}
	return result
}

// extractFunctionNames extracts function names from function calls.
func extractFunctionNames(calls []dto.FunctionCall) []string {
	if calls == nil {
		return nil
	}
	names := make([]string, len(calls))
	for i, call := range calls {
		names[i] = call.Name
	}
	return names
}
