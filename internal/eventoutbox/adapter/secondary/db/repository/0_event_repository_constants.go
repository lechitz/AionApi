// Package repository persists and queries event outbox rows in PostgreSQL.
package repository

const (
	// OutboxTracerName is the tracer name used by the outbox repository.
	OutboxTracerName = "aionapi.eventoutbox.repository"
)

const (
	// SpanOutboxSaveRepo is the span name for persisting one outbox event.
	SpanOutboxSaveRepo = "eventoutbox.repository.save"
	// SpanOutboxListPendingRepo is the span name for listing pending outbox events.
	SpanOutboxListPendingRepo = "eventoutbox.repository.list_pending"
	// SpanOutboxMarkPublishedRepo is the span name for marking one outbox event as published.
	SpanOutboxMarkPublishedRepo = "eventoutbox.repository.mark_published"
	// SpanOutboxRescheduleRepo is the span name for rescheduling one outbox event.
	SpanOutboxRescheduleRepo = "eventoutbox.repository.reschedule"
)

const (
	// OpOutboxSave is the operation value for outbox save.
	OpOutboxSave = "event_outbox_save"
	// OpOutboxListPending is the operation value for pending event lookup.
	OpOutboxListPending = "event_outbox_list_pending"
	// OpOutboxMarkPublished is the operation value for marking a row as published.
	OpOutboxMarkPublished = "event_outbox_mark_published"
	// OpOutboxReschedule is the operation value for deferring one event.
	OpOutboxReschedule = "event_outbox_reschedule"
)

const (
	// StatusOutboxSaved indicates successful outbox persistence.
	StatusOutboxSaved = "event outbox row saved successfully"
	// StatusOutboxListed indicates successful listing of outbox rows.
	StatusOutboxListed = "event outbox rows listed successfully"
	// StatusOutboxPublished indicates one outbox row was marked as published.
	StatusOutboxPublished = "event outbox row marked as published"
	// StatusOutboxRescheduled indicates one outbox row was deferred for retry.
	StatusOutboxRescheduled = "event outbox row rescheduled"
)

const (
	// ErrSaveOutboxEventMsg is used when saving an outbox event fails.
	ErrSaveOutboxEventMsg = "error saving outbox event"
	// ErrListOutboxEventsMsg is used when listing outbox events fails.
	ErrListOutboxEventsMsg = "error listing outbox events"
	// ErrMarkOutboxPublishedMsg is used when marking an event as published fails.
	ErrMarkOutboxPublishedMsg = "error marking outbox event as published"
	// ErrRescheduleOutboxEventMsg is used when rescheduling an event fails.
	ErrRescheduleOutboxEventMsg = "error rescheduling outbox event"
)
