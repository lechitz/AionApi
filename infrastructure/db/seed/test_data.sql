-- ============================================================================
-- SEED: Full Realistic Demo Profile (3 months)
-- Description:
--   - Creates one realistic demo user (testuser)
--   - Merges legacy taxonomy (old categories/tags) with new dashboard taxonomy
--   - Seeds metric_definitions and goal_templates for dashboardSnapshot
--   - Generates ~90 days with high daily volume (~50-60 records/day)
--   - Ensures timeline/reports have end time (duration_seconds) where applicable
-- Usage:
--   make seed-test
-- ============================================================================

BEGIN;

-- --------------------------------------------------------------------------
-- 0) Cleanup demo user dataset (idempotent)
-- --------------------------------------------------------------------------
DELETE FROM aion_api.goal_instances WHERE user_id = 999;
DELETE FROM aion_api.goal_templates WHERE user_id = 999;
DELETE FROM aion_api.metric_definitions WHERE user_id = 999;
DELETE FROM aion_api.records WHERE user_id = 999;
DELETE FROM aion_api.tags WHERE user_id = 999;
DELETE FROM aion_api.categories WHERE user_id = 999;
DELETE FROM aion_api.user_roles WHERE user_id = 999;
DELETE FROM aion_api.users WHERE user_id = 999;

-- --------------------------------------------------------------------------
-- 1) Demo user + role
-- --------------------------------------------------------------------------
-- Password: Test@123
INSERT INTO aion_api.users (
    user_id, name, username, password, email, locale, timezone, location, bio, created_at, updated_at
) VALUES (
    999,
    'Demo User',
    'testuser',
    '$2a$10$HR0PGkRt7bnSmR2cz5TLSus2vFFInrkaMxDgsWJxftphgjEAC3dmC',
    'test@aion.local',
    'pt-BR',
    'America/Sao_Paulo',
    'Sao Paulo, Brasil',
    'Perfil demo: 3 meses de uso realista (50-60 registros/dia).',
    NOW(),
    NOW()
);

INSERT INTO aion_api.user_roles (user_id, role_id, assigned_at)
SELECT 999, r.role_id, NOW()
FROM aion_api.roles r
WHERE r.name = 'user'
LIMIT 1;

-- --------------------------------------------------------------------------
-- 2) Categories (legacy + new dashboard taxonomy)
-- --------------------------------------------------------------------------
INSERT INTO aion_api.categories (user_id, name, description, color_hex, icon, created_at, updated_at)
VALUES
    (999, 'Saude Fisica', 'Habitos relacionados a saude fisica', '#056609', 'health/heart.svg', NOW(), NOW()),
    (999, 'Saude Mental', 'Habitos relacionados a saude mental', '#8F8424', 'mind/brain.svg', NOW(), NOW()),
    (999, 'Estudos', 'Cursos, projetos, aulas e pratica', '#E91E63', 'study/book.svg', NOW(), NOW()),
    (999, 'Trabalho', 'Rotina profissional e sessoes de trabalho', '#681179', 'work/briefcase.svg', NOW(), NOW()),
    (999, 'Saude Financeira', 'Gastos e registros financeiros', '#147FD7', 'finance/wallet.svg', NOW(), NOW()),
    (999, 'Refeicao', 'Refeicoes e bebidas do dia', '#E06D10', 'food/leaf.svg', NOW(), NOW()),
    (999, 'Hidratacao', 'Consumo de agua e bebidas', '#00ACC1', 'health/water.svg', NOW(), NOW()),
    (999, 'Projetos', 'Projetos de estudo e pratica', '#052548', 'personal/user.svg', NOW(), NOW()),
    (999, 'Pessoal', 'Registros pessoais e reflexoes', '#5C707A', 'personal/user.svg', NOW(), NOW()),
    (999, 'Sono', 'Janela de sono: dormir e acordar', '#3949AB', 'health/moon.svg', NOW(), NOW()),
    (999, 'Humor', 'Escala emocional de muito triste a muito feliz', '#F06292', 'mind/brain.svg', NOW(), NOW()),
    (999, 'Energia', 'Energia diaria em porcentagem (0 a 100)', '#FBC02D', 'mind/brain.svg', NOW(), NOW()),
    (999, 'Objetivos', 'Intencoes e check-ins diarios', '#8E24AA', 'personal/target.svg', NOW(), NOW());

-- --------------------------------------------------------------------------
-- 3) Tags (legacy + new taxonomy)
-- --------------------------------------------------------------------------
INSERT INTO aion_api.tags (user_id, name, category_id, description, icon, created_at, updated_at)
SELECT
    999,
    x.tag_name,
    c.category_id,
    x.tag_description,
    x.tag_icon,
    NOW(),
    NOW()
