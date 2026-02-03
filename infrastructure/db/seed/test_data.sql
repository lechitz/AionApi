-- ============================================================================
-- SEED: Test Data - Complete test profile
-- Description: Creates a test user with categories, tags, and realistic records
-- Usage: psql -h localhost -p 5432 -U postgres -d aion_dev -f test_data.sql
-- Cleanup: DELETE FROM aion_api.users WHERE username = 'testuser';
-- ============================================================================

BEGIN;

-- ============================================================================
-- 0. CLEANUP EXISTING TEST DATA (idempotent)
-- ============================================================================
DELETE FROM aion_api.records WHERE user_id = 999;
DELETE FROM aion_api.tags WHERE user_id = 999;
DELETE FROM aion_api.categories WHERE user_id = 999;
DELETE FROM aion_api.user_roles WHERE user_id = 999;
DELETE FROM aion_api.users WHERE user_id = 999;

-- ============================================================================
-- 1. CREATE TEST USER
-- ============================================================================
-- Password: Test@123 (hashed with bcrypt, cost 10 - DefaultCost)
INSERT INTO aion_api.users (
    user_id, name, username, password, email, 
    locale, timezone, location, bio, created_at, updated_at
) VALUES (
    999,
    'Usuário Teste',
    'testuser',
    '$2a$10$HR0PGkRt7bnSmR2cz5TLSus2vFFInrkaMxDgsWJxftphgjEAC3dmC',
    'test@aion.local',
    'pt-BR',
    'America/Sao_Paulo',
    'São Paulo, Brasil',
    'Perfil de teste para demonstração e desenvolvimento.',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);

-- Assign user role (role_id=3 is 'user' from migration 000002)
INSERT INTO aion_api.user_roles (user_id, role_id, assigned_at)
VALUES (999, 3, CURRENT_TIMESTAMP);


-- ============================================================================
-- 2. CREATE CATEGORIES (IDs: 101-108)
-- ============================================================================
INSERT INTO aion_api.categories (category_id, user_id, name, description, color_hex, icon, created_at, updated_at) VALUES
(101, 999, 'Saúde Física', 'hábitos relacionados a saúde física.', '#056609', 'health/heart.svg', '2026-02-03 01:19:03.374138', '2026-02-03 01:57:29.035405'),
(102, 999, 'Saúde Mental', 'Hábitos relacionados a saúde mental', '#8f8424', 'mind/brain.svg', '2026-02-03 01:43:23.041031', '2026-02-03 01:58:20.700668'),
(103, 999, 'Estudos', 'cursos, projetos, aulas..', '#E91E63', 'study/book.svg', '2026-02-03 01:44:23.055436', '2026-02-03 01:44:23.055436'),
(104, 999, 'Trabalho', 'Trabalhos profissionais, empresas,..', '#681179', 'work/briefcase.svg', '2026-02-03 01:46:31.325365', '2026-02-03 01:46:31.325365'),
(105, 999, 'Saúde Financeira', 'cartões, parcelamentos, etc..', '#147fd7', 'finance/wallet.svg', '2026-02-03 01:47:29.612804', '2026-02-03 01:47:29.612804'),
(106, 999, 'Saúde Alimentar', 'diferentes refeições e bebidas', '#e06d10', 'food/leaf.svg', '2026-02-03 01:49:53.387097', '2026-02-03 01:58:39.272788'),
(107, 999, 'Projetos', 'Categoria para projetos de estudos e práticas', '#052548', 'personal/user.svg', '2026-02-03 02:32:12.751278', '2026-02-03 02:32:12.751278'),
(108, 999, 'Pessoal', 'Momentos pessoais', '#5c707a', 'personal/user.svg', '2026-02-03 02:36:51.249552', '2026-02-03 02:36:51.249552');


