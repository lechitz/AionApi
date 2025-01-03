CREATE TABLE aion_api.diaries
(
    id          SERIAL PRIMARY KEY,                                   -- Unique Diary ID
    user_id     INT NOT NULL,                                          -- Foreign Key to User
    date        DATE NOT NULL,                                         -- Date of the diary
    type        VARCHAR(255) NOT NULL,                                 -- Diary type (personal or professional)
    title       VARCHAR(255),                                          -- Title of the diary
    content     TEXT,                                                  -- Content of the diary
    FOREIGN KEY (user_id) REFERENCES aion_api.users (id) ON DELETE CASCADE -- Relationship to User
);