-- Migration: 000010_dashboard_metrics_goals (rollback)

DROP TRIGGER IF EXISTS update_goal_instances_updated_at ON aion_api.goal_instances;
DROP TABLE IF EXISTS aion_api.goal_instances;

DROP TRIGGER IF EXISTS update_goal_templates_updated_at ON aion_api.goal_templates;
DROP TABLE IF EXISTS aion_api.goal_templates;

DROP TRIGGER IF EXISTS update_metric_definitions_updated_at ON aion_api.metric_definitions;
DROP TABLE IF EXISTS aion_api.metric_definitions;
