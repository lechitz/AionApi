// Package output defines interfaces for audit output ports.
package output

import (
	"context"

	"github.com/lechitz/AionApi/internal/audit/core/domain"
)

// AuditActionEventRepository defines persistence contract for immutable audit events.
type AuditActionEventRepository interface {
	// Save stores one audit event; the business operation flow must remain authoritative.
	Save(ctx context.Context, event domain.AuditActionEvent) error

	// List returns audit events using diagnostic filters (trace, draft, user and status).
	List(ctx context.Context, filter domain.AuditActionEventFilter) ([]domain.AuditActionEvent, error)
}
