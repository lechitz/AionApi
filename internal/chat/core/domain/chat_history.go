// Package domain contains core business entities for the Chat context.
package domain

import (
	"time"

	"github.com/google/uuid"
)

// ChatHistory represents a conversation entry between a user and the AI assistant.
// It stores both the user's message and the AI's response for future reference.
type ChatHistory struct {
	CreatedAt       time.Time         `json:"createdAt"`
	UpdatedAt       time.Time         `json:"updatedAt"`
	FunctionCalls   map[string]string `json:"functionCalls,omitempty"`
	DeletedAt       *time.Time        `json:"deletedAt,omitempty"`
	Message         string            `json:"message"`
	Response        string            `json:"response"`
	ChatID          uint64            `json:"chatId"`
	UserID          uint64            `json:"userId"`
	TokensUsed      int               `json:"tokensUsed"`
	SessionID       *uuid.UUID        `json:"sessionId,omitempty"`       // Groups related conversations together
	ExecutionTimeMs *int              `json:"executionTimeMs,omitempty"` // Query execution time in milliseconds
	ToolCount       int               `json:"toolCount"`                 // Number of GraphQL tools called during this interaction
	ErrorCount      int               `json:"errorCount"`                // Number of errors during execution
}
