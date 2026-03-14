-- category_generate.sql
-- Generates categories for N users using PostgreSQL generate_series
-- Usage: psql -v seed_count=5 -f category_generate.sql
-- This creates 8 categories per user with naming like: saude_fisica1, saude_fisica2, etc.

\if :{?seed_count}
\else
\set seed_count 5
\endif

-- Generate categories for each user (8 categories per user)
INSERT INTO aion_api.categories (user_id, name, description, color_hex, icon, deleted_at)
SELECT
  user_id,
  category_name || '_' || user_id AS name,
  description || '_' || user_id AS description,
  color_hex,
  icon,
  NULL AS deleted_at
FROM
  generate_series(1, :seed_count) AS user_id,
  (VALUES
    ('saude_fisica', 'Atividades físicas e condicionamento', '#E94F37', '🏋️'),
    ('meditacao', 'Práticas de meditação e atenção plena', '#9C27B0', '🧘'),
    ('saude_mental', 'Saúde mental, planejamento e bem-estar', '#F8B400', '🧠'),
    ('estudo_trabalho', 'Estudo, trabalho e desenvolvimento', '#1976D2', '💼'),
    ('idiomas', 'Atividades de aprendizado de idiomas', '#00ACC1', '🌍'),
    ('pessoal', 'Assuntos pessoais, lazer e tempo livre', '#FF6F00', '🎉'),
    ('trabalho_de_casa', 'Tarefas domésticas e manutenção da casa', '#388E3C', '🏠'),
    ('outros', 'Atividades diversas / off / viagens', '#616161', '✨')
  ) AS categories(category_name, description, color_hex, icon)
ON CONFLICT (user_id, lower(name)) WHERE deleted_at IS NULL DO NOTHING;
