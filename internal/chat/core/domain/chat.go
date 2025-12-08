// Package domain contains the chat domain models and business logic.
package domain

import "time"

// ChatMessage represents a chat message in the domain.
type ChatMessage struct {
	Timestamp  time.Time
	ID         string
	Message    string
	Response   string
	UserID     uint64
	TokensUsed int
}

// ChatResult represents the result of processing a chat message.
type ChatResult struct {
	Response      string
	Sources       []interface{}
	FunctionCalls []string
	TokensUsed    int
}
