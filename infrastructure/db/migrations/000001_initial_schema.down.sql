-- Rollback: 000001_initial_schema
-- Warning: This drops the entire schema and all data!

DROP FUNCTION IF EXISTS aion_api.update_timestamp();
DROP SCHEMA IF EXISTS aion_api CASCADE;
