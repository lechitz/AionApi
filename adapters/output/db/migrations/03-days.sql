CREATE TABLE aion_api.days
(
    id            SERIAL PRIMARY KEY,                               -- Unique Day ID
    user_id       INT NOT NULL,                                     -- Foreign Key to User
    date          DATE NOT NULL,                                    -- Day's Date
    mood          TEXT,                                             -- Mood of the day
    energy        INT CHECK (energy BETWEEN 1 AND 5),               -- Energy (1-5)
    water_intake  INT CHECK (water_intake >= 0),                     -- Water intake
    intention     TEXT,                                             -- Daily intentions
    FOREIGN KEY (user_id) REFERENCES aion_api.users (id) ON DELETE CASCADE -- Relationship to User
);