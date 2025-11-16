-- ============================================================================
-- Records Seed Data - User ID 1 (Test User)
-- ============================================================================
-- Generates exactly 128 records for user_id = 1 when the required tags exist.
-- Rows with missing tags are skipped via WHERE EXISTS guards.
-- ============================================================================

WITH tag_list(name) AS (
  VALUES
    ('Stretching'), ('Push-up'), ('ABS'), ('Run'), ('Pull-up'), ('Sunbathe'), ('Gym'), ('Walking'),
    ('Meditation'),
    ('Planning'), ('Diary'), ('Reading'), ('TheNews'), ('Emails'),
    ('Dev'), ('College'), ('Golang'), ('Notion'), ('GPT'), ('Full Cycle'), ('FreeCodeCamp'), ('Coursera'), ('Course'), ('Aion'), ('Work'), ('RD'), ('Podcast'), ('AudioBook'), ('Finance'), ('Interview'),
    ('English'), ('Spanish'), ('German'), ('French'), ('Chinese'),
    ('My'), ('Travel'), ('Movie'), ('Series'), ('Game'), ('Beach'), ('Hanging Out'),
    ('Housework'), ('Supermarket'), ('Doctor'),
    ('OFF'), ('Travelling')
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
  source, timezone, status, created_at, updated_at
)
SELECT
  1 AS user_id,
  limited.name || ' activity ' || limited.rn AS title,
  'Auto seed ' || limited.rn AS description,
  (SELECT t.category_id FROM aion_api.tags t WHERE t.name = limited.name AND t.user_id = 1 AND t.deleted_at IS NULL LIMIT 1) AS category_id,
  (SELECT t.tag_id FROM aion_api.tags t WHERE t.name = limited.name AND t.user_id = 1 AND t.deleted_at IS NULL LIMIT 1) AS tag_id,
  ('2025-01-01 06:00:00'::timestamp + (limited.rn || ' hours')::interval) AS event_time,
  ('2025-01-01 06:00:00'::timestamp + (limited.rn || ' hours')::interval) AS recorded_at,
  (1800 + ((limited.rn % 4) * 900)) AS duration_seconds,
  limited.rn::decimal AS value,
  'seed' AS source,
  'America/Sao_Paulo' AS timezone,
  'completed' AS status,
  ('2025-01-01 06:00:00'::timestamp + (limited.rn || ' hours')::interval) AS created_at,
  ('2025-01-01 06:00:00'::timestamp + (limited.rn || ' hours')::interval) AS updated_at
FROM limited
WHERE EXISTS (
  SELECT 1 FROM aion_api.tags t WHERE t.name = limited.name AND t.user_id = 1 AND t.deleted_at IS NULL
);
