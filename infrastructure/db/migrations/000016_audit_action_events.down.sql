-- Migration: 000016_audit_action_events (down)
-- Description: Drop durable audit trail table for chat-driven actions

DROP INDEX IF EXISTS aion_api.idx_audit_events_failures;
DROP INDEX IF EXISTS aion_api.idx_audit_events_status_time;
DROP INDEX IF EXISTS aion_api.idx_audit_events_draft;
DROP INDEX IF EXISTS aion_api.idx_audit_events_trace;
DROP INDEX IF EXISTS aion_api.idx_audit_events_user_time;

DROP TABLE IF EXISTS aion_api.audit_action_events;
