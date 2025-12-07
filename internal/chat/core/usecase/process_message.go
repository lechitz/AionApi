// Package usecase implements the chat use cases (business logic).
package usecase

import (
	"context"

	"github.com/lechitz/AionApi/internal/chat/adapter/primary/http/dto"
	"github.com/lechitz/AionApi/internal/chat/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ProcessMessage processes a chat message by forwarding it to the Aion-Chat service.
func (s *ChatService) ProcessMessage(ctx context.Context, userID uint64, message string) (*domain.ChatResult, error) {
	tr := otel.Tracer(TracerChatService)
	ctx, span := tr.Start(ctx, SpanProcessMessage)
	defer span.End()

	span.SetAttributes(
		attribute.Int64(AttrUserID, int64(userID)),
		attribute.String(AttrMessageLength, string(rune(len(message)))),
	)

	s.logger.InfowCtx(ctx, LogProcessingChatMessage, LogKeyUserID, userID, LogKeyMessageLength, len(message))

	req := &dto.InternalChatRequest{
		UserID:  userID,
		Message: message,
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