FROM (
    VALUES
        -- Legacy tags (preserved)
        ('Hidratacao', 'Agua', 'Ingestao de agua em ml', '💧'),
        ('Hidratacao', 'Cerveja', 'Ingestao de cerveja', '🍺'),
        ('Refeicao', 'Cafe da Manha', 'Registro de cafe da manha', '☕'),
        ('Refeicao', 'Almoco', 'Registro de almoco', '🥗'),
        ('Refeicao', 'Janta', 'Registro de jantar', '🍛'),
        ('Saude Fisica', 'Caminhar', 'Caminhadas', '🏃'),
        ('Saude Fisica', 'Flexao', 'Treino de flexoes', '🏋️'),
        ('Saude Fisica', 'Alongamento', 'Alongamento e flexibilidade', '🤸'),
        ('Saude Fisica', 'ABS', 'Treino abdominal', '🧱'),
        ('Saude Fisica', 'Corrida', 'Sessao de corrida', '🏃‍♂️'),
        ('Saude Fisica', 'Barra', 'Treino de barra fixa', '🦾'),
        ('Saude Fisica', 'Sol', 'Exposicao ao sol', '🌤️'),
        ('Saude Fisica', 'Academia', 'Treino de academia', '🏋️‍♀️'),
        ('Estudos', 'golang', 'Estudos de Go', '📚'),
        ('Estudos', 'AWS', 'Estudos de AWS', '☁️'),
        ('Estudos', 'python', 'Estudos de Python', '🐍'),
        ('Estudos', 'react', 'Estudos de React', '⚛️'),
        ('Saude Financeira', 'Debito', 'Gastos no debito', '💳'),
        ('Saude Financeira', 'Credito', 'Gastos no credito', '💸'),
        ('Refeicao', 'Cafe Preto', 'Registro de xicaras de cafe', '☕'),
        ('Trabalho', 'BRQ', 'Projetos na BRQ', '💼'),
        ('Projetos', 'Aion', 'Desenvolvimento do projeto Aion', '🧠'),
        ('Projetos', 'Kafka-stream', 'Projeto de estudos de Kafka', '📊'),
        ('Saude Mental', 'Meditacao', 'Sessoes de meditacao', '🧘'),
        ('Saude Mental', 'Auto-Hipnose', 'Sessoes de auto-hipnose', '🎯'),
        ('Pessoal', 'Relax', 'Tempo de descanso pessoal', '🛋️'),
        ('Trabalho', 'Diario Profissional', 'Resumo da jornada profissional', '📝'),
        ('Pessoal', 'Diario Pessoal', 'Reflexoes pessoais', '📓'),

        -- New operational tags
        ('Sono', 'Dormir', 'Inicio da janela de sono', '😴'),
        ('Sono', 'Acordar', 'Fim da janela de sono', '🌅'),
        ('Humor', 'Muito Triste', 'Escala de humor 0%', '😢'),
        ('Humor', 'Triste', 'Escala de humor 25%', '😕'),
        ('Humor', 'Neutro', 'Escala de humor 50%', '😐'),
        ('Humor', 'Feliz', 'Escala de humor 75%', '🙂'),
        ('Humor', 'Muito Feliz', 'Escala de humor 100%', '😄'),
        ('Energia', 'Energia %', 'Energia em porcentagem de 0 a 100', '⚡'),
        ('Refeicao', 'Lanche', 'Registro de lanche', '🍎'),
        ('Trabalho', 'Inicio Expediente', 'Inicio do expediente', '🟢'),
        ('Trabalho', 'Saida Almoco', 'Saida para almoco', '🍽️'),
        ('Trabalho', 'Retorno Almoco', 'Retorno do almoco', '🔙'),
        ('Trabalho', 'Fim Expediente', 'Fim do expediente', '🔴'),
        ('Trabalho', 'Sessao de Trabalho', 'Sessao de foco para metrificar horas trabalhadas', '🧩'),
        ('Objetivos', 'Intencao', 'Intencoes diarias', '🎯'),
        ('Objetivos', 'Check-in', 'Check-in de progresso', '✅')
) AS x(category_name, tag_name, tag_description, tag_icon)
JOIN aion_api.categories c
  ON c.user_id = 999
 AND c.name = x.category_name;

-- --------------------------------------------------------------------------
-- 4) Dashboard semantics (metric_definitions + goal_templates)
-- --------------------------------------------------------------------------
INSERT INTO aion_api.metric_definitions (
    user_id, metric_key, display_name, category_id, tag_id,
    value_source, aggregation, unit, goal_default, is_active, created_at, updated_at
)
SELECT
    999,
    md.metric_key,
    md.display_name,
    c.category_id,
    t.tag_id,
    md.value_source,
    md.aggregation,
    md.unit,
    md.goal_default,
    TRUE,
    NOW(),
    NOW()
FROM (
    VALUES
        ('sleep',        'Sono',            'duration_seconds', 'sum',    'seconds', 28800::numeric, 'Sono',       'Dormir'),
        ('work_session', 'Trabalho',        'duration_seconds', 'sum',    'seconds', 28800::numeric, 'Trabalho',   'Sessao de Trabalho'),
        ('water',        'Agua',            'value',            'sum',    'ml',      3000::numeric,  'Hidratacao', 'Agua'),
        ('nutrition',    'Refeicoes',       'count',            'count',  'meals',   4::numeric,     'Refeicao',   'Almoco'),
        ('energy',       'Energia',         'value',            'latest', 'percent', NULL::numeric,  'Energia',    'Energia %'),
        ('mood',         'Humor',           'value',            'latest', 'percent', NULL::numeric,  'Humor',      'Neutro'),
        ('intentions',   'Intencoes',       'count',            'count',  'items',   3::numeric,     'Objetivos',  'Intencao'),
        ('physical',     'Saude Fisica',    'duration_seconds', 'sum',    'seconds', 2700::numeric,  'Saude Fisica','Caminhar')
) AS md(metric_key, display_name, value_source, aggregation, unit, goal_default, category_name, tag_name)
JOIN aion_api.categories c
  ON c.user_id = 999
 AND c.name = md.category_name
JOIN aion_api.tags t
  ON t.user_id = 999
 AND t.name = md.tag_name;

-- Multi-tag metric bindings (extends primary tag_id from metric_definitions)
INSERT INTO aion_api.metric_definition_tag_bindings (user_id, metric_definition_id, tag_id)
SELECT
    999,
    md.id,
    t.tag_id
FROM aion_api.metric_definitions md
JOIN aion_api.tags t
  ON t.user_id = 999
