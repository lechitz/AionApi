CREATE TABLE IF NOT EXISTS aion_api.days
(
    id            SERIAL PRIMARY KEY,
    user_id       INT NOT NULL,
    date          DATE NOT NULL UNIQUE,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES aion_api.users (user_id) ON DELETE CASCADE
    );