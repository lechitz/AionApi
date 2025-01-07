CREATE TABLE IF NOT EXISTS aion_api.day_water_intake
(
    id            SERIAL PRIMARY KEY,
    day_id        INT NOT NULL,
    amount_ml     INT CHECK (amount_ml > 0) NOT NULL,
    consumed_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (day_id) REFERENCES aion_api.days (id) ON DELETE CASCADE
);