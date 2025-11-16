-- ============================================================================
-- Tags Seed Data - User ID 1 (Test User)
-- Idempotent via WHERE NOT EXISTS to match unique index on (user_id, lower(name))
-- ============================================================================

-- saude_fisica
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Stretching', (SELECT category_id FROM aion_api.tag_categories WHERE name='saude_fisica' AND user_id=1 LIMIT 1), 'Stretching and flexibility', '2025-01-01 09:00:00', '2025-01-01 09:00:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Stretching') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Push-up', (SELECT category_id FROM aion_api.tag_categories WHERE name='saude_fisica' AND user_id=1 LIMIT 1), 'Push-up exercise', '2025-01-01 09:01:00', '2025-01-01 09:01:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Push-up') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'ABS', (SELECT category_id FROM aion_api.tag_categories WHERE name='saude_fisica' AND user_id=1 LIMIT 1), 'Abdominal exercises', '2025-01-01 09:02:00', '2025-01-01 09:02:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('ABS') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Run', (SELECT category_id FROM aion_api.tag_categories WHERE name='saude_fisica' AND user_id=1 LIMIT 1), 'Running sessions', '2025-01-01 09:03:00', '2025-01-01 09:03:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Run') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Pull-up', (SELECT category_id FROM aion_api.tag_categories WHERE name='saude_fisica' AND user_id=1 LIMIT 1), 'Pull-up / upper body', '2025-01-01 09:04:00', '2025-01-01 09:04:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Pull-up') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Sunbathe', (SELECT category_id FROM aion_api.tag_categories WHERE name='saude_fisica' AND user_id=1 LIMIT 1), 'Sun exposure / relaxation', '2025-01-01 09:05:00', '2025-01-01 09:05:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Sunbathe') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Gym', (SELECT category_id FROM aion_api.tag_categories WHERE name='saude_fisica' AND user_id=1 LIMIT 1), 'Gym workouts', '2025-01-01 09:06:00', '2025-01-01 09:06:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Gym') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Walking', (SELECT category_id FROM aion_api.tag_categories WHERE name='saude_fisica' AND user_id=1 LIMIT 1), 'Walking sessions', '2025-01-01 09:07:00', '2025-01-01 09:07:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Walking') AND deleted_at IS NULL);

-- meditacao
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Meditation', (SELECT category_id FROM aion_api.tag_categories WHERE name='meditacao' AND user_id=1 LIMIT 1), 'Meditation practice', '2025-01-01 09:10:00', '2025-01-01 09:10:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Meditation') AND deleted_at IS NULL);

-- saude_mental
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Planning', (SELECT category_id FROM aion_api.tag_categories WHERE name='saude_mental' AND user_id=1 LIMIT 1), 'Planning and organizing', '2025-01-01 09:20:00', '2025-01-01 09:20:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Planning') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Diary', (SELECT category_id FROM aion_api.tag_categories WHERE name='saude_mental' AND user_id=1 LIMIT 1), 'Personal diary entries', '2025-01-01 09:21:00', '2025-01-01 09:21:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Diary') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Reading', (SELECT category_id FROM aion_api.tag_categories WHERE name='saude_mental' AND user_id=1 LIMIT 1), 'Reading and books', '2025-01-01 09:22:00', '2025-01-01 09:22:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Reading') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'TheNews', (SELECT category_id FROM aion_api.tag_categories WHERE name='saude_mental' AND user_id=1 LIMIT 1), 'Reading the news', '2025-01-01 09:23:00', '2025-01-01 09:23:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('TheNews') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Emails', (SELECT category_id FROM aion_api.tag_categories WHERE name='saude_mental' AND user_id=1 LIMIT 1), 'Email processing', '2025-01-01 09:24:00', '2025-01-01 09:24:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Emails') AND deleted_at IS NULL);

