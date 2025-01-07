CREATE TABLE IF NOT EXISTS aion_api.day_intentions
(
    id            SERIAL PRIMARY KEY,
    day_id        INT NOT NULL,
    intention     TEXT NOT NULL,
    is_completed  BOOLEAN DEFAULT FALSE,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (day_id) REFERENCES aion_api.days (id) ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_day_intentions_updated_at
    BEFORE UPDATE ON aion_api.day_intentions
    FOR EACH ROW
    EXECUTE FUNCTION update_timestamp();