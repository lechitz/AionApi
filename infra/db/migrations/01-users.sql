--- New Schema
CREATE SCHEMA IF NOT EXISTS aion_api;

CREATE TABLE IF NOT EXISTS aion_api.users
(
    id         SERIAL NOT NULL UNIQUE,
    name       VARCHAR(255) NOT NULL,
    username   VARCHAR(255) NOT NULL UNIQUE,
    password   VARCHAR(255) NOT NULL,
    email      VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);

-- Function to update updated_at column
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to update updated_at column
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON aion_api.users
    FOR EACH ROW
    EXECUTE FUNCTION update_timestamp();
