CREATE TABLE IF NOT EXISTS aion_api.day_energy
(
    id         SERIAL PRIMARY KEY,
    day_id     INT                                NOT NULL,
    energy     INT CHECK (energy BETWEEN 1 AND 5) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,
    CONSTRAINT fk_day_energy_day FOREIGN KEY (day_id) REFERENCES aion_api.days (id) ON DELETE CASCADE
);
