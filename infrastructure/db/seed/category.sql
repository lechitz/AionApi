-- ============================================================================
-- Category Seed Data - User ID 1 (Test User, generic names)
-- Idempotent via WHERE NOT EXISTS to match unique index on (user_id, lower(name))
-- ============================================================================

-- category 1
INSERT INTO aion_api.categories (user_id, name, description, color_hex, icon, deleted_at)
SELECT 1, 'user1_category_1', 'user1_category_1_description', '#E94F37', 'dumbbell', NULL
WHERE NOT EXISTS (
  SELECT 1 FROM aion_api.categories WHERE user_id = 1 AND lower(name) = lower('user1_category_1') AND deleted_at IS NULL
);

-- category 2
INSERT INTO aion_api.categories (user_id, name, description, color_hex, icon, deleted_at)
SELECT 1, 'user1_category_2', 'user1_category_2_description', '#9C27B0', 'spa', NULL
WHERE NOT EXISTS (
  SELECT 1 FROM aion_api.categories WHERE user_id = 1 AND lower(name) = lower('user1_category_2') AND deleted_at IS NULL
);

-- category 3
INSERT INTO aion_api.categories (user_id, name, description, color_hex, icon, deleted_at)
SELECT 1, 'user1_category_3', 'user1_category_3_description', '#F8B400', 'brain', NULL
WHERE NOT EXISTS (
  SELECT 1 FROM aion_api.categories WHERE user_id = 1 AND lower(name) = lower('user1_category_3') AND deleted_at IS NULL
);

-- category 4
INSERT INTO aion_api.categories (user_id, name, description, color_hex, icon, deleted_at)
SELECT 1, 'user1_category_4', 'user1_category_4_description', '#1976D2', 'briefcase', NULL
WHERE NOT EXISTS (
  SELECT 1 FROM aion_api.categories WHERE user_id = 1 AND lower(name) = lower('user1_category_4') AND deleted_at IS NULL
);

-- category 5
INSERT INTO aion_api.categories (user_id, name, description, color_hex, icon, deleted_at)
SELECT 1, 'user1_category_5', 'user1_category_5_description', '#00ACC1', 'globe', NULL
WHERE NOT EXISTS (
  SELECT 1 FROM aion_api.categories WHERE user_id = 1 AND lower(name) = lower('user1_category_5') AND deleted_at IS NULL
);

-- category 6
INSERT INTO aion_api.categories (user_id, name, description, color_hex, icon, deleted_at)
SELECT 1, 'user1_category_6', 'user1_category_6_description', '#FF6F00', 'user', NULL
WHERE NOT EXISTS (
  SELECT 1 FROM aion_api.categories WHERE user_id = 1 AND lower(name) = lower('user1_category_6') AND deleted_at IS NULL
);

-- category 7
INSERT INTO aion_api.categories (user_id, name, description, color_hex, icon, deleted_at)
SELECT 1, 'user1_category_7', 'user1_category_7_description', '#388E3C', 'home', NULL
WHERE NOT EXISTS (
  SELECT 1 FROM aion_api.categories WHERE user_id = 1 AND lower(name) = lower('user1_category_7') AND deleted_at IS NULL
);

-- category 8
INSERT INTO aion_api.categories (user_id, name, description, color_hex, icon, deleted_at)
SELECT 1, 'user1_category_8', 'user1_category_8_description', '#616161', 'ellipsis-h', NULL
WHERE NOT EXISTS (
  SELECT 1 FROM aion_api.categories WHERE user_id = 1 AND lower(name) = lower('user1_category_8') AND deleted_at IS NULL
);
