-- Migration: 000009_phase1_ai_optimizations
-- Description: Phase 1 AI optimizations - pg_trgm, usage stats, chat metadata
-- Impact: 10-100x performance improvement on LIKE queries and prepares for future usage-based suggestions

-- ============================================================================
-- 1. Enable pg_trgm extension for trigram similarity search
-- ============================================================================

-- pg_trgm enables fast LIKE/ILIKE queries and similarity search
CREATE EXTENSION IF NOT EXISTS pg_trgm WITH SCHEMA public;

-- Add trigram index on records.description for fast LIKE/ILIKE queries
-- This improves queries like: WHERE description ILIKE '%café%'
-- Expected improvement: 10-100x faster than sequential scan
CREATE INDEX IF NOT EXISTS idx_records_description_trgm 
  ON aion_api.records USING gin (description public.gin_trgm_ops)
  WHERE deleted_at IS NULL;

COMMENT ON INDEX aion_api.idx_records_description_trgm IS 
  'Trigram index for fast LIKE/ILIKE queries on record descriptions (AI query optimization)';

-- ============================================================================
-- 2. Add usage statistics to categories
-- ============================================================================

-- usage_count: tracks how many records reference this category (via tags)
-- last_used_at: timestamp of last record creation in this category
-- These enable LLM to prioritize frequently used categories in suggestions
ALTER TABLE aion_api.categories
  ADD COLUMN IF NOT EXISTS usage_count INTEGER DEFAULT 0,
  ADD COLUMN IF NOT EXISTS last_used_at TIMESTAMPTZ;

COMMENT ON COLUMN aion_api.categories.usage_count IS 
  'Number of records associated with this category (updated via triggers in future phases)';
COMMENT ON COLUMN aion_api.categories.last_used_at IS 
  'Timestamp of last record creation in this category (enables LLM prioritization)';

-- ============================================================================
-- 3. Add usage statistics to tags
-- ============================================================================

-- Similar to categories, enables LLM to suggest most-used tags first
ALTER TABLE aion_api.tags
  ADD COLUMN IF NOT EXISTS usage_count INTEGER DEFAULT 0,
  ADD COLUMN IF NOT EXISTS last_used_at TIMESTAMPTZ;

COMMENT ON COLUMN aion_api.tags.usage_count IS 
  'Number of records associated with this tag (updated via triggers in future phases)';
COMMENT ON COLUMN aion_api.tags.last_used_at IS 
  'Timestamp of last record creation with this tag (enables LLM prioritization)';

-- ============================================================================
-- 4. Add metadata to chat_history for performance tracking
-- ============================================================================

-- session_id: groups related conversations together
-- execution_time_ms: tracks query performance for monitoring
-- tool_count: number of GraphQL tools called during this interaction
-- error_count: number of errors during execution (debugging)
ALTER TABLE aion_api.chat_history
  ADD COLUMN IF NOT EXISTS session_id UUID,
  ADD COLUMN IF NOT EXISTS execution_time_ms INTEGER,
  ADD COLUMN IF NOT EXISTS tool_count INTEGER DEFAULT 0,
  ADD COLUMN IF NOT EXISTS error_count INTEGER DEFAULT 0;

COMMENT ON COLUMN aion_api.chat_history.session_id IS 
  'Groups related conversations together (enables context analysis)';
COMMENT ON COLUMN aion_api.chat_history.execution_time_ms IS 
  'Query execution time in milliseconds (performance monitoring)';
COMMENT ON COLUMN aion_api.chat_history.tool_count IS 
  'Number of GraphQL tools called during this interaction';
COMMENT ON COLUMN aion_api.chat_history.error_count IS 
  'Number of errors during execution (debugging aid)';

-- Add index for session-based queries (conversation threading)
CREATE INDEX IF NOT EXISTS idx_chat_history_session_time
  ON aion_api.chat_history(session_id, created_at DESC)
  WHERE deleted_at IS NULL;

COMMENT ON INDEX aion_api.idx_chat_history_session_time IS 
  'Fast retrieval of conversation threads by session';

-- ============================================================================
-- 5. Initialize usage_count for existing data (backfill)
-- ============================================================================

-- Count existing records per tag and update usage_count
UPDATE aion_api.tags t
SET usage_count = (
    SELECT COUNT(*)
    FROM aion_api.records r
    WHERE r.tag_id = t.tag_id 
      AND r.deleted_at IS NULL
),
last_used_at = (
    SELECT MAX(r.event_time)
    FROM aion_api.records r
    WHERE r.tag_id = t.tag_id 
      AND r.deleted_at IS NULL
)
WHERE t.deleted_at IS NULL;

-- Count existing records per category (via tags) and update usage_count
UPDATE aion_api.categories c
SET usage_count = (
    SELECT COUNT(*)
    FROM aion_api.records r
    INNER JOIN aion_api.tags t ON r.tag_id = t.tag_id
    WHERE t.category_id = c.category_id 
      AND r.deleted_at IS NULL
      AND t.deleted_at IS NULL
),
last_used_at = (
    SELECT MAX(r.event_time)
    FROM aion_api.records r
    INNER JOIN aion_api.tags t ON r.tag_id = t.tag_id
    WHERE t.category_id = c.category_id 
      AND r.deleted_at IS NULL
      AND t.deleted_at IS NULL
)
WHERE c.deleted_at IS NULL;

-- ============================================================================
-- Summary
-- ============================================================================

-- Phase 1 optimizations completed:
-- ✅ pg_trgm extension enabled
-- ✅ Trigram index on records.description (10-100x faster LIKE queries)
-- ✅ usage_count + last_used_at on categories (LLM prioritization)
-- ✅ usage_count + last_used_at on tags (LLM prioritization)
-- ✅ session_id, execution_time_ms, tool_count, error_count on chat_history (monitoring)
-- ✅ Backfilled usage stats from existing data

-- Expected impact:
-- - LIKE/ILIKE queries: 10-100x faster
-- - Zero overhead (columns empty until Phase 2 triggers)
-- - Prepares for future LLM enhancements

-- Next phases (future):
-- Phase 2: Triggers to auto-update usage_count (when volume > 10k records)
-- Phase 3: Vector embeddings for semantic search (when chat history > 10k messages)
