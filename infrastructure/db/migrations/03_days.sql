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
