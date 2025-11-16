-- migration: create records table (single tag per record)
-- filepath: infrastructure/db/migrations/11_records.sql

CREATE TABLE IF NOT EXISTS aion_api.records (
    id                BIGSERIAL PRIMARY KEY,
    user_id           BIGINT NOT NULL REFERENCES aion_api.users (user_id) ON DELETE CASCADE,
    category_id       BIGINT NOT NULL REFERENCES aion_api.tag_categories (category_id) ON DELETE RESTRICT,
    tag_id            BIGINT NOT NULL REFERENCES aion_api.tags (tag_id) ON DELETE RESTRICT,
    title             VARCHAR(255) NOT NULL,
    description       TEXT,
    value             NUMERIC(12,2),
    duration_seconds  INTEGER,
    event_time        TIMESTAMPTZ NOT NULL,
    recorded_at       TIMESTAMPTZ,
    source            VARCHAR(50),
    timezone          VARCHAR(64),
    status            VARCHAR(32) DEFAULT 'published',

    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ DEFAULT NULL
);

-- indices
CREATE INDEX IF NOT EXISTS idx_records_user_event_time ON aion_api.records (user_id, event_time DESC);
CREATE INDEX IF NOT EXISTS idx_records_user_recorded_at ON aion_api.records (user_id, recorded_at DESC);
CREATE INDEX IF NOT EXISTS idx_records_category ON aion_api.records (category_id);
CREATE INDEX IF NOT EXISTS idx_records_tag ON aion_api.records (tag_id);

-- trigger to update updated_at on update (assumes function aion_api.update_timestamp() exists)
DROP TRIGGER IF EXISTS update_records_updated_at ON aion_api.records;
CREATE TRIGGER update_records_updated_at
    BEFORE UPDATE ON aion_api.records
    FOR EACH ROW
    EXECUTE FUNCTION aion_api.update_timestamp();
