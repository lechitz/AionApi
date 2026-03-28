// Package mapper provides conversion functions between domain and database models for the audit context.
package mapper

import (
	"encoding/json"

	"github.com/lechitz/aion-api/internal/audit/adapter/secondary/db/model"
	"github.com/lechitz/aion-api/internal/audit/core/domain"
)

// AuditActionEventToDB converts a domain audit event into the DB model.
func AuditActionEventToDB(event domain.AuditActionEvent) model.AuditActionEventDB {
	row := model.AuditActionEventDB{
		EventID:                 event.EventID,
		TimestampUTC:            event.TimestampUTC,
		UserID:                  event.UserID,
		Source:                  event.Source,
		TraceID:                 event.TraceID,
		RequestID:               event.RequestID,
		UIActionType:            event.UIActionType,
		DraftID:                 event.DraftID,
		Action:                  event.Action,
		Entity:                  event.Entity,
		Operation:               event.Operation,
		Status:                  event.Status,
		EntityID:                event.EntityID,
		ConsentRequired:         event.ConsentRequired,
		ConsentConfirmed:        event.ConsentConfirmed,
		ConsentPolicyVersion:    event.ConsentPolicyVersion,
		QuickAddContractVersion: event.QuickAddContractVersion,
		QuickAddIdempotencyKey:  event.QuickAddIdempotencyKey,
		MessageCode:             event.MessageCode,
	}

	if len(event.PayloadRedacted) > 0 {
		if payload, err := json.Marshal(event.PayloadRedacted); err == nil {
			row.PayloadRedacted = payload
		}
	}

	return row
}

// AuditActionEventFromDB converts DB row into domain audit event.
func AuditActionEventFromDB(row model.AuditActionEventDB) domain.AuditActionEvent {
	event := domain.AuditActionEvent{
		EventID:                 row.EventID,
		TimestampUTC:            row.TimestampUTC,
		UserID:                  row.UserID,
		Source:                  row.Source,
		TraceID:                 row.TraceID,
		RequestID:               row.RequestID,
		UIActionType:            row.UIActionType,
		DraftID:                 row.DraftID,
		Action:                  row.Action,
		Entity:                  row.Entity,
		Operation:               row.Operation,
		Status:                  row.Status,
		EntityID:                row.EntityID,
		ConsentRequired:         row.ConsentRequired,
		ConsentConfirmed:        row.ConsentConfirmed,
		ConsentPolicyVersion:    row.ConsentPolicyVersion,
		QuickAddContractVersion: row.QuickAddContractVersion,
		QuickAddIdempotencyKey:  row.QuickAddIdempotencyKey,
		MessageCode:             row.MessageCode,
	}

	if len(row.PayloadRedacted) > 0 {
		var payload map[string]interface{}
		if err := json.Unmarshal(row.PayloadRedacted, &payload); err == nil {
			event.PayloadRedacted = payload
		}
	}

	return event
}

// AuditActionEventsFromDB converts DB rows into domain audit events.
func AuditActionEventsFromDB(rows []model.AuditActionEventDB) []domain.AuditActionEvent {
	out := make([]domain.AuditActionEvent, len(rows))
	for i := range rows {
		out[i] = AuditActionEventFromDB(rows[i])
	}
	return out
}
