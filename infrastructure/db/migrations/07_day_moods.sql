CREATE TABLE IF NOT EXISTS aion_api.day_moods
(
    id         SERIAL PRIMARY KEY,
    day_id     INT  NOT NULL,
    mood       TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (day_id) REFERENCES aion_api.days (id) ON DELETE CASCADE
);
