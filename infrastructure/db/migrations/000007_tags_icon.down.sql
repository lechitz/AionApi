-- Migration: 000007_tags_icon
-- Description: Remove emoji icon from tags

ALTER TABLE aion_api.tags
DROP COLUMN IF EXISTS icon;
