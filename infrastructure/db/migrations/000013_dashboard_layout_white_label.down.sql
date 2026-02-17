-- Rollback: 000013_dashboard_layout_white_label

DROP TRIGGER IF EXISTS update_dashboard_widgets_updated_at ON aion_api.dashboard_widgets;
DROP TABLE IF EXISTS aion_api.dashboard_widgets;

DROP TRIGGER IF EXISTS update_dashboard_views_updated_at ON aion_api.dashboard_views;
DROP TABLE IF EXISTS aion_api.dashboard_views;
