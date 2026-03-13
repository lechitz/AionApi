DROP INDEX IF EXISTS aion_api.idx_event_outbox_pending_only;
DROP INDEX IF EXISTS aion_api.idx_event_outbox_trace;
DROP INDEX IF EXISTS aion_api.idx_event_outbox_aggregate;
DROP INDEX IF EXISTS aion_api.idx_event_outbox_status_available;

DROP TABLE IF EXISTS aion_api.event_outbox;
