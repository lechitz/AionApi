-- Migration: 000015_user_onboarding (down)
-- Description: Remove onboarding_completed flag from users table

ALTER TABLE aion_api.users
    DROP COLUMN IF EXISTS onboarding_completed;
