// Package controller provides mapping helpers between GraphQL models and core commands/domain for the Chat context.
package controller

import (
	"encoding/json"
	"strconv"

	gmodel "github.com/lechitz/aion-api/internal/adapter/primary/graphql/model"
	"github.com/lechitz/aion-api/internal/chat/core/domain"
)

// safeInt32 safely converts int to int32 with range validation.
func safeInt32(v int) int32 {
	if v >= -2147483648 && v <= 2147483647 {
		return int32(v) // #nosec G115 - validated range
	}
	return 0 // fallback for overflow
}

// toModelOut converts a domain.ChatHistory to a GraphQL model.ChatMessage.
func toModelOut(ch domain.ChatHistory) *gmodel.ChatMessage {
	out := &gmodel.ChatMessage{
		ID:         strconv.FormatUint(ch.ChatID, 10),
		UserID:     strconv.FormatUint(ch.UserID, 10),
		Message:    ch.Message,
		Response:   ch.Response,
		TokensUsed: safeInt32(ch.TokensUsed),
		CreatedAt:  ch.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  ch.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Convert FunctionCalls map to JSON (if exists)
	if len(ch.FunctionCalls) > 0 {
		if jsonBytes, err := json.Marshal(ch.FunctionCalls); err == nil {
			jsonStr := string(jsonBytes)
			out.FunctionCalls = &jsonStr
		}
	}

	return out
}

// toModelOutSlice converts a slice of domain.ChatHistory to a slice of GraphQL model.ChatMessage pointers.
func toModelOutSlice(histories []domain.ChatHistory) []*gmodel.ChatMessage {
	result := make([]*gmodel.ChatMessage, len(histories))
	for i, ch := range histories {
		result[i] = toModelOut(ch)
	}
	return result
}