-- ============================================================================
-- 3. CREATE TAGS (IDs: 201-222)
-- ============================================================================
INSERT INTO aion_api.tags (tag_id, user_id, name, category_id, description, icon, created_at, updated_at) VALUES
(201, 999, 'Água', 106, 'ingestão de copos de água. A quantidade varia conforme descrição de cada registro.', '💧', '2026-02-03 01:51:14.003933', '2026-02-03 01:51:14.003933'),
(202, 999, 'Cerveja', 106, 'ingestão de cerveja, quantidade definida em cada registro.', '🤝', '2026-02-03 01:52:13.189578', '2026-02-03 01:52:13.189578'),
(203, 999, 'Café da manhã', 106, 'Para registrar refeições de café da manhã. Pratos e porções, ficam a critério dos registros.', '☕', '2026-02-03 01:53:32.521603', '2026-02-03 01:53:32.521603'),
(204, 999, 'Almoço', 106, 'refeições próximas ao horário do almoço ou assim consideradas. Cardápio em cada registro.', '🥗', '2026-02-03 01:54:08.312942', '2026-02-03 01:54:08.312942'),
(205, 999, 'Jantar', 106, 'separado para refeições a noite.', '🥗', '2026-02-03 01:55:45.346861', '2026-02-03 01:55:45.346861'),
(206, 999, 'Caminhar', 101, 'Caminhadas no quarteirão', '🏃', '2026-02-03 01:56:45.139013', '2026-02-03 01:56:45.139013'),
(207, 999, 'Flexão', 101, '', '🏋️', '2026-02-03 01:57:02.524692', '2026-02-03 01:57:02.524692'),
(208, 999, 'golang', 103, 'estudos, livros, cursos, projetos', '📚', '2026-02-03 01:59:52.109558', '2026-02-03 01:59:52.109558'),
(209, 999, 'AWS', 103, 'cursos, projetos, estudos', '📉', '2026-02-03 02:00:13.068235', '2026-02-03 02:00:13.068235'),
(210, 999, 'python', 103, 'cursos, projetos relacionados a python', '🛠️', '2026-02-03 02:01:16.588112', '2026-02-03 02:01:16.588112'),
(211, 999, 'Débito', 105, 'gastos no débito', '✈️', '2026-02-03 02:02:18.038721', '2026-02-03 02:02:18.038721'),
(212, 999, 'Crédito', 105, 'gastos no cartão de crédito, parcelados e a vista. Variação da bandeira ficará em cada registro.', '📚', '2026-02-03 02:03:01.030963', '2026-02-03 02:03:01.030963'),
(213, 999, 'react', 103, 'estudos de front, em react', '📖', '2026-02-03 02:04:50.913264', '2026-02-03 02:04:50.913264'),
(214, 999, 'Café Preto', 106, 'registros de xícaras de café.', '📚', '2026-02-03 02:05:25.240279', '2026-02-03 02:05:25.240279'),
(215, 999, 'BRQ', 104, 'todos os trabalhos  em projetos com a BRQ', '💼', '2026-02-03 02:29:59.923310', '2026-02-03 02:29:59.923310'),
(216, 999, 'Aion', 107, 'Estudo e desenvolvimento do Projeto Aion', '🧠', '2026-02-03 02:33:07.164218', '2026-02-03 02:33:07.164218'),
(217, 999, 'Kafka-stream', 107, 'Projeto de estudos de Kafka', '📊', '2026-02-03 02:33:36.044320', '2026-02-03 02:33:36.044320'),
(218, 999, 'Meditação', 102, 'sessões de meditação', '🧘', '2026-02-03 02:35:08.082966', '2026-02-03 02:35:08.082966'),
(219, 999, 'Auto-Hipnose', 102, 'Sessões de auto-hipnose', '🎯', '2026-02-03 02:35:32.692669', '2026-02-03 02:35:32.692669'),
(220, 999, 'Relax', 108, 'Tag para quando eu quiser tirar um tempinho pra mim, curtir minha paz, escutar uma música na sala..', '✈️', '2026-02-03 02:39:47.462553', '2026-02-03 02:39:47.462553'),
(221, 999, 'Diário Profissional', 104, 'Tag para registros que irei escrever para detalhar partes da jornada diária.', '📝', '2026-02-03 03:15:57.041377', '2026-02-03 03:15:57.041377'),
(222, 999, 'Diário Pessoal', 108, 'registro pessoais', '📖', '2026-02-03 03:16:59.060696', '2026-02-03 03:16:59.060696');


-- ============================================================================
-- 4. CREATE RECORDS (50 records across last 3 days)
-- ============================================================================
-- Rules:
-- - Água (201) and Café Preto (214): 5-10 minutes duration
-- - Events must be before current time or still in progress
-- - 3 intentional overlaps for multi-row rendering test
-- - All marked with source='test_seed' for easy cleanup
-- ============================================================================

