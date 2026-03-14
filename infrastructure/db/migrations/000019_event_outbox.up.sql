-- Migration: 000019_event_outbox
-- Description: Create durable event outbox table for canonical domain events

CREATE TABLE IF NOT EXISTS aion_api.event_outbox (
    id BIGSERIAL PRIMARY KEY,
    event_id UUID NOT NULL UNIQUE,
    aggregate_type VARCHAR(64) NOT NULL,
    aggregate_id VARCHAR(128) NOT NULL,
    event_type VARCHAR(128) NOT NULL,
    event_version VARCHAR(16) NOT NULL,
    source VARCHAR(32) NOT NULL,
    trace_id VARCHAR(64),
    request_id VARCHAR(64),
    status VARCHAR(24) NOT NULL DEFAULT 'pending',
    attempt_count INTEGER NOT NULL DEFAULT 0,
    available_at_utc TIMESTAMPTZ NOT NULL,
    published_at_utc TIMESTAMPTZ,
    last_error TEXT,
    payload_json JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_event_outbox_status
        CHECK (status IN ('pending', 'publishing', 'published', 'failed'))
);

CREATE INDEX IF NOT EXISTS idx_event_outbox_status_available
    ON aion_api.event_outbox(status, available_at_utc ASC);

CREATE INDEX IF NOT EXISTS idx_event_outbox_aggregate
    ON aion_api.event_outbox(aggregate_type, aggregate_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_event_outbox_trace
    ON aion_api.event_outbox(trace_id);

CREATE INDEX IF NOT EXISTS idx_event_outbox_pending_only
    ON aion_api.event_outbox(available_at_utc ASC)
    WHERE status = 'pending';

COMMENT ON TABLE aion_api.event_outbox IS 'Durable outbox for canonical domain events pending publication to Kafka';
COMMENT ON COLUMN aion_api.event_outbox.event_id IS 'Globally unique event identifier (UUID)';
COMMENT ON COLUMN aion_api.event_outbox.aggregate_type IS 'Canonical aggregate family, e.g. record or category';
COMMENT ON COLUMN aion_api.event_outbox.aggregate_id IS 'Aggregate identifier serialized as string for cross-language consumption';
COMMENT ON COLUMN aion_api.event_outbox.event_type IS 'Semantic event type, e.g. record.created';
COMMENT ON COLUMN aion_api.event_outbox.source IS 'Canonical source system that produced the event';
COMMENT ON COLUMN aion_api.event_outbox.status IS 'Publication lifecycle state managed by outbox publisher';
COMMENT ON COLUMN aion_api.event_outbox.available_at_utc IS 'Earliest timestamp the publisher may attempt delivery';
COMMENT ON COLUMN aion_api.event_outbox.payload_json IS 'Canonical event payload serialized as JSON';
