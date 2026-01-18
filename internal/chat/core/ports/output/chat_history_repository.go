// Package output defines interfaces for chat history repository operations.
package output

import (
	"context"

	"github.com/lechitz/AionApi/internal/chat/core/domain"
)

// ChatHistoryRepository defines an interface for managing chat history persistence.
type ChatHistoryRepository interface {
	// Save persists a new chat history entry to the database.
	Save(ctx context.Context, chatHistory domain.ChatHistory) (domain.ChatHistory, error)

	// GetLatest retrieves the N most recent chat history entries for a user.
	GetLatest(ctx context.Context, userID uint64, limit int) ([]domain.ChatHistory, error)

	// GetByUserID retrieves chat history entries for a user with pagination.
	GetByUserID(ctx context.Context, userID uint64, limit, offset int) ([]domain.ChatHistory, error)
}
