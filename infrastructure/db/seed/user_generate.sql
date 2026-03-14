-- user_generate.sql
-- Generates N users using Postgres' pgcrypto to hash passwords server-side.
-- Usage with psql:
--  psql -v seed_count=60 -v user_seed_password_plain='password123' -f user_generate.sql
-- When using docker exec via Makefile, pass -v seed_count and -v user_seed_password_plain.

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Provide safe defaults when vars are missing (older psql versions may not support \error)
\if :{?user_seed_password_plain}
\else
\set user_seed_password_plain 'testpassword123'
\endif

\if :{?seed_count}
\else
\set seed_count 10
\endif

INSERT INTO aion_api.users (name, username, password, email)
SELECT
  ('user_name_' || i) AS name,
  ('user' || i) AS username,
  crypt(:'user_seed_password_plain', gen_salt('bf', 10)) AS password,
  ('user' || i || '@aion.com') AS email
FROM generate_series(1, :seed_count) AS s(i)
ON CONFLICT (username) DO NOTHING;

