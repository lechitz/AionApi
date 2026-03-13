package usecase

import "errors"

const (
	TracerName = "aionapi.eventoutbox.usecase"

	SpanEnqueue = "eventoutbox.enqueue"

	EventNormalizeInput = "eventoutbox.input.normalize"
	EventValidateInput  = "eventoutbox.input.validate"
	EventRepositorySave = "eventoutbox.repository.save"
	EventSuccess        = "eventoutbox.success"

	DefaultEventVersion = "v1"
	EventVersionV1      = DefaultEventVersion
	DefaultEventStatus  = "pending"
	DefaultSource       = "aionapi"

	StatusEventQueued = "event queued"

	LogKeyEventID       = "event_id"
	LogKeyAggregateType = "aggregate_type"
	LogKeyAggregateID   = "aggregate_id"
	LogKeyEventType     = "event_type"
	LogKeyEventVersion  = "event_version"

	LogEnqueueingEvent      = "enqueueing outbox event"
	LogOutboxEventQueued    = "outbox event queued"
	LogFailedToEnqueueEvent = "failed to enqueue outbox event"

	EventIDRequired       = "event id is required"
	AggregateTypeRequired = "aggregate type is required"
	AggregateIDRequired   = "aggregate id is required"
	EventTypeRequired     = "event type is required"
	EventVersionRequired  = "event version is required"
	SourceRequired        = "source is required"
	PayloadRequired       = "payload is required"
)

var (
	ErrEventIDRequired       = errors.New(EventIDRequired)
	ErrAggregateTypeRequired = errors.New(AggregateTypeRequired)
	ErrAggregateIDRequired   = errors.New(AggregateIDRequired)
	ErrEventTypeRequired     = errors.New(EventTypeRequired)
	ErrEventVersionRequired  = errors.New(EventVersionRequired)
	ErrSourceRequired        = errors.New(SourceRequired)
	ErrPayloadRequired       = errors.New(PayloadRequired)
)
