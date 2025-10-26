CREATE SCHEMA IF NOT EXISTS aion_api;

SET search_path TO aion_api;

CREATE OR REPLACE FUNCTION aion_api.update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;