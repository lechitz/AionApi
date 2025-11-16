CREATE TABLE IF NOT EXISTS aion_api.day_tag_summary
(
    id           SERIAL PRIMARY KEY,
    day_id       INT NOT NULL,
    tag_id       INT NOT NULL,
    summary      TEXT,
    created_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (day_id) REFERENCES aion_api.days (id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES aion_api.tags (tag_id) ON DELETE CASCADE
    );