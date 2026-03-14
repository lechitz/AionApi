// Package domain contains chat domain models and value objects.
package domain

import "time"

// AuditActionEvent represents one immutable audit event for a chat-driven action.
type AuditActionEvent struct {
	EventID                 string
	TimestampUTC            time.Time
	UserID                  uint64
	Source                  string
	TraceID                 string
	RequestID               string
	UIActionType            string
	DraftID                 string
	Action                  string
	Entity                  string
	Operation               string
	Status                  string
	EntityID                string
	ConsentRequired         bool
	ConsentConfirmed        bool
	ConsentPolicyVersion    string
	QuickAddContractVersion string
	QuickAddIdempotencyKey  string
	MessageCode             string
	PayloadRedacted         map[string]interface{}
}

// AuditActionEventFilter defines query fields for internal diagnostics endpoints.
type AuditActionEventFilter struct {
	UserID   *uint64
	TraceID  string
	DraftID  string
	Statuses []string
	FromUTC  *time.Time
	ToUTC    *time.Time
	Limit    int
	Offset   int
}
