-- Migration: 000007_tags_icon
-- Description: Add emoji icon to tags

ALTER TABLE aion_api.tags
ADD COLUMN IF NOT EXISTS icon VARCHAR(50) NOT NULL DEFAULT '⬜';
