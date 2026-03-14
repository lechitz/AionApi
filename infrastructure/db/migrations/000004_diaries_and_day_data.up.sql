-- Migration: 000004_diaries_and_day_data
-- Description: Create diary and day-related tables
-- This consolidates: 04_personal_diaries.sql, 05_professional_diaries.sql,
--                    06_day_tag_summary.sql, 07_day_moods.sql, 08_day_energy.sql,
--                    09_day_water_intake.sql, 10_day_intentions.sql

-- Personal diaries
CREATE TABLE IF NOT EXISTS aion_api.personal_diaries
(
    id         SERIAL PRIMARY KEY,
    day_id     INT       NOT NULL,
    content      TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (day_id) REFERENCES aion_api.days (id) ON DELETE CASCADE
);

-- Professional diaries
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

-- Day tag summary
CREATE TABLE IF NOT EXISTS aion_api.day_tag_summary
(
    id           SERIAL PRIMARY KEY,
    day_id       INT NOT NULL,
    tag_id       INT NOT NULL,
    summary      TEXT,
    created_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (day_id) REFERENCES aion_api.days (id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES aion_api.tags (tag_id) ON DELETE CASCADE
);

-- Day moods
CREATE TABLE IF NOT EXISTS aion_api.day_moods
(
    id         SERIAL PRIMARY KEY,
    day_id     INT  NOT NULL,
    mood       TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (day_id) REFERENCES aion_api.days (id) ON DELETE CASCADE
);

-- Day energy
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

-- Day water intake
CREATE TABLE IF NOT EXISTS aion_api.day_water_intake
(
    id            SERIAL PRIMARY KEY,
    day_id        INT NOT NULL,
    amount_ml     INT CHECK (amount_ml > 0) NOT NULL,
    consumed_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (day_id) REFERENCES aion_api.days (id) ON DELETE CASCADE
);

-- Day intentions
CREATE TABLE IF NOT EXISTS aion_api.day_intentions
(
    id            SERIAL PRIMARY KEY,
    day_id        INT NOT NULL,
    intention     TEXT NOT NULL,
    is_completed  BOOLEAN DEFAULT FALSE,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (day_id) REFERENCES aion_api.days (id) ON DELETE CASCADE
);

DROP TRIGGER IF EXISTS update_day_intentions_updated_at ON aion_api.day_intentions;
CREATE TRIGGER update_day_intentions_updated_at
    BEFORE UPDATE ON aion_api.day_intentions
    FOR EACH ROW
    EXECUTE FUNCTION aion_api.update_timestamp();
