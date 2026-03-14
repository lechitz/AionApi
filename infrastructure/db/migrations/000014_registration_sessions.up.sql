-- Migration: 000014_registration_sessions
-- Description: Add staged public registration sessions for multi-step onboarding

CREATE TABLE IF NOT EXISTS aion_api.registration_sessions
(
    registration_id UUID PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    username        VARCHAR(255) NOT NULL,
    email           VARCHAR(255) NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    locale          VARCHAR(16) DEFAULT NULL,
    timezone        VARCHAR(64) DEFAULT NULL,
    location        VARCHAR(255) DEFAULT NULL,
    bio             TEXT DEFAULT NULL,
    avatar_url      TEXT DEFAULT NULL,
    current_step    SMALLINT NOT NULL DEFAULT 1 CHECK (current_step BETWEEN 1 AND 3),
    status          VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'expired', 'canceled')),
    expires_at      TIMESTAMP NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_registration_sessions_username
    ON aion_api.registration_sessions (username);

CREATE INDEX IF NOT EXISTS idx_registration_sessions_email
    ON aion_api.registration_sessions (email);

CREATE INDEX IF NOT EXISTS idx_registration_sessions_status_expires_at
    ON aion_api.registration_sessions (status, expires_at);

