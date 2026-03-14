-- Migration: 000001_initial_schema
-- Description: Create initial schema and utility functions
-- This consolidates: 00_schema.sql

CREATE SCHEMA IF NOT EXISTS aion_api;

SET search_path TO aion_api;

CREATE OR REPLACE FUNCTION aion_api.update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
