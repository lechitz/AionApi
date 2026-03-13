// Package mapper provides conversion functions between domain and database models for the outbox context.
package mapper

import (
	"github.com/lechitz/AionApi/internal/eventoutbox/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/eventoutbox/core/domain"
)

// EventFromDB converts one DB row into a domain outbox event.
func EventFromDB(row model.EventDB) domain.Event {
	return domain.Event{
		EventID:        row.EventID,
		AggregateType:  row.AggregateType,
		AggregateID:    row.AggregateID,
		EventType:      row.EventType,
		EventVersion:   row.EventVersion,
		Source:         row.Source,
		TraceID:        row.TraceID,
		RequestID:      row.RequestID,
		Status:         row.Status,
		AttemptCount:   row.AttemptCount,
		AvailableAtUTC: row.AvailableAtUTC,
		PublishedAtUTC: row.PublishedAtUTC,
		LastError:      row.LastError,
		PayloadJSON:    row.PayloadJSON,
		CreatedAt:      row.CreatedAt,
	}
}

// EventsFromDB converts DB rows into domain outbox events.
func EventsFromDB(rows []model.EventDB) []domain.Event {
	out := make([]domain.Event, len(rows))
	for i := range rows {
		out[i] = EventFromDB(rows[i])
	}
	return out
}
