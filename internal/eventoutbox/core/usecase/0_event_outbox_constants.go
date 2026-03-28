// Package usecase implements event outbox validation, enqueueing, and publishing flows.
package usecase

import (
	"errors"
	"time"
)

const (
	// TracerName is the tracer name used by the event outbox use case layer.
	TracerName = "aion-api.eventoutbox.usecase"

	// SpanEnqueue covers input validation and persistence of one outbox event.
	SpanEnqueue = "eventoutbox.enqueue"
	// SpanPublishPending covers one publish loop over pending outbox rows.
	SpanPublishPending = "eventoutbox.publish_pending"

	// EventNormalizeInput records normalization of one enqueue request.
	EventNormalizeInput = "eventoutbox.input.normalize"
	// EventValidateInput records validation of one enqueue request.
	EventValidateInput = "eventoutbox.input.validate"
	// EventRepositorySave records persistence of one outbox row.
	EventRepositorySave = "eventoutbox.repository.save"
	// EventRepositoryList records retrieval of pending outbox rows.
	EventRepositoryList = "eventoutbox.repository.list_pending"
	// EventRepositoryPublish records marking one outbox row as published.
	EventRepositoryPublish = "eventoutbox.repository.mark_published"
	// EventPublish records the attempt to publish one outbox event externally.
	EventPublish = "eventoutbox.publish"
	// EventSuccess records successful completion of one outbox flow.
	EventSuccess = "eventoutbox.success"

	// DefaultEventVersion is the current canonical event version for newly enqueued outbox events.
	DefaultEventVersion = "v1"
	// EventVersionV1 aliases the current canonical event version for compatibility with existing callers.
	EventVersionV1 = DefaultEventVersion
	// DefaultEventStatus is the initial persisted status for newly enqueued outbox rows.
	DefaultEventStatus = "pending"
	// DefaultSource is the canonical source label used for outbox events emitted by aion-api.
	DefaultSource = "aion-api"
	// DefaultPublishBackoff is the retry delay applied between publish attempts for transient failures.
	DefaultPublishBackoff = 5 * time.Second

	// StatusEventQueued is the success status attached after one outbox event is persisted.
	StatusEventQueued = "event queued"

	// LogKeyEventID stores the structured logger field name for the outbox event id.
	LogKeyEventID = "event_id"
	// LogKeyAggregateType stores the structured logger field name for the aggregate type.
	LogKeyAggregateType = "aggregate_type"
	// LogKeyAggregateID stores the structured logger field name for the aggregate id.
	LogKeyAggregateID = "aggregate_id"
	// LogKeyEventType stores the structured logger field name for the canonical event type.
	LogKeyEventType = "event_type"
	// LogKeyEventVersion stores the structured logger field name for the canonical event version.
	LogKeyEventVersion = "event_version"

	// LogEnqueueingEvent is emitted before the use case validates and persists an outbox event.
	LogEnqueueingEvent = "enqueueing outbox event"
	// LogOutboxEventQueued is emitted after one outbox row is stored successfully.
	LogOutboxEventQueued = "outbox event queued"
	// LogFailedToEnqueueEvent is emitted when enqueueing an outbox event fails.
	LogFailedToEnqueueEvent = "failed to enqueue outbox event"
	// LogPublishPendingEvents is emitted when the publisher loop starts processing pending rows.
	LogPublishPendingEvents = "publishing pending outbox events"
	// LogOutboxEventPublished is emitted after one outbox row is published successfully.
	LogOutboxEventPublished = "outbox event published"
	// LogOutboxEventPublishFailed is emitted when publishing an outbox row fails.
	LogOutboxEventPublishFailed = "failed to publish outbox event"

	// EventIDRequired is the validation message returned when the event id is missing.
	EventIDRequired = "event id is required"
	// AggregateTypeRequired is the validation message returned when the aggregate type is missing.
	AggregateTypeRequired = "aggregate type is required"
	// AggregateIDRequired is the validation message returned when the aggregate id is missing.
	AggregateIDRequired = "aggregate id is required"
	// EventTypeRequired is the validation message returned when the event type is missing.
	EventTypeRequired = "event type is required"
	// EventVersionRequired is the validation message returned when the event version is missing.
	EventVersionRequired = "event version is required"
	// SourceRequired is the validation message returned when the event source is missing.
	SourceRequired = "source is required"
	// PayloadRequired is the validation message returned when the event payload is missing.
	PayloadRequired = "payload is required"
)

var (
	// ErrEventIDRequired indicates that the outbox event is missing its id.
	ErrEventIDRequired = errors.New(EventIDRequired)
	// ErrAggregateTypeRequired indicates that the outbox event is missing its aggregate type.
	ErrAggregateTypeRequired = errors.New(AggregateTypeRequired)
	// ErrAggregateIDRequired indicates that the outbox event is missing its aggregate id.
	ErrAggregateIDRequired = errors.New(AggregateIDRequired)
	// ErrEventTypeRequired indicates that the outbox event is missing its type.
	ErrEventTypeRequired = errors.New(EventTypeRequired)
	// ErrEventVersionRequired indicates that the outbox event is missing its version.
	ErrEventVersionRequired = errors.New(EventVersionRequired)
	// ErrSourceRequired indicates that the outbox event is missing its source.
	ErrSourceRequired = errors.New(SourceRequired)
	// ErrPayloadRequired indicates that the outbox event is missing its payload.
	ErrPayloadRequired = errors.New(PayloadRequired)
)
