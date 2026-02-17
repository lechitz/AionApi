-- Migration: 000013_dashboard_layout_white_label
-- Description: Persist user-defined dashboard views/widgets for white-label dashboard layout

CREATE TABLE IF NOT EXISTS aion_api.dashboard_views (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES aion_api.users (user_id) ON DELETE CASCADE,
    name        VARCHAR(120) NOT NULL,
    is_default  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_dashboard_views_user
    ON aion_api.dashboard_views (user_id);

CREATE INDEX IF NOT EXISTS idx_dashboard_views_user_default
    ON aion_api.dashboard_views (user_id, is_default);

CREATE UNIQUE INDEX IF NOT EXISTS ux_dashboard_views_user_name
    ON aion_api.dashboard_views (user_id, lower(name));

DROP TRIGGER IF EXISTS update_dashboard_views_updated_at ON aion_api.dashboard_views;
CREATE TRIGGER update_dashboard_views_updated_at
    BEFORE UPDATE ON aion_api.dashboard_views
    FOR EACH ROW
    EXECUTE FUNCTION aion_api.update_timestamp();

CREATE TABLE IF NOT EXISTS aion_api.dashboard_widgets (
    id                      BIGSERIAL PRIMARY KEY,
    user_id                 BIGINT NOT NULL REFERENCES aion_api.users (user_id) ON DELETE CASCADE,
    view_id                 BIGINT NOT NULL REFERENCES aion_api.dashboard_views (id) ON DELETE CASCADE,
    metric_definition_id    BIGINT NOT NULL REFERENCES aion_api.metric_definitions (id) ON DELETE CASCADE,
    widget_type             VARCHAR(32) NOT NULL,
    size                    VARCHAR(16) NOT NULL,
    order_index             INTEGER NOT NULL DEFAULT 0,
    title_override          VARCHAR(120),
    config_json             JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_active               BOOLEAN NOT NULL DEFAULT TRUE,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_dashboard_widgets_size
        CHECK (size IN ('small', 'medium', 'large')),
    CONSTRAINT chk_dashboard_widgets_type
        CHECK (widget_type IN ('kpi_number', 'goal_progress', 'trend_line', 'checklist'))
);

CREATE INDEX IF NOT EXISTS idx_dashboard_widgets_user_view_active
    ON aion_api.dashboard_widgets (user_id, view_id, is_active);

CREATE INDEX IF NOT EXISTS idx_dashboard_widgets_user_metric
    ON aion_api.dashboard_widgets (user_id, metric_definition_id);

DROP TRIGGER IF EXISTS update_dashboard_widgets_updated_at ON aion_api.dashboard_widgets;
CREATE TRIGGER update_dashboard_widgets_updated_at
    BEFORE UPDATE ON aion_api.dashboard_widgets
    FOR EACH ROW
    EXECUTE FUNCTION aion_api.update_timestamp();
