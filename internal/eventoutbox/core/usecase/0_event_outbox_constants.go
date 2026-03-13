package usecase

import (
	"errors"
	"time"
)

const (
	TracerName = "aionapi.eventoutbox.usecase"

	SpanEnqueue        = "eventoutbox.enqueue"
	SpanPublishPending = "eventoutbox.publish_pending"

	EventNormalizeInput    = "eventoutbox.input.normalize"
	EventValidateInput     = "eventoutbox.input.validate"
	EventRepositorySave    = "eventoutbox.repository.save"
	EventRepositoryList    = "eventoutbox.repository.list_pending"
	EventRepositoryPublish = "eventoutbox.repository.mark_published"
	EventPublish           = "eventoutbox.publish"
	EventSuccess           = "eventoutbox.success"

	DefaultEventVersion   = "v1"
	EventVersionV1        = DefaultEventVersion
	DefaultEventStatus    = "pending"
	DefaultSource         = "aionapi"
	DefaultPublishBackoff = 5 * time.Second

	StatusEventQueued = "event queued"

	LogKeyEventID       = "event_id"
	LogKeyAggregateType = "aggregate_type"
	LogKeyAggregateID   = "aggregate_id"
	LogKeyEventType     = "event_type"
	LogKeyEventVersion  = "event_version"

	LogEnqueueingEvent          = "enqueueing outbox event"
	LogOutboxEventQueued        = "outbox event queued"
	LogFailedToEnqueueEvent     = "failed to enqueue outbox event"
	LogPublishPendingEvents     = "publishing pending outbox events"
	LogOutboxEventPublished     = "outbox event published"
	LogOutboxEventPublishFailed = "failed to publish outbox event"

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
