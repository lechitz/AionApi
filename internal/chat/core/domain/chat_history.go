// Package domain contains core business entities for the Chat context.
package domain

import "time"

// ChatHistory represents a conversation entry between a user and the AI assistant.
// It stores both the user's message and the AI's response for future reference.
type ChatHistory struct {
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt"`
	FunctionCalls map[string]string `json:"functionCalls,omitempty"`
	DeletedAt     *time.Time        `json:"deletedAt,omitempty"`
	Message       string            `json:"message"`
	Response      string            `json:"response"`
	ChatID        uint64            `json:"chatId"`
	UserID        uint64            `json:"userId"`
	TokensUsed    int               `json:"tokensUsed"`
}
