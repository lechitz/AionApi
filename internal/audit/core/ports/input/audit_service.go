// Package input defines use case interfaces for the audit context.
package input

import (
	"context"

	"github.com/lechitz/AionApi/internal/audit/core/domain"
)

// Service defines audit use cases for writing and querying immutable action events.
type Service interface {
	// WriteEvent validates and persists one audit event.
	WriteEvent(ctx context.Context, event domain.AuditActionEvent) error

	// ListEvents returns audit events for internal diagnostics according to filters.
	ListEvents(ctx context.Context, filter domain.AuditActionEventFilter) ([]domain.AuditActionEvent, error)
}
