-- ============================================================================
-- Category Seed Data - User ID 1 (Test User)
-- Idempotent via WHERE NOT EXISTS to match unique index on (user_id, lower(name))
-- ============================================================================

-- saude_fisica
INSERT INTO aion_api.tag_categories (user_id, name, description, color_hex, icon, created_at, updated_at, deleted_at)
SELECT 1, 'saude_fisica', 'Atividades físicas e condicionamento', '#E94F37', 'dumbbell', '2025-01-01 08:00:00', '2025-01-01 08:00:00', NULL
WHERE NOT EXISTS (
  SELECT 1 FROM aion_api.tag_categories WHERE user_id = 1 AND lower(name) = lower('saude_fisica') AND deleted_at IS NULL
);

-- meditacao
INSERT INTO aion_api.tag_categories (user_id, name, description, color_hex, icon, created_at, updated_at, deleted_at)
SELECT 1, 'meditacao', 'Práticas de meditação e atenção plena', '#9C27B0', 'spa', '2025-01-01 08:05:00', '2025-01-01 08:05:00', NULL
WHERE NOT EXISTS (
  SELECT 1 FROM aion_api.tag_categories WHERE user_id = 1 AND lower(name) = lower('meditacao') AND deleted_at IS NULL
);

-- saude_mental
INSERT INTO aion_api.tag_categories (user_id, name, description, color_hex, icon, created_at, updated_at, deleted_at)
SELECT 1, 'saude_mental', 'Saúde mental, planejamento e bem-estar', '#F8B400', 'brain', '2025-01-01 08:10:00', '2025-01-01 08:10:00', NULL
WHERE NOT EXISTS (
  SELECT 1 FROM aion_api.tag_categories WHERE user_id = 1 AND lower(name) = lower('saude_mental') AND deleted_at IS NULL
);

-- estudo_trabalho
INSERT INTO aion_api.tag_categories (user_id, name, description, color_hex, icon, created_at, updated_at, deleted_at)
SELECT 1, 'estudo_trabalho', 'Estudo, trabalho e desenvolvimento', '#1976D2', 'briefcase', '2025-01-01 08:15:00', '2025-01-01 08:15:00', NULL
WHERE NOT EXISTS (
  SELECT 1 FROM aion_api.tag_categories WHERE user_id = 1 AND lower(name) = lower('estudo_trabalho') AND deleted_at IS NULL
);

-- idiomas
INSERT INTO aion_api.tag_categories (user_id, name, description, color_hex, icon, created_at, updated_at, deleted_at)
SELECT 1, 'idiomas', 'Atividades de aprendizado de idiomas', '#00ACC1', 'globe', '2025-01-01 08:20:00', '2025-01-01 08:20:00', NULL
WHERE NOT EXISTS (
  SELECT 1 FROM aion_api.tag_categories WHERE user_id = 1 AND lower(name) = lower('idiomas') AND deleted_at IS NULL
);

-- pessoal
INSERT INTO aion_api.tag_categories (user_id, name, description, color_hex, icon, created_at, updated_at, deleted_at)
SELECT 1, 'pessoal', 'Assuntos pessoais, lazer e tempo livre', '#FF6F00', 'user', '2025-01-01 08:25:00', '2025-01-01 08:25:00', NULL
WHERE NOT EXISTS (
  SELECT 1 FROM aion_api.tag_categories WHERE user_id = 1 AND lower(name) = lower('pessoal') AND deleted_at IS NULL
);

-- trabalho_de_casa
INSERT INTO aion_api.tag_categories (user_id, name, description, color_hex, icon, created_at, updated_at, deleted_at)
SELECT 1, 'trabalho_de_casa', 'Tarefas domésticas e manutenção da casa', '#388E3C', 'home', '2025-01-01 08:30:00', '2025-01-01 08:30:00', NULL
WHERE NOT EXISTS (
  SELECT 1 FROM aion_api.tag_categories WHERE user_id = 1 AND lower(name) = lower('trabalho_de_casa') AND deleted_at IS NULL
);

-- outros
INSERT INTO aion_api.tag_categories (user_id, name, description, color_hex, icon, created_at, updated_at, deleted_at)
SELECT 1, 'outros', 'Atividades diversas / off / viagens', '#616161', 'ellipsis-h', '2025-01-01 08:35:00', '2025-01-01 08:35:00', NULL
WHERE NOT EXISTS (
  SELECT 1 FROM aion_api.tag_categories WHERE user_id = 1 AND lower(name) = lower('outros') AND deleted_at IS NULL
);
