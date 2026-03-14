-- Migration: 000016_audit_action_events
-- Description: Create durable audit trail table for chat-driven actions

CREATE TABLE IF NOT EXISTS aion_api.audit_action_events (
    id BIGSERIAL PRIMARY KEY,
    event_id UUID NOT NULL UNIQUE,
    timestamp_utc TIMESTAMPTZ NOT NULL,
    user_id BIGINT NOT NULL,
    source VARCHAR(32) NOT NULL,
    trace_id VARCHAR(64) NOT NULL,
    request_id VARCHAR(64),
    ui_action_type VARCHAR(32) NOT NULL,
    draft_id VARCHAR(128) NOT NULL,
    action VARCHAR(64) NOT NULL,
    entity VARCHAR(32) NOT NULL,
    operation VARCHAR(16) NOT NULL,
    status VARCHAR(32) NOT NULL,
    entity_id VARCHAR(64),
    consent_required BOOLEAN NOT NULL DEFAULT FALSE,
    consent_confirmed BOOLEAN NOT NULL DEFAULT FALSE,
    consent_policy_version VARCHAR(32),
    quick_add_contract_version VARCHAR(32),
    quick_add_idempotency_key VARCHAR(192),
    message_code VARCHAR(64),
    payload_redacted JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_audit_action_events_ui_action_type
        CHECK (ui_action_type IN ('draft_accept', 'draft_cancel')),
    CONSTRAINT chk_audit_action_events_operation
        CHECK (operation IN ('create', 'update', 'delete', 'cancel')),
    CONSTRAINT chk_audit_action_events_status
        CHECK (status IN ('success', 'failed', 'canceled', 'needs_input', 'blocked'))
);

CREATE INDEX IF NOT EXISTS idx_audit_events_user_time
    ON aion_api.audit_action_events(user_id, timestamp_utc DESC);

CREATE INDEX IF NOT EXISTS idx_audit_events_trace
    ON aion_api.audit_action_events(trace_id);

CREATE INDEX IF NOT EXISTS idx_audit_events_draft
    ON aion_api.audit_action_events(draft_id);

CREATE INDEX IF NOT EXISTS idx_audit_events_status_time
    ON aion_api.audit_action_events(status, timestamp_utc DESC);

CREATE INDEX IF NOT EXISTS idx_audit_events_failures
    ON aion_api.audit_action_events(timestamp_utc DESC)
    WHERE status IN ('failed', 'blocked', 'needs_input');

COMMENT ON TABLE aion_api.audit_action_events IS 'Immutable audit log for chat-driven actions and outcomes';
COMMENT ON COLUMN aion_api.audit_action_events.event_id IS 'Globally unique event identifier (UUID)';
COMMENT ON COLUMN aion_api.audit_action_events.timestamp_utc IS 'Canonical event occurrence timestamp in UTC';
COMMENT ON COLUMN aion_api.audit_action_events.trace_id IS 'Cross-service correlation id for observability';
COMMENT ON COLUMN aion_api.audit_action_events.draft_id IS 'UI draft action identifier for lifecycle reconstruction';
COMMENT ON COLUMN aion_api.audit_action_events.payload_redacted IS 'Allow-listed metadata payload with sensitive values redacted';
