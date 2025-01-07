CREATE TABLE IF NOT EXISTS aion_api.personal_diaries
(
    id         SERIAL PRIMARY KEY,
    day_id     INT       NOT NULL,
    content      TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (day_id) REFERENCES aion_api.days (id) ON DELETE CASCADE
);