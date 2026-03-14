-- Migration: 000015_user_onboarding
-- Description: Add onboarding_completed flag to users table

ALTER TABLE aion_api.users
    ADD COLUMN IF NOT EXISTS onboarding_completed BOOLEAN NOT NULL DEFAULT FALSE;

COMMENT ON COLUMN aion_api.users.onboarding_completed IS 'Indicates if the user has completed the onboarding flow';
