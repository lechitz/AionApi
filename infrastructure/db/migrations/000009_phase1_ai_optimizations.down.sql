-- Migration: 000009_phase1_ai_optimizations (DOWN)
-- Description: Rollback Phase 1 AI optimizations

-- ============================================================================
-- 5. Remove usage_count initialization (no-op, data will remain)
-- ============================================================================

-- Nothing to rollback for UPDATE statements

-- ============================================================================
-- 4. Remove chat_history metadata columns
-- ============================================================================

DROP INDEX IF EXISTS aion_api.idx_chat_history_session_time;

ALTER TABLE aion_api.chat_history
  DROP COLUMN IF EXISTS error_count,
  DROP COLUMN IF EXISTS tool_count,
  DROP COLUMN IF EXISTS execution_time_ms,
  DROP COLUMN IF EXISTS session_id;

-- ============================================================================
-- 3. Remove usage statistics from tags
-- ============================================================================

ALTER TABLE aion_api.tags
  DROP COLUMN IF EXISTS last_used_at,
  DROP COLUMN IF EXISTS usage_count;

-- ============================================================================
-- 2. Remove usage statistics from categories
-- ============================================================================

ALTER TABLE aion_api.categories
  DROP COLUMN IF EXISTS last_used_at,
  DROP COLUMN IF EXISTS usage_count;

-- ============================================================================
-- 1. Remove pg_trgm extension and index
-- ============================================================================

DROP INDEX IF EXISTS aion_api.idx_records_description_trgm;

-- Note: We don't drop pg_trgm extension because other databases might use it
-- and dropping extensions can be risky. If you really want to drop it:
-- DROP EXTENSION IF EXISTS pg_trgm CASCADE;

-- ============================================================================
-- Summary
-- ============================================================================

-- Phase 1 optimizations rolled back:
-- ✅ Removed trigram index
-- ✅ Removed usage_count + last_used_at from categories
-- ✅ Removed usage_count + last_used_at from tags
-- ✅ Removed chat_history metadata columns
-- ⚠️  pg_trgm extension NOT dropped (safe default)
