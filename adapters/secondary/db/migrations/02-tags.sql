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
    FOREIGN KEY (user_id) REFERENCES aion_api.users (id) ON DELETE CASCADE
    );

CREATE TRIGGER update_tags_updated_at
    BEFORE UPDATE ON aion_api.tags
    FOR EACH ROW
    EXECUTE FUNCTION aion_api.update_timestamp();