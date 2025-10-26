CREATE TABLE IF NOT EXISTS aion_api.tags (
    tag_id        SERIAL PRIMARY KEY,
    user_id       INT NOT NULL,
    name          VARCHAR(255) NOT NULL,
    category_id   INT NULL,
    description   TEXT,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES aion_api.users (user_id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES aion_api.tag_categories (category_id) ON DELETE SET NULL
    );

CREATE UNIQUE INDEX IF NOT EXISTS ux_tags_user_name_ci
ON aion_api.tags (user_id, lower(name))
WHERE deleted_at IS NULL;

CREATE TRIGGER update_tags_updated_at
    BEFORE UPDATE ON aion_api.tags
    FOR EACH ROW
    EXECUTE FUNCTION aion_api.update_timestamp();