WHERE md.user_id = 999
  AND (
      (md.metric_key = 'mood' AND t.name IN ('Muito Triste', 'Triste', 'Neutro', 'Feliz', 'Muito Feliz'))
      OR
      (md.metric_key = 'sleep' AND t.name IN ('Dormir', 'Acordar'))
      OR
      (md.metric_key = 'work_session' AND t.name IN ('Sessao de Trabalho', 'Inicio Expediente', 'Saida Almoco', 'Retorno Almoco', 'Fim Expediente'))
      OR
      (md.metric_key = 'nutrition' AND t.name IN ('Cafe da Manha', 'Almoco', 'Lanche', 'Janta'))
      OR
      (md.metric_key = 'physical' AND t.name IN ('Caminhar', 'Flexao', 'Alongamento', 'ABS', 'Corrida', 'Barra', 'Sol', 'Academia'))
  );

INSERT INTO aion_api.goal_templates (
    user_id, metric_key, title, target_value, comparison, period, is_active, created_at, updated_at
)
VALUES
    (999, 'sleep',        'Dormir pelo menos 8h',           28800, 'gte', 'day', TRUE, NOW(), NOW()),
    (999, 'work_session', 'Trabalhar pelo menos 8h',        28800, 'gte', 'day', TRUE, NOW(), NOW()),
    (999, 'water',        'Consumir 3000 ml de agua',        3000, 'gte', 'day', TRUE, NOW(), NOW()),
    (999, 'nutrition',    'Registrar 4 refeicoes',              4, 'gte', 'day', TRUE, NOW(), NOW()),
    (999, 'intentions',   'Registrar 3 intencoes',              3, 'gte', 'day', TRUE, NOW(), NOW()),
    (999, 'physical',     'Saude fisica por 45 minutos',     2700, 'gte', 'day', TRUE, NOW(), NOW());

-- --------------------------------------------------------------------------
-- 5) Historical records (about 3 months, 50-60 records/day)
-- --------------------------------------------------------------------------
DO $$
DECLARE
    v_user_id BIGINT := 999;
    v_tz      TEXT := 'America/Sao_Paulo';
    v_day     DATE;
    v_day_idx INT := 0;
    v_target_count INT;
    v_inserted INT;
    v_is_weekend BOOLEAN;
    v_i INT;
    v_ts TIMESTAMPTZ;
    v_duration INT;
    v_value NUMERIC(12,2);
    v_mood_name TEXT;
    v_extra_idx INT;

    t_agua BIGINT; t_cafe_manha BIGINT; t_almoco BIGINT; t_lanche BIGINT; t_janta BIGINT; t_cafe_preto BIGINT;
    t_caminhar BIGINT; t_flexao BIGINT; t_alongamento BIGINT; t_abs BIGINT; t_corrida BIGINT; t_barra BIGINT; t_sol BIGINT; t_academia BIGINT;
    t_golang BIGINT; t_aws BIGINT; t_python BIGINT; t_react BIGINT;
    t_brq BIGINT; t_aion BIGINT; t_kafka BIGINT; t_meditacao BIGINT; t_auto_hipnose BIGINT; t_relax BIGINT;
    t_diario_p BIGINT; t_diario_w BIGINT; t_dormir BIGINT; t_acordar BIGINT; t_energia BIGINT;
    t_muito_triste BIGINT; t_triste BIGINT; t_neutro BIGINT; t_feliz BIGINT; t_muito_feliz BIGINT;
    t_inicio_exp BIGINT; t_saida_alm BIGINT; t_retorno_alm BIGINT; t_fim_exp BIGINT; t_work_session BIGINT;
    t_intencao BIGINT; t_checkin BIGINT; t_cerveja BIGINT; t_debito BIGINT; t_credito BIGINT;
