-- records_generate.sql
-- Generates records for N users based on their tags
-- Usage: psql -v seed_count=5 -v days=7 -f records_generate.sql
-- This creates daily records for each user's tags

\if :{?seed_count}
\else
\set seed_count 5
\endif

\if :{?days}
\else
\set days 7
\endif

-- Generate 0-5 records per tag for each user (spread over :days days)
INSERT INTO aion_api.records (
  user_id,
  category_id,
  tag_id,
  title,
  description,
  event_time,
  recorded_at,
  duration_seconds,
  value,
  source,
  timezone,
  status
)
SELECT
  u.user_id,
  t.category_id,
  t.tag_id,
  t.name || ' record ' || rec_idx AS title,
  'auto_seed_' || rec_idx AS description,
  ts.event_ts AS event_time,
  ts.event_ts + interval '30 minutes' AS recorded_at,
  1800 AS duration_seconds,
  (10 + rec_idx) AS value,
  'seed' AS source,
  'UTC' AS timezone,
  'completed' AS status
FROM
  generate_series(1, :seed_count) AS u(user_id)
INNER JOIN aion_api.tags t
  ON t.user_id = u.user_id
  AND t.deleted_at IS NULL
INNER JOIN aion_api.categories c
  ON c.category_id = t.category_id
  AND c.deleted_at IS NULL
CROSS JOIN LATERAL generate_series(0, 5) AS rec_idx
CROSS JOIN LATERAL (
  SELECT ((CURRENT_DATE - (rec_idx % GREATEST(:days,1)))::timestamp + (interval '1 hour' * (8 + (rec_idx % 12)))) AT TIME ZONE 'UTC' AS event_ts
) AS ts
ORDER BY u.user_id, t.tag_id, rec_idx;
