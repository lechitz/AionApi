-- Rollback: 000003_categories_tags_days

DROP TRIGGER IF EXISTS update_tags_updated_at ON aion_api.tags;
DROP INDEX IF EXISTS aion_api.ux_tags_user_name_ci;
DROP TABLE IF EXISTS aion_api.tags;

DROP TABLE IF EXISTS aion_api.days;

DROP TRIGGER IF EXISTS update_categories_updated_at ON aion_api.categories;
DROP INDEX IF EXISTS aion_api.ux_categories_user_name_ci;
DROP TABLE IF EXISTS aion_api.categories;