BEGIN
    SELECT tag_id INTO t_agua FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Agua';
    SELECT tag_id INTO t_cafe_manha FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Cafe da Manha';
    SELECT tag_id INTO t_almoco FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Almoco';
    SELECT tag_id INTO t_lanche FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Lanche';
    SELECT tag_id INTO t_janta FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Janta';
    SELECT tag_id INTO t_cafe_preto FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Cafe Preto';
    SELECT tag_id INTO t_caminhar FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Caminhar';
    SELECT tag_id INTO t_flexao FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Flexao';
    SELECT tag_id INTO t_alongamento FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Alongamento';
    SELECT tag_id INTO t_abs FROM aion_api.tags WHERE user_id = v_user_id AND name = 'ABS';
    SELECT tag_id INTO t_corrida FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Corrida';
    SELECT tag_id INTO t_barra FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Barra';
    SELECT tag_id INTO t_sol FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Sol';
    SELECT tag_id INTO t_academia FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Academia';
    SELECT tag_id INTO t_golang FROM aion_api.tags WHERE user_id = v_user_id AND name = 'golang';
    SELECT tag_id INTO t_aws FROM aion_api.tags WHERE user_id = v_user_id AND name = 'AWS';
    SELECT tag_id INTO t_python FROM aion_api.tags WHERE user_id = v_user_id AND name = 'python';
    SELECT tag_id INTO t_react FROM aion_api.tags WHERE user_id = v_user_id AND name = 'react';
    SELECT tag_id INTO t_brq FROM aion_api.tags WHERE user_id = v_user_id AND name = 'BRQ';
    SELECT tag_id INTO t_aion FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Aion';
    SELECT tag_id INTO t_kafka FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Kafka-stream';
    SELECT tag_id INTO t_meditacao FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Meditacao';
    SELECT tag_id INTO t_auto_hipnose FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Auto-Hipnose';
    SELECT tag_id INTO t_relax FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Relax';
    SELECT tag_id INTO t_diario_p FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Diario Pessoal';
    SELECT tag_id INTO t_diario_w FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Diario Profissional';
    SELECT tag_id INTO t_dormir FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Dormir';
    SELECT tag_id INTO t_acordar FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Acordar';
    SELECT tag_id INTO t_energia FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Energia %';
    SELECT tag_id INTO t_muito_triste FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Muito Triste';
    SELECT tag_id INTO t_triste FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Triste';
    SELECT tag_id INTO t_neutro FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Neutro';
    SELECT tag_id INTO t_feliz FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Feliz';
    SELECT tag_id INTO t_muito_feliz FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Muito Feliz';
    SELECT tag_id INTO t_inicio_exp FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Inicio Expediente';
    SELECT tag_id INTO t_saida_alm FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Saida Almoco';
    SELECT tag_id INTO t_retorno_alm FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Retorno Almoco';
    SELECT tag_id INTO t_fim_exp FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Fim Expediente';
    SELECT tag_id INTO t_work_session FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Sessao de Trabalho';
    SELECT tag_id INTO t_intencao FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Intencao';
    SELECT tag_id INTO t_checkin FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Check-in';
    SELECT tag_id INTO t_cerveja FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Cerveja';
    SELECT tag_id INTO t_debito FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Debito';
    SELECT tag_id INTO t_credito FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Credito';

    FOR v_day IN
        SELECT generate_series(CURRENT_DATE - INTERVAL '90 days', CURRENT_DATE - INTERVAL '2 days', INTERVAL '1 day')::date
    LOOP
        v_day_idx := v_day_idx + 1;
        v_is_weekend := EXTRACT(ISODOW FROM v_day) IN (6, 7);
        v_target_count := CASE WHEN v_is_weekend THEN 50 + (v_day_idx % 4) ELSE 56 + (v_day_idx % 5) END;
        v_inserted := 0;

        -- Sleep window (Dormir + Acordar)
        v_duration := 23000 + ((v_day_idx * 97) % 11000);
        v_ts := ((v_day + TIME '00:10:00')::timestamp AT TIME ZONE v_tz);
        INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
        VALUES (v_user_id, t_dormir, 'Inicio do sono principal', v_duration, v_ts, v_ts + INTERVAL '2 minutes', 'test_seed_realistic', 'published');
        v_inserted := v_inserted + 1;

        INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
        VALUES (v_user_id, t_acordar, 'Fim do sono principal', 60, v_ts + (v_duration || ' seconds')::interval, v_ts + (v_duration || ' seconds')::interval + INTERVAL '2 minutes', 'test_seed_realistic', 'published');
        v_inserted := v_inserted + 1;

        -- Work markers
        v_ts := ((v_day + TIME '08:30:00')::timestamp AT TIME ZONE v_tz);
        INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
        VALUES (v_user_id, t_inicio_exp, 'Inicio de expediente', 60, v_ts, v_ts + INTERVAL '1 minute', 'test_seed_realistic', 'published');
        v_inserted := v_inserted + 1;

        v_ts := ((v_day + TIME '12:05:00')::timestamp AT TIME ZONE v_tz);
        INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
        VALUES (v_user_id, t_saida_alm, 'Saida para almoco', 60, v_ts, v_ts + INTERVAL '1 minute', 'test_seed_realistic', 'published');
        v_inserted := v_inserted + 1;

        v_ts := ((v_day + TIME '13:05:00')::timestamp AT TIME ZONE v_tz);
        INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
        VALUES (v_user_id, t_retorno_alm, 'Retorno do almoco', 60, v_ts, v_ts + INTERVAL '1 minute', 'test_seed_realistic', 'published');
        v_inserted := v_inserted + 1;

        v_ts := ((v_day + TIME '18:00:00')::timestamp AT TIME ZONE v_tz);
        INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
        VALUES (v_user_id, t_fim_exp, 'Fim do expediente', 60, v_ts, v_ts + INTERVAL '1 minute', 'test_seed_realistic', 'published');
        v_inserted := v_inserted + 1;

        -- Work sessions (metric backbone)
        FOR v_i IN 1..CASE WHEN v_is_weekend THEN 4 + (v_day_idx % 2) ELSE 10 + (v_day_idx % 3) END LOOP
            v_ts := ((v_day + TIME '08:40:00')::timestamp AT TIME ZONE v_tz) + ((v_i - 1) * INTERVAL '48 minutes');
            v_duration := 1500 + ((v_day_idx * 19 + v_i * 23) % 1800);
            INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
            VALUES (v_user_id, t_work_session, 'Sessao de foco/trabalho', v_duration, v_ts, v_ts + INTERVAL '3 minutes', 'test_seed_realistic', 'published');
            v_inserted := v_inserted + 1;
        END LOOP;

        -- Hydration
        FOR v_i IN 1..CASE WHEN v_is_weekend THEN 10 + (v_day_idx % 3) ELSE 12 + (v_day_idx % 3) END LOOP
            v_ts := ((v_day + TIME '06:40:00')::timestamp AT TIME ZONE v_tz) + ((v_i - 1) * INTERVAL '75 minutes');
            v_value := 180 + ((v_day_idx * 11 + v_i * 13) % 170);
            INSERT INTO aion_api.records (user_id, tag_id, description, value, duration_seconds, event_time, recorded_at, source, status)
            VALUES (v_user_id, t_agua, 'Copo de agua', v_value, 120, v_ts, v_ts + INTERVAL '2 minutes', 'test_seed_realistic', 'published');
            v_inserted := v_inserted + 1;
        END LOOP;

        -- Meals
        v_ts := ((v_day + TIME '07:35:00')::timestamp AT TIME ZONE v_tz);
        INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
        VALUES (v_user_id, t_cafe_manha, 'Cafe da manha', 1200, v_ts, v_ts + INTERVAL '5 minutes', 'test_seed_realistic', 'published');
        v_inserted := v_inserted + 1;

        v_ts := ((v_day + TIME '12:15:00')::timestamp AT TIME ZONE v_tz);
        INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
        VALUES (v_user_id, t_almoco, 'Almoco', 1800, v_ts, v_ts + INTERVAL '4 minutes', 'test_seed_realistic', 'published');
        v_inserted := v_inserted + 1;

        v_ts := ((v_day + TIME '16:30:00')::timestamp AT TIME ZONE v_tz);
        INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
        VALUES (v_user_id, t_lanche, 'Lanche da tarde', 900, v_ts, v_ts + INTERVAL '3 minutes', 'test_seed_realistic', 'published');
        v_inserted := v_inserted + 1;

        v_ts := ((v_day + TIME '20:00:00')::timestamp AT TIME ZONE v_tz);
        INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
        VALUES (v_user_id, t_janta, 'Janta', 1500, v_ts, v_ts + INTERVAL '4 minutes', 'test_seed_realistic', 'published');
        v_inserted := v_inserted + 1;

        -- Coffee breaks
        FOR v_i IN 1..3 LOOP
            v_ts := ((v_day + TIME '09:45:00')::timestamp AT TIME ZONE v_tz) + ((v_i - 1) * INTERVAL '2 hours 35 minutes');
            INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
            VALUES (v_user_id, t_cafe_preto, 'Cafe preto', 480, v_ts, v_ts + INTERVAL '2 minutes', 'test_seed_realistic', 'published');
            v_inserted := v_inserted + 1;
        END LOOP;

        -- Energy checkpoints (0-100)
        FOR v_i IN 1..4 LOOP
            v_ts := ((v_day + TIME '08:10:00')::timestamp AT TIME ZONE v_tz) + ((v_i - 1) * INTERVAL '3 hours 30 minutes');
            v_value := 45 + ((v_day_idx * 7 + v_i * 11) % 56);
            INSERT INTO aion_api.records (user_id, tag_id, description, value, duration_seconds, event_time, recorded_at, source, status)
            VALUES (v_user_id, t_energia, 'Energia em porcentagem', v_value, 60, v_ts, v_ts + INTERVAL '1 minute', 'test_seed_realistic', 'published');
            v_inserted := v_inserted + 1;
        END LOOP;

        -- Mood checkpoints (tag scale + numeric value)
        FOR v_i IN 1..4 LOOP
            v_ts := ((v_day + TIME '08:20:00')::timestamp AT TIME ZONE v_tz) + ((v_i - 1) * INTERVAL '3 hours 20 minutes');
            v_value := ((v_day_idx * 5 + v_i * 9) % 101);
            v_mood_name := CASE
                WHEN v_value < 20 THEN 'Muito Triste'
                WHEN v_value < 40 THEN 'Triste'
                WHEN v_value < 60 THEN 'Neutro'
                WHEN v_value < 80 THEN 'Feliz'
                ELSE 'Muito Feliz'
            END;

            INSERT INTO aion_api.records (user_id, tag_id, description, value, duration_seconds, event_time, recorded_at, source, status)
            VALUES (
                v_user_id,
                CASE v_mood_name
                    WHEN 'Muito Triste' THEN t_muito_triste
                    WHEN 'Triste' THEN t_triste
                    WHEN 'Neutro' THEN t_neutro
                    WHEN 'Feliz' THEN t_feliz
                    ELSE t_muito_feliz
                END,
                'Humor em escala emocional',
                v_value,
                60,
                v_ts,
                v_ts + INTERVAL '1 minute',
                'test_seed_realistic',
                'published'
            );
            v_inserted := v_inserted + 1;
        END LOOP;

        -- Physical and mental health
        v_ts := ((v_day + TIME '18:30:00')::timestamp AT TIME ZONE v_tz);
        INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
        VALUES (v_user_id, t_caminhar, 'Caminhada do dia', 1800 + ((v_day_idx * 5) % 1800), v_ts, v_ts + INTERVAL '4 minutes', 'test_seed_realistic', 'published');
        v_inserted := v_inserted + 1;

        IF (v_day_idx % 2) = 0 THEN
            v_ts := ((v_day + TIME '07:00:00')::timestamp AT TIME ZONE v_tz);
            INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
            VALUES (v_user_id, t_flexao, 'Treino de flexoes', 600 + ((v_day_idx * 17) % 600), v_ts, v_ts + INTERVAL '3 minutes', 'test_seed_realistic', 'published');
            v_inserted := v_inserted + 1;
        END IF;

        IF (v_day_idx % 3) = 0 THEN
            v_ts := ((v_day + TIME '06:35:00')::timestamp AT TIME ZONE v_tz);
            INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
            VALUES (v_user_id, t_alongamento, 'Alongamento matinal', 480 + ((v_day_idx * 7) % 360), v_ts, v_ts + INTERVAL '2 minutes', 'test_seed_realistic', 'published');
            v_inserted := v_inserted + 1;
        END IF;

        IF (v_day_idx % 4) = 0 THEN
            v_ts := ((v_day + TIME '18:05:00')::timestamp AT TIME ZONE v_tz);
            INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
            VALUES (v_user_id, t_corrida, 'Corrida leve', 1500 + ((v_day_idx * 11) % 1200), v_ts, v_ts + INTERVAL '3 minutes', 'test_seed_realistic', 'published');
            v_inserted := v_inserted + 1;
        END IF;

        IF (v_day_idx % 5) = 0 THEN
            v_ts := ((v_day + TIME '07:20:00')::timestamp AT TIME ZONE v_tz);
            INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
            VALUES (
                v_user_id,
                CASE WHEN (v_day_idx % 10) = 0 THEN t_barra ELSE t_abs END,
                CASE WHEN (v_day_idx % 10) = 0 THEN 'Treino de barra fixa' ELSE 'Treino de abdominal' END,
                600 + ((v_day_idx * 13) % 900),
                v_ts,
                v_ts + INTERVAL '3 minutes',
                'test_seed_realistic',
                'published'
            );
            v_inserted := v_inserted + 1;
        END IF;

        IF (v_day_idx % 6) = 0 THEN
            v_ts := ((v_day + TIME '11:40:00')::timestamp AT TIME ZONE v_tz);
            INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
            VALUES (v_user_id, t_sol, 'Exposicao ao sol', 900, v_ts, v_ts + INTERVAL '2 minutes', 'test_seed_realistic', 'published');
            v_inserted := v_inserted + 1;
        END IF;

        IF (v_day_idx % 7) = 0 THEN
            v_ts := ((v_day + TIME '19:00:00')::timestamp AT TIME ZONE v_tz);
            INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
            VALUES (v_user_id, t_academia, 'Treino de academia', 2100 + ((v_day_idx * 3) % 1800), v_ts, v_ts + INTERVAL '3 minutes', 'test_seed_realistic', 'published');
            v_inserted := v_inserted + 1;
        END IF;

        IF (v_day_idx % 3) = 0 THEN
            v_ts := ((v_day + TIME '22:00:00')::timestamp AT TIME ZONE v_tz);
            INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
            VALUES (v_user_id, t_meditacao, 'Sessao de meditacao', 900, v_ts, v_ts + INTERVAL '2 minutes', 'test_seed_realistic', 'published');
            v_inserted := v_inserted + 1;
        END IF;

        IF (v_day_idx % 5) = 0 THEN
            v_ts := ((v_day + TIME '22:30:00')::timestamp AT TIME ZONE v_tz);
            INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
            VALUES (v_user_id, t_auto_hipnose, 'Sessao de auto-hipnose', 900, v_ts, v_ts + INTERVAL '2 minutes', 'test_seed_realistic', 'published');
            v_inserted := v_inserted + 1;
        END IF;

        -- Studies/projects/work context
        v_ts := ((v_day + TIME '20:15:00')::timestamp AT TIME ZONE v_tz);
        INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
        VALUES (
            v_user_id,
            CASE (v_day_idx % 4)
                WHEN 0 THEN t_golang
                WHEN 1 THEN t_aws
                WHEN 2 THEN t_python
                ELSE t_react
            END,
            'Bloco de estudo tecnico',
            3600 + ((v_day_idx * 13) % 2700),
            v_ts,
            v_ts + INTERVAL '5 minutes',
            'test_seed_realistic',
            'published'
        );
        v_inserted := v_inserted + 1;

        IF (v_day_idx % 2) = 1 THEN
            v_ts := ((v_day + TIME '19:10:00')::timestamp AT TIME ZONE v_tz);
            INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
            VALUES (v_user_id, t_aion, 'Projeto Aion', 3000 + ((v_day_idx * 11) % 2400), v_ts, v_ts + INTERVAL '3 minutes', 'test_seed_realistic', 'published');
            v_inserted := v_inserted + 1;
        ELSE
            v_ts := ((v_day + TIME '19:10:00')::timestamp AT TIME ZONE v_tz);
            INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
            VALUES (v_user_id, t_kafka, 'Projeto Kafka-stream', 3000 + ((v_day_idx * 11) % 2400), v_ts, v_ts + INTERVAL '3 minutes', 'test_seed_realistic', 'published');
            v_inserted := v_inserted + 1;
        END IF;

        IF NOT v_is_weekend THEN
            v_ts := ((v_day + TIME '14:40:00')::timestamp AT TIME ZONE v_tz);
            INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
            VALUES (v_user_id, t_brq, 'Atividade BRQ', 2400 + ((v_day_idx * 17) % 2400), v_ts, v_ts + INTERVAL '4 minutes', 'test_seed_realistic', 'published');
            v_inserted := v_inserted + 1;
        END IF;

        -- Intentions and check-ins
        FOR v_i IN 1..4 LOOP
            v_ts := ((v_day + TIME '07:05:00')::timestamp AT TIME ZONE v_tz) + ((v_i - 1) * INTERVAL '4 hours 20 minutes');
            INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
            VALUES (v_user_id, t_intencao, 'Intencao do dia #' || v_i, 90, v_ts, v_ts + INTERVAL '1 minute', 'test_seed_realistic', CASE WHEN (v_day_idx + v_i) % 5 = 0 THEN 'pending' ELSE 'published' END);
            v_inserted := v_inserted + 1;
        END LOOP;

        -- Diaries
        v_ts := ((v_day + TIME '21:35:00')::timestamp AT TIME ZONE v_tz);
        INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
        VALUES (v_user_id, t_diario_w, 'Resumo profissional do dia', 780, v_ts, v_ts + INTERVAL '4 minutes', 'test_seed_realistic', 'published');
        v_inserted := v_inserted + 1;

        v_ts := ((v_day + TIME '22:15:00')::timestamp AT TIME ZONE v_tz);
        INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
        VALUES (v_user_id, t_diario_p, 'Resumo pessoal do dia', 900, v_ts, v_ts + INTERVAL '4 minutes', 'test_seed_realistic', 'published');
        v_inserted := v_inserted + 1;

        -- Occasional legacy finance/leisure beverage
        IF (v_day_idx % 6) = 0 THEN
            v_ts := ((v_day + TIME '13:20:00')::timestamp AT TIME ZONE v_tz);
            INSERT INTO aion_api.records (user_id, tag_id, description, value, duration_seconds, event_time, recorded_at, source, status)
            VALUES (v_user_id, t_debito, 'Compra no debito', 38 + ((v_day_idx * 3) % 120), 120, v_ts, v_ts + INTERVAL '2 minutes', 'test_seed_realistic', 'published');
            v_inserted := v_inserted + 1;
        END IF;

        IF (v_day_idx % 8) = 0 THEN
            v_ts := ((v_day + TIME '13:50:00')::timestamp AT TIME ZONE v_tz);
            INSERT INTO aion_api.records (user_id, tag_id, description, value, duration_seconds, event_time, recorded_at, source, status)
            VALUES (v_user_id, t_credito, 'Compra no credito', 65 + ((v_day_idx * 5) % 200), 120, v_ts, v_ts + INTERVAL '2 minutes', 'test_seed_realistic', 'published');
            v_inserted := v_inserted + 1;
        END IF;

        IF v_is_weekend AND (v_day_idx % 3) = 0 THEN
            v_ts := ((v_day + TIME '20:40:00')::timestamp AT TIME ZONE v_tz);
            INSERT INTO aion_api.records (user_id, tag_id, description, value, duration_seconds, event_time, recorded_at, source, status)
            VALUES (v_user_id, t_cerveja, 'Consumo de cerveja', 355, 1200, v_ts, v_ts + INTERVAL '3 minutes', 'test_seed_realistic', 'published');
            v_inserted := v_inserted + 1;
        END IF;

        v_ts := ((v_day + TIME '23:10:00')::timestamp AT TIME ZONE v_tz);
        INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
        VALUES (v_user_id, t_relax, 'Relaxamento noturno', 1500, v_ts, v_ts + INTERVAL '2 minutes', 'test_seed_realistic', 'published');
        v_inserted := v_inserted + 1;

        -- Fill until target_count with quick check-ins (always with duration)
        v_extra_idx := 0;
        WHILE v_inserted < v_target_count LOOP
            v_extra_idx := v_extra_idx + 1;
            v_ts := ((v_day + TIME '06:00:00')::timestamp AT TIME ZONE v_tz) + (v_extra_idx * INTERVAL '17 minutes');

            IF (v_extra_idx + v_day_idx) % 3 = 0 THEN
                INSERT INTO aion_api.records (user_id, tag_id, description, value, duration_seconds, event_time, recorded_at, source, status)
                VALUES (v_user_id, t_agua, 'Check-in rapido de hidratacao', 100 + ((v_extra_idx * 9) % 140), 60, v_ts, v_ts + INTERVAL '1 minute', 'test_seed_realistic', 'published');
            ELSIF (v_extra_idx + v_day_idx) % 3 = 1 THEN
                INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
                VALUES (v_user_id, t_work_session, 'Micro-sessao de foco', 480 + ((v_extra_idx * 13) % 900), v_ts, v_ts + INTERVAL '1 minute', 'test_seed_realistic', 'published');
            ELSE
                INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
                VALUES (v_user_id, t_checkin, 'Check-in rapido de progresso', 60, v_ts, v_ts + INTERVAL '1 minute', 'test_seed_realistic', 'published');
            END IF;
            v_inserted := v_inserted + 1;
        END LOOP;
    END LOOP;
