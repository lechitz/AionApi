-- Migration: 000010_dashboard_metrics_goals
-- Description: Add metric definitions and daily goals tables for dashboard projections

CREATE TABLE IF NOT EXISTS aion_api.metric_definitions (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES aion_api.users (user_id) ON DELETE CASCADE,
    metric_key      VARCHAR(64) NOT NULL,
    display_name    VARCHAR(120) NOT NULL,
    category_id     BIGINT NULL REFERENCES aion_api.categories (category_id) ON DELETE SET NULL,
    tag_id          BIGINT NOT NULL REFERENCES aion_api.tags (tag_id) ON DELETE RESTRICT,
    value_source    VARCHAR(32) NOT NULL, -- value | duration_seconds | count | latest_value
    aggregation     VARCHAR(32) NOT NULL, -- sum | avg | latest | count
    unit            VARCHAR(24) NOT NULL DEFAULT 'count',
    goal_default    NUMERIC(12,2),
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_metric_definitions_user_metric_key
    ON aion_api.metric_definitions (user_id, metric_key);

CREATE INDEX IF NOT EXISTS idx_metric_definitions_user_active
    ON aion_api.metric_definitions (user_id, is_active);

DROP TRIGGER IF EXISTS update_metric_definitions_updated_at ON aion_api.metric_definitions;
CREATE TRIGGER update_metric_definitions_updated_at
    BEFORE UPDATE ON aion_api.metric_definitions
    FOR EACH ROW
    EXECUTE FUNCTION aion_api.update_timestamp();

CREATE TABLE IF NOT EXISTS aion_api.goal_templates (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES aion_api.users (user_id) ON DELETE CASCADE,
    metric_key      VARCHAR(64) NOT NULL,
    title           VARCHAR(200) NOT NULL,
    target_value    NUMERIC(12,2) NOT NULL,
    comparison      VARCHAR(16) NOT NULL DEFAULT 'gte', -- gte | lte | eq
    period          VARCHAR(16) NOT NULL DEFAULT 'day',
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_goal_templates_user_active
    ON aion_api.goal_templates (user_id, is_active);

DROP TRIGGER IF EXISTS update_goal_templates_updated_at ON aion_api.goal_templates;
CREATE TRIGGER update_goal_templates_updated_at
    BEFORE UPDATE ON aion_api.goal_templates
    FOR EACH ROW
    EXECUTE FUNCTION aion_api.update_timestamp();

CREATE TABLE IF NOT EXISTS aion_api.goal_instances (
    id                  BIGSERIAL PRIMARY KEY,
    user_id             BIGINT NOT NULL REFERENCES aion_api.users (user_id) ON DELETE CASCADE,
    goal_template_id    BIGINT NOT NULL REFERENCES aion_api.goal_templates (id) ON DELETE CASCADE,
    date                DATE NOT NULL,
    target_value        NUMERIC(12,2) NOT NULL,
    current_value       NUMERIC(12,2) NOT NULL DEFAULT 0,
    status              VARCHAR(16) NOT NULL DEFAULT 'pending', -- pending | completed | failed
    progress_pct        NUMERIC(5,2) NOT NULL DEFAULT 0,
    computed_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_goal_instances_user_template_day
    ON aion_api.goal_instances (user_id, goal_template_id, date);

CREATE INDEX IF NOT EXISTS idx_goal_instances_user_date
    ON aion_api.goal_instances (user_id, date);

DROP TRIGGER IF EXISTS update_goal_instances_updated_at ON aion_api.goal_instances;
CREATE TRIGGER update_goal_instances_updated_at
    BEFORE UPDATE ON aion_api.goal_instances
    FOR EACH ROW
    EXECUTE FUNCTION aion_api.update_timestamp();
