// Package output defines interfaces for chat history cache operations.
package output

import (
	"context"

	"github.com/lechitz/AionApi/internal/chat/core/domain"
)

// ChatHistoryCache defines an interface for managing chat history in cache (Redis).
// This is used for fast access to recent conversation history for LLM context.
type ChatHistoryCache interface {
	// GetLatest retrieves the N most recent chat history entries from cache.
	GetLatest(ctx context.Context, userID uint64, limit int) ([]domain.ChatHistory, error)

	// Add adds a new chat history entry to cache.
	Add(ctx context.Context, userID uint64, history domain.ChatHistory) error

	// Clear removes chat history cache for a user.
	Clear(ctx context.Context, userID uint64) error
}