-- estudo_trabalho
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Dev', (SELECT category_id FROM aion_api.tag_categories WHERE name='estudo_trabalho' AND user_id=1 LIMIT 1), 'Development work', '2025-01-01 09:30:00', '2025-01-01 09:30:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Dev') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'College', (SELECT category_id FROM aion_api.tag_categories WHERE name='estudo_trabalho' AND user_id=1 LIMIT 1), 'College / study', '2025-01-01 09:31:00', '2025-01-01 09:31:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('College') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Golang', (SELECT category_id FROM aion_api.tag_categories WHERE name='estudo_trabalho' AND user_id=1 LIMIT 1), 'Golang practice', '2025-01-01 09:32:00', '2025-01-01 09:32:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Golang') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Notion', (SELECT category_id FROM aion_api.tag_categories WHERE name='estudo_trabalho' AND user_id=1 LIMIT 1), 'Notion usage', '2025-01-01 09:33:00', '2025-01-01 09:33:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Notion') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'GPT', (SELECT category_id FROM aion_api.tag_categories WHERE name='estudo_trabalho' AND user_id=1 LIMIT 1), 'GPT interactions', '2025-01-01 09:34:00', '2025-01-01 09:34:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('GPT') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Full Cycle', (SELECT category_id FROM aion_api.tag_categories WHERE name='estudo_trabalho' AND user_id=1 LIMIT 1), 'Full Cycle learning', '2025-01-01 09:35:00', '2025-01-01 09:35:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Full Cycle') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'FreeCodeCamp', (SELECT category_id FROM aion_api.tag_categories WHERE name='estudo_trabalho' AND user_id=1 LIMIT 1), 'FreeCodeCamp tasks', '2025-01-01 09:36:00', '2025-01-01 09:36:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('FreeCodeCamp') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Coursera', (SELECT category_id FROM aion_api.tag_categories WHERE name='estudo_trabalho' AND user_id=1 LIMIT 1), 'Coursera courses', '2025-01-01 09:37:00', '2025-01-01 09:37:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Coursera') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Course', (SELECT category_id FROM aion_api.tag_categories WHERE name='estudo_trabalho' AND user_id=1 LIMIT 1), 'Generic course activity', '2025-01-01 09:38:00', '2025-01-01 09:38:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Course') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Aion', (SELECT category_id FROM aion_api.tag_categories WHERE name='estudo_trabalho' AND user_id=1 LIMIT 1), 'Aion related work', '2025-01-01 09:39:00', '2025-01-01 09:39:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Aion') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Work', (SELECT category_id FROM aion_api.tag_categories WHERE name='estudo_trabalho' AND user_id=1 LIMIT 1), 'Work tasks', '2025-01-01 09:40:00', '2025-01-01 09:40:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Work') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'RD', (SELECT category_id FROM aion_api.tag_categories WHERE name='estudo_trabalho' AND user_id=1 LIMIT 1), 'Research & development', '2025-01-01 09:41:00', '2025-01-01 09:41:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('RD') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Podcast', (SELECT category_id FROM aion_api.tag_categories WHERE name='estudo_trabalho' AND user_id=1 LIMIT 1), 'Podcast listening', '2025-01-01 09:42:00', '2025-01-01 09:42:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Podcast') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'AudioBook', (SELECT category_id FROM aion_api.tag_categories WHERE name='estudo_trabalho' AND user_id=1 LIMIT 1), 'AudioBook listening', '2025-01-01 09:43:00', '2025-01-01 09:43:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('AudioBook') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Finance', (SELECT category_id FROM aion_api.tag_categories WHERE name='estudo_trabalho' AND user_id=1 LIMIT 1), 'Finance related course', '2025-01-01 09:44:00', '2025-01-01 09:44:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Finance') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Interview', (SELECT category_id FROM aion_api.tag_categories WHERE name='estudo_trabalho' AND user_id=1 LIMIT 1), 'Interview prep', '2025-01-01 09:45:00', '2025-01-01 09:45:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Interview') AND deleted_at IS NULL);

