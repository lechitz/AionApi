CREATE TABLE IF NOT EXISTS aion_api.professional_diaries
(
    id                SERIAL PRIMARY KEY,
    day_id            INT NOT NULL,
    work_start        TIME NOT NULL,
    lunch_start       TIME,
    lunch_end         TIME,
    work_end          TIME NOT NULL,
    content           TEXT,
    created_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at        TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (day_id) REFERENCES aion_api.days (id) ON DELETE CASCADE
);