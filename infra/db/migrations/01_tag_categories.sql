CREATE TABLE IF NOT EXISTS aion_api.tag_categories (
    category_id   SERIAL PRIMARY KEY,
    user_id       INT NOT NULL,
    name          VARCHAR(40) NOT NULL UNIQUE,
    description   VARCHAR(200),
    color_hex     VARCHAR(7),
    icon          VARCHAR(50),
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMP DEFAULT NULL
);