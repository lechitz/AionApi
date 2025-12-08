// Package output defines the outbound ports (interfaces to external services).
package output

import (
	"context"

	"github.com/lechitz/AionApi/internal/chat/adapter/primary/http/dto"
)

// AionChatClient defines the interface for communicating with the Aion-Chat service (Python).
type AionChatClient interface {
	SendMessage(ctx context.Context, req *dto.InternalChatRequest) (*dto.InternalChatResponse, error)
}
