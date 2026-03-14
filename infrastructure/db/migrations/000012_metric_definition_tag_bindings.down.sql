-- Rollback: 000012_metric_definition_tag_bindings

DROP INDEX IF EXISTS aion_api.idx_metric_definition_tag_bindings_tag;
DROP INDEX IF EXISTS aion_api.idx_metric_definition_tag_bindings_user_metric;
DROP INDEX IF EXISTS aion_api.ux_metric_definition_tag_bindings_unique;
DROP TABLE IF EXISTS aion_api.metric_definition_tag_bindings;
