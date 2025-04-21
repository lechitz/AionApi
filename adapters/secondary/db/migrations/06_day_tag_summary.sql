CREATE TABLE IF NOT EXISTS aion_api.day_tag_summary
(
    id           SERIAL PRIMARY KEY,
    day_id       INT NOT NULL,
    tag_id       INT NOT NULL,
    summary      TEXT,
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (day_id) REFERENCES aion_api.days (id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES aion_api.tags (id) ON DELETE CASCADE
);