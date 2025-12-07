// Package domain contains the chat domain models and business logic.
package domain

import "time"

// ChatMessage represents a chat message in the domain.
type ChatMessage struct {
	ID         string
	UserID     uint64
	Message    string
	Response   string
	Timestamp  time.Time
	TokensUsed int
}

// ChatResult represents the result of processing a chat message.
type ChatResult struct {
	Response      string
	Sources       []interface{}
	TokensUsed    int
	FunctionCalls []string
}
