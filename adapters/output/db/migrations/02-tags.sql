CREATE TABLE aion_api.tags
(
    id               SERIAL PRIMARY KEY,                              -- Unique Tag ID
    user_id          INT NOT NULL,                                    -- Foreign Key to User
    name             VARCHAR(255) NOT NULL,                            -- Tag name
    category         VARCHAR(255) NOT NULL,                            -- Tag category (e.g., Exercise, Work, etc.)
    description      TEXT,                                             -- Tag description
    creation_date    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,              -- Tag creation date
    update_date      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,              -- Tag last update date
    FOREIGN KEY (user_id) REFERENCES aion_api.users (id) ON DELETE CASCADE -- Relationship to User
);