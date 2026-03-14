-- Migration: 000005_records
-- Description: Create records table with full-text search
-- This consolidates: 11_records.sql

CREATE TABLE IF NOT EXISTS aion_api.records (
    id                BIGSERIAL PRIMARY KEY,
    user_id           BIGINT NOT NULL REFERENCES aion_api.users (user_id) ON DELETE CASCADE,
    tag_id            BIGINT NOT NULL REFERENCES aion_api.tags (tag_id) ON DELETE RESTRICT,
    description       TEXT,
    value             NUMERIC(12,2),
    duration_seconds  INTEGER,
    event_time        TIMESTAMPTZ NOT NULL,
    recorded_at       TIMESTAMPTZ,
    source            VARCHAR(50),
    timezone          VARCHAR(64),
    status            VARCHAR(32) DEFAULT 'published',
    search_vector     TSVECTOR,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ DEFAULT NULL
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_records_user_event_time
    ON aion_api.records (user_id, event_time DESC);
CREATE INDEX IF NOT EXISTS idx_records_user_recorded_at
    ON aion_api.records (user_id, recorded_at DESC);
CREATE INDEX IF NOT EXISTS idx_records_tag
    ON aion_api.records (tag_id);
CREATE INDEX IF NOT EXISTS idx_records_user_active_event_time
    ON aion_api.records (user_id, deleted_at, event_time DESC);
CREATE INDEX IF NOT EXISTS idx_records_search_vector
    ON aion_api.records USING GIN (search_vector);

-- Trigger to auto-update updated_at timestamp
DROP TRIGGER IF EXISTS update_records_updated_at ON aion_api.records;
CREATE TRIGGER update_records_updated_at
    BEFORE UPDATE ON aion_api.records
    FOR EACH ROW
    EXECUTE FUNCTION aion_api.update_timestamp();

-- Trigger to auto-update search_vector from description (Portuguese config)
CREATE OR REPLACE FUNCTION aion_api.update_records_search_vector()
RETURNS TRIGGER AS $$
BEGIN
    NEW.search_vector := to_tsvector('portuguese', COALESCE(NEW.description, ''));
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_records_search_vector ON aion_api.records;
CREATE TRIGGER update_records_search_vector
    BEFORE INSERT OR UPDATE OF description ON aion_api.records
    FOR EACH ROW
    EXECUTE FUNCTION aion_api.update_records_search_vector();

-- Table and column comments
COMMENT ON TABLE aion_api.records IS
    'User activity records (water intake, exercise, meals, etc.)';
COMMENT ON COLUMN aion_api.records.id IS
    'Primary key - unique record identifier';
COMMENT ON COLUMN aion_api.records.tag_id IS
    'Tag classifier (category obtained via JOIN with tags table)';
COMMENT ON COLUMN aion_api.records.description IS
    'Free-text description of the activity';
COMMENT ON COLUMN aion_api.records.value IS
    'Numeric value (e.g., liters of water, duration in minutes)';
COMMENT ON COLUMN aion_api.records.duration_seconds IS
    'Duration of activity in seconds';
COMMENT ON COLUMN aion_api.records.event_time IS
    'When the event actually occurred';
COMMENT ON COLUMN aion_api.records.recorded_at IS
    'When the user recorded this event (may differ from event_time)';
COMMENT ON COLUMN aion_api.records.search_vector IS
    'Full-text search vector (Portuguese) - auto-generated from description';
