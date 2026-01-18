-- Rollback: 000006_chat_history

DROP INDEX IF EXISTS aion_api.idx_chat_history_user_created;
DROP INDEX IF EXISTS aion_api.idx_chat_history_created_at;
DROP INDEX IF EXISTS aion_api.idx_chat_history_user_id;
DROP TABLE IF EXISTS aion_api.chat_history;
