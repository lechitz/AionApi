// Package model contains database models for the audit context.
package model

import "time"

const (
	// AuditActionEventsTable is the fully qualified table name for audit action events.
	AuditActionEventsTable = "aion_api.audit_action_events"
)

// AuditActionEventDB represents one row in the immutable audit action events table.
type AuditActionEventDB struct {
	ID                      uint64    `gorm:"primaryKey;column:id;autoIncrement"`
	EventID                 string    `gorm:"column:event_id;type:uuid;not null;uniqueIndex"`
	TimestampUTC            time.Time `gorm:"column:timestamp_utc;not null"`
	UserID                  uint64    `gorm:"column:user_id;not null;index"`
	Source                  string    `gorm:"column:source;size:32;not null"`
	TraceID                 string    `gorm:"column:trace_id;size:64;not null;index"`
	RequestID               string    `gorm:"column:request_id;size:64"`
	UIActionType            string    `gorm:"column:ui_action_type;size:32;not null"`
	DraftID                 string    `gorm:"column:draft_id;size:128;not null;index"`
	Action                  string    `gorm:"column:action;size:64;not null"`
	Entity                  string    `gorm:"column:entity;size:32;not null"`
	Operation               string    `gorm:"column:operation;size:16;not null"`
	Status                  string    `gorm:"column:status;size:32;not null;index"`
	EntityID                string    `gorm:"column:entity_id;size:64"`
	ConsentRequired         bool      `gorm:"column:consent_required;not null;default:false"`
	ConsentConfirmed        bool      `gorm:"column:consent_confirmed;not null;default:false"`
	ConsentPolicyVersion    string    `gorm:"column:consent_policy_version;size:32"`
	QuickAddContractVersion string    `gorm:"column:quick_add_contract_version;size:32"`
	QuickAddIdempotencyKey  string    `gorm:"column:quick_add_idempotency_key;size:192"`
	MessageCode             string    `gorm:"column:message_code;size:64"`
	PayloadRedacted         []byte    `gorm:"column:payload_redacted;type:jsonb"`
	CreatedAt               time.Time `gorm:"column:created_at;autoCreateTime"`
}

// TableName returns the fully qualified database table name for AuditActionEventDB.
func (AuditActionEventDB) TableName() string {
	return AuditActionEventsTable
}
