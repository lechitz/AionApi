// Package model contains database models for the Chat context.
package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	// ChatHistoryTable is the fully qualified table name for chat_history.
	ChatHistoryTable = "aion_api.chat_history"
)

// ChatHistoryDB represents the database row for a chat history entry.
type ChatHistoryDB struct {
	ChatID          uint64         `gorm:"primaryKey;column:chat_id;autoIncrement"`
	UserID          uint64         `gorm:"column:user_id;not null;index"`
	Message         string         `gorm:"column:message;type:text;not null"`
	Response        string         `gorm:"column:response;type:text;not null"`
	TokensUsed      int            `gorm:"column:tokens_used;default:0"`
	FunctionCalls   []byte         `gorm:"column:function_calls;type:jsonb"`
	CreatedAt       time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index"`
	SessionID       *uuid.UUID     `gorm:"column:session_id;type:uuid"`
	ExecutionTimeMs *int           `gorm:"column:execution_time_ms"`
	ToolCount       int            `gorm:"column:tool_count;default:0"`
	ErrorCount      int            `gorm:"column:error_count;default:0"`
}

// TableName implements GORM's tabler interface and returns the fully qualified
// database table name for ChatHistoryDB.
func (ChatHistoryDB) TableName() string {
	return ChatHistoryTable
}