-- idiomas
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'English', (SELECT category_id FROM aion_api.tag_categories WHERE name='idiomas' AND user_id=1 LIMIT 1), 'English study', '2025-01-01 09:50:00', '2025-01-01 09:50:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('English') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Spanish', (SELECT category_id FROM aion_api.tag_categories WHERE name='idiomas' AND user_id=1 LIMIT 1), 'Spanish study', '2025-01-01 09:51:00', '2025-01-01 09:51:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Spanish') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'German', (SELECT category_id FROM aion_api.tag_categories WHERE name='idiomas' AND user_id=1 LIMIT 1), 'German study', '2025-01-01 09:52:00', '2025-01-01 09:52:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('German') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'French', (SELECT category_id FROM aion_api.tag_categories WHERE name='idiomas' AND user_id=1 LIMIT 1), 'French study', '2025-01-01 09:53:00', '2025-01-01 09:53:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('French') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Chinese', (SELECT category_id FROM aion_api.tag_categories WHERE name='idiomas' AND user_id=1 LIMIT 1), 'Chinese study', '2025-01-01 09:54:00', '2025-01-01 09:54:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Chinese') AND deleted_at IS NULL);

-- pessoal
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'My', (SELECT category_id FROM aion_api.tag_categories WHERE name='pessoal' AND user_id=1 LIMIT 1), 'Personal tasks', '2025-01-01 10:00:00', '2025-01-01 10:00:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('My') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Travel', (SELECT category_id FROM aion_api.tag_categories WHERE name='pessoal' AND user_id=1 LIMIT 1), 'Travel planning', '2025-01-01 10:01:00', '2025-01-01 10:01:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Travel') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Movie', (SELECT category_id FROM aion_api.tag_categories WHERE name='pessoal' AND user_id=1 LIMIT 1), 'Watching movies', '2025-01-01 10:02:00', '2025-01-01 10:02:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Movie') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Series', (SELECT category_id FROM aion_api.tag_categories WHERE name='pessoal' AND user_id=1 LIMIT 1), 'Watching series', '2025-01-01 10:03:00', '2025-01-01 10:03:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Series') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Game', (SELECT category_id FROM aion_api.tag_categories WHERE name='pessoal' AND user_id=1 LIMIT 1), 'Gaming', '2025-01-01 10:04:00', '2025-01-01 10:04:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Game') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Beach', (SELECT category_id FROM aion_api.tag_categories WHERE name='pessoal' AND user_id=1 LIMIT 1), 'Beach time', '2025-01-01 10:05:00', '2025-01-01 10:05:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Beach') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Hanging Out', (SELECT category_id FROM aion_api.tag_categories WHERE name='pessoal' AND user_id=1 LIMIT 1), 'Hanging out', '2025-01-01 10:06:00', '2025-01-01 10:06:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Hanging Out') AND deleted_at IS NULL);

-- trabalho_de_casa
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Housework', (SELECT category_id FROM aion_api.tag_categories WHERE name='trabalho_de_casa' AND user_id=1 LIMIT 1), 'House chores', '2025-01-01 10:10:00', '2025-01-01 10:10:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Housework') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Supermarket', (SELECT category_id FROM aion_api.tag_categories WHERE name='trabalho_de_casa' AND user_id=1 LIMIT 1), 'Grocery shopping', '2025-01-01 10:11:00', '2025-01-01 10:11:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Supermarket') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Doctor', (SELECT category_id FROM aion_api.tag_categories WHERE name='trabalho_de_casa' AND user_id=1 LIMIT 1), 'Doctor visits', '2025-01-01 10:12:00', '2025-01-01 10:12:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Doctor') AND deleted_at IS NULL);

-- outros
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'OFF', (SELECT category_id FROM aion_api.tag_categories WHERE name='outros' AND user_id=1 LIMIT 1), 'Days off / rest', '2025-01-01 10:20:00', '2025-01-01 10:20:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('OFF') AND deleted_at IS NULL);
INSERT INTO aion_api.tags (user_id, name, category_id, description, created_at, updated_at, deleted_at)
SELECT 1, 'Travelling', (SELECT category_id FROM aion_api.tag_categories WHERE name='outros' AND user_id=1 LIMIT 1), 'Travel activities', '2025-01-01 10:21:00', '2025-01-01 10:21:00', NULL
WHERE NOT EXISTS (SELECT 1 FROM aion_api.tags WHERE user_id=1 AND lower(name)=lower('Travelling') AND deleted_at IS NULL);

-- NOTE: Other users' tags have been removed to match the request.
