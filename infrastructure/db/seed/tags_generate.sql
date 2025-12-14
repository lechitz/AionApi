-- tags_generate.sql
-- Generates tags for N users based on their categories
-- Usage: psql -v seed_count=5 -f tags_generate.sql
-- This creates tags with naming like: Stretching1, Stretching2, etc.

\if :{?seed_count}
\else
\set seed_count 5
\endif

-- Generate tags for each user (7 tags per category, 8 categories = 56 tags per user)
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT
  u.user_id,
  t.tag_name || '_' || u.user_id AS name,
  c.category_id,
  t.description || '_' || u.user_id AS description,
  NOW() - (interval '1 day' * (10 - u.user_id)) AS created_at,
  NOW() - (interval '1 day' * (10 - u.user_id)) AS updated_at,
  NULL AS deleted_at
FROM
  generate_series(1, :seed_count) AS u(user_id)
CROSS JOIN (
  VALUES
    -- saude_fisica tags
    ('Stretching', 'Stretching and flexibility', 'saude_fisica'),
    ('Push-up', 'Push-up exercise', 'saude_fisica'),
    ('ABS', 'Abdominal exercises', 'saude_fisica'),
    ('Run', 'Running sessions', 'saude_fisica'),
    ('Pull-up', 'Pull-up / upper body', 'saude_fisica'),
    ('Sunbathe', 'Sun exposure / relaxation', 'saude_fisica'),
    ('Gym', 'Gym workouts', 'saude_fisica'),
    ('Walking', 'Walking sessions', 'saude_fisica'),

    -- meditacao tags
    ('Meditation', 'Meditation practice', 'meditacao'),

    -- saude_mental tags
    ('Planning', 'Planning and organizing', 'saude_mental'),
    ('Diary', 'Personal diary entries', 'saude_mental'),
    ('Reading', 'Reading and books', 'saude_mental'),
    ('TheNews', 'Reading the news', 'saude_mental'),
    ('Emails', 'Email processing', 'saude_mental'),
    ('Thinking', 'Reflection and thinking', 'saude_mental'),

    -- estudo_trabalho tags
    ('DevStudy', 'Development and coding study', 'estudo_trabalho'),
    ('Coding', 'Programming and coding', 'estudo_trabalho'),
    ('SoftSkills', 'Soft skills development', 'estudo_trabalho'),
    ('Search', 'Research and investigation', 'estudo_trabalho'),
    ('WriteDoc', 'Writing documentation', 'estudo_trabalho'),

    -- idiomas tags
    ('English', 'English language practice', 'idiomas'),
    ('German', 'German language practice', 'idiomas'),
    ('Japanese', 'Japanese language practice', 'idiomas'),
    ('Chinese', 'Chinese language practice', 'idiomas'),
    ('Libras', 'Brazilian sign language', 'idiomas'),

    -- pessoal tags
    ('YouTube', 'Watching YouTube videos', 'pessoal'),
    ('Movies', 'Watching movies', 'pessoal'),
    ('Series', 'Watching series', 'pessoal'),
    ('Netflix', 'Watching Netflix', 'pessoal'),
    ('Gaming', 'Playing games', 'pessoal'),
    ('Music', 'Listening to music', 'pessoal'),
    ('Podcasts', 'Listening to podcasts', 'pessoal'),
    ('Social', 'Social media', 'pessoal'),
    ('Shopping', 'Shopping activities', 'pessoal'),

    -- trabalho_de_casa tags
    ('Cleaning', 'House cleaning', 'trabalho_de_casa'),
    ('Cooking', 'Cooking meals', 'trabalho_de_casa'),
    ('Laundry', 'Laundry tasks', 'trabalho_de_casa'),
    ('Dishes', 'Washing dishes', 'trabalho_de_casa'),
    ('Organizing', 'Organizing spaces', 'trabalho_de_casa'),

    -- outros tags
    ('Off', 'Time off / break', 'outros'),
    ('Travel', 'Traveling', 'outros'),
    ('Driving', 'Driving / commuting', 'outros'),
    ('Doctor', 'Medical appointments', 'outros'),
    ('Dentist', 'Dental appointments', 'outros'),
    ('Barber', 'Haircut / grooming', 'outros'),
    ('Bank', 'Banking errands', 'outros'),
    ('Market', 'Grocery shopping', 'outros')
) AS t(tag_name, description, category_name)
INNER JOIN aion_api.tag_categories c
  ON c.user_id = u.user_id
  AND c.name = category_name || '_' || u.user_id
  AND c.deleted_at IS NULL
ON CONFLICT (user_id, lower(name)) WHERE deleted_at IS NULL DO NOTHING;
