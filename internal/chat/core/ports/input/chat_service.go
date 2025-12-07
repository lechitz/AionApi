// Package input defines the inbound ports (use cases) for the chat module.
package input

import (
	"context"

	"github.com/lechitz/AionApi/internal/chat/core/domain"
)

// ChatService defines the interface for chat operations (use cases).
type ChatService interface {
	ProcessMessage(ctx context.Context, userID uint64, message string) (*domain.ChatResult, error)
}
