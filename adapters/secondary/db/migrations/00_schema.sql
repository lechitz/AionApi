-- ===========================================================
-- Arquivo: 00_schema.sql
-- Propósito: Criação do schema principal e funções utilitárias
-- ===========================================================

-- Cria o schema principal do AionAPI, se ainda não existir
CREATE SCHEMA IF NOT EXISTS aion_api;

-- Define o search_path para garantir que tudo execute no schema correto
SET search_path TO aion_api;

-- ===========================================================
-- Função: update_timestamp
-- Descrição: Atualiza o campo `updated_at` em triggers BEFORE UPDATE
-- ===========================================================

CREATE OR REPLACE FUNCTION aion_api.update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ===========================================================
-- Outras funções globais (caso necessário) devem vir aqui
-- ===========================================================
