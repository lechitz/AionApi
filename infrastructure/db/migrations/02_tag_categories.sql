CREATE TABLE IF NOT EXISTS aion_api.tag_categories (
    category_id   SERIAL PRIMARY KEY,
    user_id       INT NOT NULL,
    name          VARCHAR(40) NOT NULL,
    description   VARCHAR(200),
    color_hex     VARCHAR(7),
    icon          VARCHAR(50),
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMP DEFAULT NULL,
    CONSTRAINT fk_tag_categories_user FOREIGN KEY (user_id) REFERENCES aion_api.users (user_id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_tag_categories_user_name_ci
ON aion_api.tag_categories (user_id, lower(name))
WHERE deleted_at IS NULL;

CREATE TRIGGER update_tag_categories_updated_at
    BEFORE UPDATE ON aion_api.tag_categories
    FOR EACH ROW
    EXECUTE FUNCTION aion_api.update_timestamp();
