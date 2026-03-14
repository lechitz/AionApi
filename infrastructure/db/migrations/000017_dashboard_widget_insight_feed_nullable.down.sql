DELETE FROM aion_api.dashboard_widgets
WHERE metric_definition_id IS NULL;

ALTER TABLE aion_api.dashboard_widgets
    ALTER COLUMN metric_definition_id SET NOT NULL;