-- Helper function to insert records only if event_time is in the past
DO $$
DECLARE
    v_base_date DATE := CURRENT_DATE - INTERVAL '2 days';
    v_day1 DATE := v_base_date;
    v_day2 DATE := v_base_date + INTERVAL '1 day';
    v_day3 DATE := v_base_date + INTERVAL '2 days';
    v_now TIMESTAMPTZ := NOW();
    v_tz TEXT := 'America/Sao_Paulo';
BEGIN

-- Helper function to create timestamptz from date and time
-- Uses the format: (date + time)::timestamp AT TIME ZONE 'America/Sao_Paulo'

-- ========================================
-- DAY 1 (2 days ago) - 16 records
-- ========================================

-- Morning routine
IF ((v_day1 + TIME '07:00:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 201, '1 copo de água ao acordar', 300, (v_day1 + TIME '07:00:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day1 + TIME '07:30:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 203, 'Pão integral com queijo e café', 1200, (v_day1 + TIME '07:30:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day1 + TIME '08:00:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 214, 'Café preto médio', 420, (v_day1 + TIME '08:00:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

-- Work morning
IF ((v_day1 + TIME '09:00:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 215, 'Daily meeting + code review', 7200, (v_day1 + TIME '09:00:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day1 + TIME '10:30:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 201, '1 copo de água', 240, (v_day1 + TIME '10:30:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day1 + TIME '11:00:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 214, 'Café preto', 360, (v_day1 + TIME '11:00:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

-- Lunch
IF ((v_day1 + TIME '12:30:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 204, 'Arroz, feijão, frango grelhado e salada', 1800, v_day1 + TIME '12:30:00', 'test_seed', 'published');
END IF;

-- Afternoon work
IF ((v_day1 + TIME '14:00:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 215, 'Implementação de features + testes', 10800, (v_day1 + TIME '14:00:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day1 + TIME '15:30:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 201, '1 copo de água', 300, (v_day1 + TIME '15:30:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

-- Evening
IF ((v_day1 + TIME '18:00:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 206, 'Caminhada no quarteirão', 1800, (v_day1 + TIME '18:00:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day1 + TIME '19:00:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 205, 'Sopa de legumes com torradas', 1500, (v_day1 + TIME '19:00:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day1 + TIME '20:00:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 208, 'Estudando Go: goroutines e channels', 5400, (v_day1 + TIME '20:00:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day1 + TIME '21:30:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 201, '1 copo de água', 240, (v_day1 + TIME '21:30:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day1 + TIME '22:00:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 218, 'Meditação guiada 20 minutos', 1200, (v_day1 + TIME '22:00:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day1 + TIME '22:30:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 222, 'Diário: Reflexões do dia', 900, (v_day1 + TIME '22:30:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day1 + TIME '23:00:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 220, 'Assistindo série favorita', 3600, (v_day1 + TIME '23:00:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;


-- ========================================
-- DAY 2 (Yesterday) - 18 records with 1 overlap
-- ========================================

IF ((v_day2 + TIME '06:30:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 201, '1 copo de água ao acordar', 300, (v_day2 + TIME '06:30:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day2 + TIME '07:00:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 207, '3 séries de 15 flexões', 600, (v_day2 + TIME '07:00:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day2 + TIME '07:30:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 203, 'Tapioca com queijo e presunto', 1200, (v_day2 + TIME '07:30:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day2 + TIME '08:00:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 214, 'Café preto', 480, (v_day2 + TIME '08:00:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

-- Work with overlap (personal study in parallel)
IF ((v_day2 + TIME '09:00:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 215, 'Sprint planning meeting', 5400, (v_day2 + TIME '09:00:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

-- OVERLAP 1: Study during break in meeting
IF ((v_day2 + TIME '09:45:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 213, 'Lendo artigo sobre React hooks', 1800, (v_day2 + TIME '09:45:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day2 + TIME '10:45:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 201, '1 copo de água', 300, (v_day2 + TIME '10:45:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day2 + TIME '11:15:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 214, 'Café preto', 420, (v_day2 + TIME '11:15:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day2 + TIME '12:00:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 204, 'Macarrão com molho branco e salada', 1800, (v_day2 + TIME '12:00:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day2 + TIME '13:30:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 215, 'Code review e pair programming', 9000, (v_day2 + TIME '13:30:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day2 + TIME '16:00:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 201, '1 copo de água', 240, (v_day2 + TIME '16:00:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day2 + TIME '16:30:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 211, 'Supermercado: R$ 87,50', 0, v_day2 + TIME '16:30:00', 'test_seed', 'published');
END IF;

IF ((v_day2 + TIME '18:30:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 205, 'Frango assado com batatas', 2100, (v_day2 + TIME '18:30:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day2 + TIME '20:00:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 216, 'Projeto Aion: Ajustes na timeline', 7200, (v_day2 + TIME '20:00:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day2 + TIME '21:00:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 201, '1 copo de água', 300, (v_day2 + TIME '21:00:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day2 + TIME '22:00:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 219, 'Auto-hipnose para relaxamento', 1800, (v_day2 + TIME '22:00:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day2 + TIME '22:45:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 221, 'Diário: Avanços no projeto BRQ', 1200, (v_day2 + TIME '22:45:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day2 + TIME '23:30:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 220, 'Lendo livro antes de dormir', 1800, (v_day2 + TIME '23:30:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;


-- ========================================
-- DAY 3 (Today) - 16 records with 2 overlaps
-- ========================================

IF ((v_day3 + TIME '06:45:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 201, '1 copo de água ao acordar', 300, (v_day3 + TIME '06:45:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day3 + TIME '07:15:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 206, 'Caminhada matinal 30 min', 1800, (v_day3 + TIME '07:15:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day3 + TIME '08:00:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 203, 'Ovos mexidos com pão integral', 1200, (v_day3 + TIME '08:00:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day3 + TIME '08:30:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 214, 'Café preto grande', 540, (v_day3 + TIME '08:30:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

-- Morning overlap: work + personal project
IF ((v_day3 + TIME '09:00:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 215, 'Desenvolvimento de API REST', 10800, (v_day3 + TIME '09:00:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

-- OVERLAP 2: Quick personal task during work
IF ((v_day3 + TIME '10:00:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 212, 'Pagamento fatura cartão: R$ 450,00', 0, v_day3 + TIME '10:00:00', 'test_seed', 'published');
END IF;

IF ((v_day3 + TIME '10:30:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 201, '1 copo de água', 300, (v_day3 + TIME '10:30:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day3 + TIME '11:30:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 214, 'Café preto', 360, (v_day3 + TIME '11:30:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day3 + TIME '12:15:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 204, 'Prato feito: arroz, carne e legumes', 1800, v_day3 + TIME '12:15:00', 'test_seed', 'published');
END IF;

-- Afternoon study session with multiple topics
IF ((v_day3 + TIME '14:00:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 210, 'Python: Estudando FastAPI', 5400, (v_day3 + TIME '14:00:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

-- OVERLAP 3: AWS study starts before Python session ends
IF ((v_day3 + TIME '14:45:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 209, 'AWS: Lambda e API Gateway', 3600, (v_day3 + TIME '14:45:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day3 + TIME '16:00:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 201, '1 copo de água', 240, (v_day3 + TIME '16:00:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day3 + TIME '17:00:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 217, 'Kafka Streams: Processamento em tempo real', 7200, (v_day3 + TIME '17:00:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day3 + TIME '19:00:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 202, 'Cerveja com amigos - 3 long necks', 5400, (v_day3 + TIME '19:00:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day3 + TIME '20:30:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 205, 'Pizza margherita', 1800, (v_day3 + TIME '20:30:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

IF ((v_day3 + TIME '22:00:00')::timestamp AT TIME ZONE v_tz) < v_now THEN
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, source, status)
    VALUES (999, 222, 'Diário: Reflexões sobre aprendizado', 1200, (v_day3 + TIME '22:00:00')::timestamp AT TIME ZONE v_tz, 'test_seed', 'published');
END IF;

END $$;

COMMIT;

-- ============================================================================
-- CLEANUP INSTRUCTIONS
-- ============================================================================
-- To remove all test data:
-- DELETE FROM aion_api.records WHERE user_id = 999;
-- DELETE FROM aion_api.tags WHERE user_id = 999;
-- DELETE FROM aion_api.categories WHERE user_id = 999;
-- DELETE FROM aion_api.user_roles WHERE user_id = 999;
-- DELETE FROM aion_api.users WHERE user_id = 999;
-- ============================================================================
