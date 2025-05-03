CREATE TABLE IF NOT EXISTS aion_api.tags (
    tag_id             SERIAL PRIMARY KEY,
    user_id        INT NOT NULL,
    name           VARCHAR(255) NOT NULL,
    category_id    INT NOT NULL,
    description    TEXT,
    creation_date  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at     TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES aion_api.users (id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES aion_api.tag_categories (id) ON DELETE RESTRICT
    );

CREATE TRIGGER update_tags_updated_at
    BEFORE UPDATE ON aion_api.tags
    FOR EACH ROW
    EXECUTE FUNCTION aion_api.update_timestamp();
