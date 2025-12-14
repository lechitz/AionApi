-- ============================================================================
-- Records Seed Data - User ID 1 (Test User, generic names)
-- ============================================================================
-- Generates up to 128 records for user_id = 1 when the required tags exist.
-- Rows with missing tags are skipped via WHERE EXISTS guards.
-- ============================================================================

WITH tag_list(name) AS (
  VALUES
    ('user1_tag_01'), ('user1_tag_02'),
    ('user1_tag_03'), ('user1_tag_04'),
    ('user1_tag_05'), ('user1_tag_06'),
    ('user1_tag_07'), ('user1_tag_08'),
    ('user1_tag_09'), ('user1_tag_10'),
    ('user1_tag_11'), ('user1_tag_12'),
    ('user1_tag_13'), ('user1_tag_14'),
    ('user1_tag_15'), ('user1_tag_16')
),
expanded AS (
  SELECT name
  FROM tag_list
  CROSS JOIN generate_series(1, 3) AS rep
),
numbered AS (
  SELECT name, ROW_NUMBER() OVER (ORDER BY name, (SELECT 1)) AS rn
  FROM expanded
),
limited AS (
  SELECT name, rn
  FROM numbered
  WHERE rn <= 128
)
INSERT INTO aion_api.records (
  user_id, title, description, category_id, tag_id,
  event_time, recorded_at, duration_seconds, value,
  source, timezone, status
)
SELECT
  1 AS user_id,
  limited.name || ' activity ' || limited.rn AS title,
  'auto_seed_' || limited.rn AS description,
  (SELECT t.category_id FROM aion_api.tags t WHERE t.name = limited.name AND t.user_id = 1 AND t.deleted_at IS NULL LIMIT 1) AS category_id,
  (SELECT t.tag_id FROM aion_api.tags t WHERE t.name = limited.name AND t.user_id = 1 AND t.deleted_at IS NULL LIMIT 1) AS tag_id,
  ('2025-01-01 06:00:00'::timestamp + (limited.rn || ' hours')::interval) AS event_time,
  ('2025-01-01 06:00:00'::timestamp + (limited.rn || ' hours')::interval) AS recorded_at,
  (1800 + ((limited.rn % 4) * 900)) AS duration_seconds,
  limited.rn::decimal AS value,
  'seed' AS source,
  'America/Sao_Paulo' AS timezone,
  'completed' AS status
FROM limited
WHERE EXISTS (
  SELECT 1 FROM aion_api.tags t WHERE t.name = limited.name AND t.user_id = 1 AND t.deleted_at IS NULL
);
