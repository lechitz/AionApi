-- Migration: 000012_metric_definition_tag_bindings
-- Description: Allow dashboard metrics to aggregate multiple tags

CREATE TABLE IF NOT EXISTS aion_api.metric_definition_tag_bindings (
    id                   BIGSERIAL PRIMARY KEY,
    user_id              BIGINT NOT NULL REFERENCES aion_api.users (user_id) ON DELETE CASCADE,
    metric_definition_id BIGINT NOT NULL REFERENCES aion_api.metric_definitions (id) ON DELETE CASCADE,
    tag_id               BIGINT NOT NULL REFERENCES aion_api.tags (tag_id) ON DELETE CASCADE,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_metric_definition_tag_bindings_unique
    ON aion_api.metric_definition_tag_bindings (metric_definition_id, tag_id);

CREATE INDEX IF NOT EXISTS idx_metric_definition_tag_bindings_user_metric
    ON aion_api.metric_definition_tag_bindings (user_id, metric_definition_id);

CREATE INDEX IF NOT EXISTS idx_metric_definition_tag_bindings_tag
    ON aion_api.metric_definition_tag_bindings (tag_id);
