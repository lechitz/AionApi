CREATE TABLE IF NOT EXISTS aion_api.day_energy
(
    id         SERIAL PRIMARY KEY,                                       -- Unique Energy ID
    day_id     INT                                NOT NULL,              -- Foreign Key to Day
    energy     INT CHECK (energy BETWEEN 1 AND 5) NOT NULL,              -- Energy level (1-5)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,                      -- Timestamp of creation
    FOREIGN KEY (day_id) REFERENCES aion_api.days (id) ON DELETE CASCADE -- Relationship to Day
);
