CREATE TABLE IF NOT EXISTS aion_api.tags
(
    id               SERIAL PRIMARY KEY,
    user_id          INT NOT NULL,
    name             VARCHAR(255) NOT NULL,
    category         VARCHAR(255) NOT NULL,
    description      TEXT,
    creation_date    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at       TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES aion_api.users (id) ON DELETE CASCADE -- Relationship to User
);

CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_tags_updated_at
    BEFORE UPDATE ON aion_api.tags
    FOR EACH ROW
    EXECUTE FUNCTION update_timestamp();
