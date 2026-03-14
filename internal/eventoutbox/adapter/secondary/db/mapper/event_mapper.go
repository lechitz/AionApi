// Package mapper provides conversion functions between domain and database models for the outbox context.
package mapper

import (
	"github.com/lechitz/AionApi/internal/eventoutbox/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/eventoutbox/core/domain"
)

// EventToDB converts a domain outbox event into the DB model.
func EventToDB(event domain.Event) model.EventDB {
	return model.EventDB{
		EventID:        event.EventID,
		AggregateType:  event.AggregateType,
		AggregateID:    event.AggregateID,
		EventType:      event.EventType,
		EventVersion:   event.EventVersion,
		Source:         event.Source,
		TraceID:        event.TraceID,
		RequestID:      event.RequestID,
		Status:         event.Status,
		AttemptCount:   event.AttemptCount,
		AvailableAtUTC: event.AvailableAtUTC,
		PublishedAtUTC: event.PublishedAtUTC,
		LastError:      event.LastError,
		PayloadJSON:    event.PayloadJSON,
		CreatedAt:      event.CreatedAt,
	}
}
