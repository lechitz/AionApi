-- Migration: 000002_users_and_roles
-- Description: Create users, roles, and user_roles tables
-- This consolidates: 01a_roles.sql, 01b_users.sql, 01c_user_roles.sql

-- Roles table: stores all valid system roles
CREATE TABLE IF NOT EXISTS aion_api.roles
(
    role_id     SERIAL PRIMARY KEY,
    name        VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Insert default roles with hierarchy: owner > admin > user > blocked
INSERT INTO aion_api.roles (name, description) VALUES
    ('owner', 'System owner with highest privileges'),
    ('admin', 'Administrator with full system access'),
    ('user', 'Default user role with basic access'),
    ('blocked', 'Blocked user with no access')
ON CONFLICT (name) DO NOTHING;

-- Users table
CREATE TABLE IF NOT EXISTS aion_api.users
(
    user_id    SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    username   VARCHAR(255) NOT NULL UNIQUE,
    password   VARCHAR(255) NOT NULL,
    email      VARCHAR(255) NOT NULL UNIQUE,
    locale     VARCHAR(16) DEFAULT NULL,
    timezone   VARCHAR(64) DEFAULT NULL,
    location   VARCHAR(255) DEFAULT NULL,
    bio        TEXT DEFAULT NULL,
    avatar_url TEXT DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);

-- User Roles junction table: many-to-many relationship between users and roles
CREATE TABLE IF NOT EXISTS aion_api.user_roles
(
    user_role_id SERIAL PRIMARY KEY,
    user_id      INTEGER NOT NULL REFERENCES aion_api.users(user_id) ON DELETE CASCADE,
    role_id      INTEGER NOT NULL REFERENCES aion_api.roles(role_id) ON DELETE CASCADE,
    assigned_by  INTEGER REFERENCES aion_api.users(user_id) ON DELETE SET NULL,
    assigned_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, role_id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON aion_api.user_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON aion_api.user_roles(role_id);
