package repository

const (
	// OutboxTracerName is the tracer name used by the outbox repository.
	OutboxTracerName = "aionapi.eventoutbox.repository"
)

const (
	// SpanOutboxSaveRepo is the span name for persisting one outbox event.
	SpanOutboxSaveRepo = "eventoutbox.repository.save"
)

const (
	// OpOutboxSave is the operation value for outbox save.
	OpOutboxSave = "event_outbox_save"
)

const (
	// StatusOutboxSaved indicates successful outbox persistence.
	StatusOutboxSaved = "event outbox row saved successfully"
)

const (
	// ErrSaveOutboxEventMsg is used when saving an outbox event fails.
	ErrSaveOutboxEventMsg = "error saving outbox event"
)
