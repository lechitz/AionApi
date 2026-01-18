// Package domain contains core business entities for the Chat context.
package domain

import "time"

// ChatHistory represents a conversation entry between a user and the AI assistant.
// It stores both the user's message and the AI's response for future reference.
type ChatHistory struct {
	CreatedAt     time.Time
	UpdatedAt     time.Time
	FunctionCalls map[string]string
	DeletedAt     *time.Time
	Message       string
	Response      string
	ChatID        uint64
	UserID        uint64
	TokensUsed    int
}
