-- Migration: 000008_category_icon_length
-- Description: Revert categories icon key length

ALTER TABLE aion_api.categories
ALTER COLUMN icon TYPE VARCHAR(50);
