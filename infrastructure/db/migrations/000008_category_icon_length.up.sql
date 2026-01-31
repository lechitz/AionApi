-- Migration: 000008_category_icon_length
-- Description: Increase categories icon key length for SVG paths

ALTER TABLE aion_api.categories
ALTER COLUMN icon TYPE VARCHAR(120);
