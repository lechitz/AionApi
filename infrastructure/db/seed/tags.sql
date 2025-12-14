-- ============================================================================
-- Tags Seed Data - User ID 1 (Test User, generic names)
-- Idempotent via WHERE NOT EXISTS to match unique index on (user_id, lower(name))
-- ============================================================================

-- category 1 tags
INSERT INTO aion_api.tags (user_id, name, category_id, description, deleted_at)
SELECT 1, 'user1_tag_01', (SELECT category_id FROM aion_api.tag_categories WHERE name='user1_category_1' AND user_id=1 LIMIT 1), 'user1_tag_01_description', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('user1_tag_01') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, deleted_at)
SELECT 1, 'user1_tag_02', (SELECT category_id FROM aion_api.tag_categories WHERE name='user1_category_1' AND user_id=1 LIMIT 1), 'user1_tag_02_description', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('user1_tag_02') AND deleted_at IS NULL);

-- category 2 tags
INSERT INTO aion_api.tags (user_id, name, category_id, description, deleted_at)
SELECT 1, 'user1_tag_03', (SELECT category_id FROM aion_api.tag_categories WHERE name='user1_category_2' AND user_id=1 LIMIT 1), 'user1_tag_03_description', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('user1_tag_03') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, deleted_at)
SELECT 1, 'user1_tag_04', (SELECT category_id FROM aion_api.tag_categories WHERE name='user1_category_2' AND user_id=1 LIMIT 1), 'user1_tag_04_description', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('user1_tag_04') AND deleted_at IS NULL);

-- category 3 tags
INSERT INTO aion_api.tags (user_id, name, category_id, description, deleted_at)
SELECT 1, 'user1_tag_05', (SELECT category_id FROM aion_api.tag_categories WHERE name='user1_category_3' AND user_id=1 LIMIT 1), 'user1_tag_05_description', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('user1_tag_05') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, deleted_at)
SELECT 1, 'user1_tag_06', (SELECT category_id FROM aion_api.tag_categories WHERE name='user1_category_3' AND user_id=1 LIMIT 1), 'user1_tag_06_description', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('user1_tag_06') AND deleted_at IS NULL);

-- category 4 tags
INSERT INTO aion_api.tags (user_id, name, category_id, description, deleted_at)
SELECT 1, 'user1_tag_07', (SELECT category_id FROM aion_api.tag_categories WHERE name='user1_category_4' AND user_id=1 LIMIT 1), 'user1_tag_07_description', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('user1_tag_07') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, deleted_at)
SELECT 1, 'user1_tag_08', (SELECT category_id FROM aion_api.tag_categories WHERE name='user1_category_4' AND user_id=1 LIMIT 1), 'user1_tag_08_description', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('user1_tag_08') AND deleted_at IS NULL);

-- category 5 tags
INSERT INTO aion_api.tags (user_id, name, category_id, description, deleted_at)
SELECT 1, 'user1_tag_09', (SELECT category_id FROM aion_api.tag_categories WHERE name='user1_category_5' AND user_id=1 LIMIT 1), 'user1_tag_09_description', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('user1_tag_09') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, deleted_at)
SELECT 1, 'user1_tag_10', (SELECT category_id FROM aion_api.tag_categories WHERE name='user1_category_5' AND user_id=1 LIMIT 1), 'user1_tag_10_description', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('user1_tag_10') AND deleted_at IS NULL);

-- category 6 tags
INSERT INTO aion_api.tags (user_id, name, category_id, description, deleted_at)
SELECT 1, 'user1_tag_11', (SELECT category_id FROM aion_api.tag_categories WHERE name='user1_category_6' AND user_id=1 LIMIT 1), 'user1_tag_11_description', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('user1_tag_11') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, deleted_at)
SELECT 1, 'user1_tag_12', (SELECT category_id FROM aion_api.tag_categories WHERE name='user1_category_6' AND user_id=1 LIMIT 1), 'user1_tag_12_description', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('user1_tag_12') AND deleted_at IS NULL);

-- category 7 tags
INSERT INTO aion_api.tags (user_id, name, category_id, description, deleted_at)
SELECT 1, 'user1_tag_13', (SELECT category_id FROM aion_api.tag_categories WHERE name='user1_category_7' AND user_id=1 LIMIT 1), 'user1_tag_13_description', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('user1_tag_13') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, deleted_at)
SELECT 1, 'user1_tag_14', (SELECT category_id FROM aion_api.tag_categories WHERE name='user1_category_7' AND user_id=1 LIMIT 1), 'user1_tag_14_description', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('user1_tag_14') AND deleted_at IS NULL);

-- category 8 tags
INSERT INTO aion_api.tags (user_id, name, category_id, description, deleted_at)
SELECT 1, 'user1_tag_15', (SELECT category_id FROM aion_api.tag_categories WHERE name='user1_category_8' AND user_id=1 LIMIT 1), 'user1_tag_15_description', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('user1_tag_15') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, deleted_at)
SELECT 1, 'user1_tag_16', (SELECT category_id FROM aion_api.tag_categories WHERE name='user1_category_8' AND user_id=1 LIMIT 1), 'user1_tag_16_description', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('user1_tag_16') AND deleted_at IS NULL);
