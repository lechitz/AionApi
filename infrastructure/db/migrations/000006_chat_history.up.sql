-- Migration: 000006_chat_history
-- Description: Create chat_history table for AI conversation storage
-- This consolidates: 12_chat_history.sql

CREATE TABLE IF NOT EXISTS aion_api.chat_history (
    chat_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    message TEXT NOT NULL,
    response TEXT NOT NULL,
    tokens_used INTEGER DEFAULT 0,
    function_calls JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Foreign key constraint
    CONSTRAINT fk_chat_history_user
        FOREIGN KEY (user_id)
        REFERENCES aion_api.users(user_id)
        ON DELETE CASCADE
);

-- Indexes for query performance (using partial indexes for active records only)
CREATE INDEX IF NOT EXISTS idx_chat_history_user_id
    ON aion_api.chat_history(user_id)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_chat_history_created_at
    ON aion_api.chat_history(created_at DESC)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_chat_history_user_created
    ON aion_api.chat_history(user_id, created_at DESC)
    WHERE deleted_at IS NULL;

-- Table and column comments
COMMENT ON TABLE aion_api.chat_history IS 'Stores conversation history between users and AI assistant';
COMMENT ON COLUMN aion_api.chat_history.chat_id IS 'Primary key - unique chat message identifier';
COMMENT ON COLUMN aion_api.chat_history.user_id IS 'Foreign key to users table';
COMMENT ON COLUMN aion_api.chat_history.message IS 'User message sent to AI';
COMMENT ON COLUMN aion_api.chat_history.response IS 'AI response to user message';
COMMENT ON COLUMN aion_api.chat_history.tokens_used IS 'Number of tokens consumed by this chat interaction';
COMMENT ON COLUMN aion_api.chat_history.function_calls IS 'JSON array of function calls made by AI (if any)';
COMMENT ON COLUMN aion_api.chat_history.created_at IS 'Timestamp when the chat message was created';
COMMENT ON COLUMN aion_api.chat_history.updated_at IS 'Timestamp when the chat message was last updated';
COMMENT ON COLUMN aion_api.chat_history.deleted_at IS 'Soft delete timestamp (NULL if active)';
