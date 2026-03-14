-- Rollback: 000014_registration_sessions

DROP INDEX IF EXISTS aion_api.idx_registration_sessions_status_expires_at;
DROP INDEX IF EXISTS aion_api.idx_registration_sessions_email;
DROP INDEX IF EXISTS aion_api.idx_registration_sessions_username;
DROP TABLE IF EXISTS aion_api.registration_sessions;