END $$;

-- --------------------------------------------------------------------------
-- 6) Today records (real-time visibility, no future timestamps)
-- --------------------------------------------------------------------------
DO $$
DECLARE
    v_user_id BIGINT := 999;
    v_now TIMESTAMPTZ := NOW();
    v_tz TEXT := 'America/Sao_Paulo';
    v_base TIMESTAMPTZ := NOW() - INTERVAL '12 hours';
    v_i INT;
    v_ts TIMESTAMPTZ;
    t_dormir BIGINT; t_acordar BIGINT; t_work_session BIGINT; t_agua BIGINT; t_cafe BIGINT; t_almoco BIGINT;
    t_lanche BIGINT; t_janta BIGINT; t_energia BIGINT; t_neutro BIGINT; t_feliz BIGINT; t_intencao BIGINT;
    t_inicio_exp BIGINT; t_saida_alm BIGINT; t_retorno_alm BIGINT; t_fim_exp BIGINT; t_checkin BIGINT;
    t_diario_w BIGINT; t_diario_p BIGINT; t_caminhar BIGINT;
BEGIN
    SELECT tag_id INTO t_dormir FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Dormir';
    SELECT tag_id INTO t_acordar FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Acordar';
    SELECT tag_id INTO t_work_session FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Sessao de Trabalho';
    SELECT tag_id INTO t_agua FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Agua';
    SELECT tag_id INTO t_cafe FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Cafe da Manha';
    SELECT tag_id INTO t_almoco FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Almoco';
    SELECT tag_id INTO t_lanche FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Lanche';
    SELECT tag_id INTO t_janta FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Janta';
    SELECT tag_id INTO t_energia FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Energia %';
    SELECT tag_id INTO t_neutro FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Neutro';
    SELECT tag_id INTO t_feliz FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Feliz';
    SELECT tag_id INTO t_intencao FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Intencao';
    SELECT tag_id INTO t_inicio_exp FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Inicio Expediente';
    SELECT tag_id INTO t_saida_alm FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Saida Almoco';
    SELECT tag_id INTO t_retorno_alm FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Retorno Almoco';
    SELECT tag_id INTO t_fim_exp FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Fim Expediente';
    SELECT tag_id INTO t_checkin FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Check-in';
    SELECT tag_id INTO t_diario_w FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Diario Profissional';
    SELECT tag_id INTO t_diario_p FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Diario Pessoal';
    SELECT tag_id INTO t_caminhar FROM aion_api.tags WHERE user_id = v_user_id AND name = 'Caminhar';

    -- Sleep (ended today)
    v_ts := v_now - INTERVAL '11 hours 30 minutes';
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
    VALUES (v_user_id, t_dormir, 'Sono da madrugada (hoje)', 27000, v_ts, v_ts + INTERVAL '2 minutes', 'test_seed_realistic', 'published');

    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
    VALUES (v_user_id, t_acordar, 'Acordar (hoje)', 60, v_ts + INTERVAL '27000 seconds', v_ts + INTERVAL '27002 seconds', 'test_seed_realistic', 'published');

    -- Work markers today
    FOR v_i IN 1..4 LOOP
        v_ts := CASE v_i
            WHEN 1 THEN v_now - INTERVAL '8 hours 30 minutes'
            WHEN 2 THEN v_now - INTERVAL '5 hours'
            WHEN 3 THEN v_now - INTERVAL '4 hours'
            ELSE v_now - INTERVAL '1 hour 30 minutes'
        END;
        IF v_ts > v_now - INTERVAL '2 minutes' THEN
            v_ts := v_now - INTERVAL '2 minutes';
        END IF;
        INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
        VALUES (
            v_user_id,
            CASE v_i WHEN 1 THEN t_inicio_exp WHEN 2 THEN t_saida_alm WHEN 3 THEN t_retorno_alm ELSE t_fim_exp END,
            'Marcador de expediente (hoje)',
            60,
            v_ts,
            v_ts + INTERVAL '1 minute',
            'test_seed_realistic',
            'published'
        );
    END LOOP;

    -- Work sessions today
    FOR v_i IN 1..8 LOOP
        v_ts := v_now - INTERVAL '9 hours' + ((v_i - 1) * INTERVAL '45 minutes');
        IF v_ts > v_now - INTERVAL '2 minutes' THEN
            v_ts := v_now - INTERVAL '2 minutes';
        END IF;
        INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
        VALUES (v_user_id, t_work_session, 'Sessao de trabalho atual #' || v_i, 2100 + ((v_i * 73) % 900), v_ts, v_ts + INTERVAL '2 minutes', 'test_seed_realistic', 'published');
    END LOOP;

    -- Water today
    FOR v_i IN 1..12 LOOP
        v_ts := v_base + ((v_i - 1) * INTERVAL '50 minutes');
        IF v_ts > v_now - INTERVAL '2 minutes' THEN
            v_ts := v_now - INTERVAL '2 minutes';
        END IF;
        INSERT INTO aion_api.records (user_id, tag_id, description, value, duration_seconds, event_time, recorded_at, source, status)
        VALUES (v_user_id, t_agua, 'Agua hoje', 180 + ((v_i * 17) % 160), 120, v_ts, v_ts + INTERVAL '1 minute', 'test_seed_realistic', 'published');
    END LOOP;

    -- Meals today
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
    VALUES (v_user_id, t_cafe, 'Cafe da manha de hoje', 1200, v_now - INTERVAL '10 hours 20 minutes', v_now - INTERVAL '10 hours 15 minutes', 'test_seed_realistic', 'published');
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
    VALUES (v_user_id, t_almoco, 'Almoco de hoje', 1800, v_now - INTERVAL '6 hours 10 minutes', v_now - INTERVAL '6 hours', 'test_seed_realistic', 'published');
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
    VALUES (v_user_id, t_lanche, 'Lanche de hoje', 900, v_now - INTERVAL '3 hours 40 minutes', v_now - INTERVAL '3 hours 35 minutes', 'test_seed_realistic', 'published');
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
    VALUES (v_user_id, t_janta, 'Janta de hoje', 1500, v_now - INTERVAL '1 hour 20 minutes', v_now - INTERVAL '1 hour 15 minutes', 'test_seed_realistic', 'published');

    -- Energy + mood today
    FOR v_i IN 1..4 LOOP
        v_ts := v_now - INTERVAL '9 hours' + ((v_i - 1) * INTERVAL '2 hours 20 minutes');
        IF v_ts > v_now - INTERVAL '2 minutes' THEN
            v_ts := v_now - INTERVAL '2 minutes';
        END IF;
        INSERT INTO aion_api.records (user_id, tag_id, description, value, duration_seconds, event_time, recorded_at, source, status)
        VALUES (v_user_id, t_energia, 'Energia hoje', 52 + ((v_i * 9) % 41), 60, v_ts, v_ts + INTERVAL '1 minute', 'test_seed_realistic', 'published');
        INSERT INTO aion_api.records (user_id, tag_id, description, value, duration_seconds, event_time, recorded_at, source, status)
        VALUES (v_user_id, CASE WHEN (v_i % 2) = 0 THEN t_feliz ELSE t_neutro END, 'Humor hoje', 45 + ((v_i * 13) % 51), 60, v_ts + INTERVAL '8 minutes', v_ts + INTERVAL '9 minutes', 'test_seed_realistic', 'published');
    END LOOP;

    -- Intentions today
    FOR v_i IN 1..4 LOOP
        v_ts := v_now - INTERVAL '11 hours' + ((v_i - 1) * INTERVAL '2 hours 35 minutes');
        INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
        VALUES (v_user_id, t_intencao, 'Intencao de hoje #' || v_i, 90, v_ts, v_ts + INTERVAL '1 minute', 'test_seed_realistic', CASE WHEN v_i = 4 THEN 'pending' ELSE 'published' END);
    END LOOP;

    -- Health + diaries today
    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
    VALUES (v_user_id, t_caminhar, 'Treino de hoje', 2400, v_now - INTERVAL '2 hours 55 minutes', v_now - INTERVAL '2 hours 50 minutes', 'test_seed_realistic', 'published');

    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
    VALUES (v_user_id, t_diario_w, 'Resumo profissional de hoje', 840, v_now - INTERVAL '1 hour 40 minutes', v_now - INTERVAL '1 hour 35 minutes', 'test_seed_realistic', 'published');

    INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
    VALUES (v_user_id, t_diario_p, 'Resumo pessoal de hoje', 960, v_now - INTERVAL '1 hour 10 minutes', v_now - INTERVAL '1 hour 5 minutes', 'test_seed_realistic', 'published');

    -- Dense day with additional check-ins
    FOR v_i IN 1..12 LOOP
        v_ts := v_now - INTERVAL '6 hours' + ((v_i - 1) * INTERVAL '23 minutes');
        IF v_ts > v_now - INTERVAL '2 minutes' THEN
            v_ts := v_now - INTERVAL '2 minutes';
        END IF;
        INSERT INTO aion_api.records (user_id, tag_id, description, duration_seconds, event_time, recorded_at, source, status)
        VALUES (v_user_id, t_checkin, 'Check-in rapido de hoje', 60, v_ts, v_ts + INTERVAL '1 minute', 'test_seed_realistic', 'published');
    END LOOP;
END $$;

COMMIT;
