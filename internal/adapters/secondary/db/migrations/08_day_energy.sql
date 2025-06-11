CREATE TABLE IF NOT EXISTS aion_api.day_energy
(
    id         SERIAL PRIMARY KEY,
    day_id     INT                                NOT NULL,
    energy     INT CHECK (energy BETWEEN 1 AND 5) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (day_id) REFERENCES aion_api.days (id) ON DELETE CASCADE
);
