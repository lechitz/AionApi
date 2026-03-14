-- Rollback: 000005_records

DROP TRIGGER IF EXISTS update_records_search_vector ON aion_api.records;
DROP TRIGGER IF EXISTS update_records_updated_at ON aion_api.records;
DROP FUNCTION IF EXISTS aion_api.update_records_search_vector();
DROP INDEX IF EXISTS aion_api.idx_records_search_vector;
DROP INDEX IF EXISTS aion_api.idx_records_user_active_event_time;
DROP INDEX IF EXISTS aion_api.idx_records_tag;
DROP INDEX IF EXISTS aion_api.idx_records_user_recorded_at;
DROP INDEX IF EXISTS aion_api.idx_records_user_event_time;
DROP TABLE IF EXISTS aion_api.records;
