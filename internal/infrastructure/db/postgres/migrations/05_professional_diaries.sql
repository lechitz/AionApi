CREATE TABLE IF NOT EXISTS aion_api.professional_diaries
(
    id                SERIAL PRIMARY KEY,
    day_id INT NOT NULL,
    start_time        TIME NOT NULL,
    lunch_start       TIME,
    lunch_end         TIME,
    end_time          TIME NOT NULL,
    content             TEXT, --TODO: PODE SER UM VALOR NULO:
    FOREIGN KEY (day_id) REFERENCES aion_api.days (id) ON DELETE CASCADE
);