-- Rollback: 000004_diaries_and_day_data

DROP TRIGGER IF EXISTS update_day_intentions_updated_at ON aion_api.day_intentions;
DROP TABLE IF EXISTS aion_api.day_intentions;
DROP TABLE IF EXISTS aion_api.day_water_intake;
DROP TABLE IF EXISTS aion_api.day_energy;
DROP TABLE IF EXISTS aion_api.day_moods;
DROP TABLE IF EXISTS aion_api.day_tag_summary;
DROP TABLE IF EXISTS aion_api.professional_diaries;
DROP TABLE IF EXISTS aion_api.personal_diaries;
