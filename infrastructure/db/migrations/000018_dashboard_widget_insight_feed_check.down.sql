ALTER TABLE aion_api.dashboard_widgets
    DROP CONSTRAINT IF EXISTS chk_dashboard_widgets_type;

ALTER TABLE aion_api.dashboard_widgets
    ADD CONSTRAINT chk_dashboard_widgets_type
    CHECK (widget_type IN ('kpi_number', 'goal_progress', 'trend_line', 'checklist'));
