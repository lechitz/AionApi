CREATE TABLE IF NOT EXISTS aion_api.day_tag_summary
(
    id           SERIAL PRIMARY KEY,                               -- Unique Summary ID
    day_id       INT NOT NULL,                                     -- Foreign Key to Day
    tag_id       INT NOT NULL,                                     -- Foreign Key to Tag
    summary      TEXT,                                             -- Summary of the tag on the day
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,             -- Timestamp of creation
    FOREIGN KEY (day_id) REFERENCES aion_api.days (id) ON DELETE CASCADE, -- Relationship to Day
    FOREIGN KEY (tag_id) REFERENCES aion_api.tags (id) ON DELETE CASCADE  -- Relationship to Tag
);