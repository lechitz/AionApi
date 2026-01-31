-- Migration: 000003_categories_tags_days
-- Description: Create categories, tags, and days tables
-- This consolidates: 02_categories.sql, 03_tags.sql, 03_days.sql

-- Categories table
CREATE TABLE IF NOT EXISTS aion_api.categories (
    category_id   SERIAL PRIMARY KEY,
    user_id       INT NOT NULL,
    name          VARCHAR(40) NOT NULL,
    description   VARCHAR(200),
    color_hex     VARCHAR(7),
    icon          VARCHAR(50),
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMP DEFAULT NULL,
    CONSTRAINT fk_categories_user FOREIGN KEY (user_id) REFERENCES aion_api.users (user_id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_categories_user_name_ci
ON aion_api.categories (user_id, lower(name))
WHERE deleted_at IS NULL;

DROP TRIGGER IF EXISTS update_categories_updated_at ON aion_api.categories;
CREATE TRIGGER update_categories_updated_at
    BEFORE UPDATE ON aion_api.categories
    FOR EACH ROW
    EXECUTE FUNCTION aion_api.update_timestamp();

-- Days table
CREATE TABLE IF NOT EXISTS aion_api.days
(
    id            SERIAL PRIMARY KEY,
    user_id       INT NOT NULL,
    date          DATE NOT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES aion_api.users (user_id) ON DELETE CASCADE
);

-- Tags table
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
    FOREIGN KEY (category_id) REFERENCES aion_api.categories (category_id) ON DELETE SET NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_tags_user_name_ci
ON aion_api.tags (user_id, lower(name))
WHERE deleted_at IS NULL;

DROP TRIGGER IF EXISTS update_tags_updated_at ON aion_api.tags;
CREATE TRIGGER update_tags_updated_at
    BEFORE UPDATE ON aion_api.tags
    FOR EACH ROW
    EXECUTE FUNCTION aion_api.update_timestamp();
