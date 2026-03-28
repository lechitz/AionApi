// Package repository provides methods for interacting with the chat_history database.
package repository

import (
	"github.com/lechitz/aion-api/internal/platform/ports/output/db"
	"github.com/lechitz/aion-api/internal/platform/ports/output/logger"
)

// ChatHistoryRepository manages database operations related to chat history entities.
// Depends on db.DB interface (not *gorm.DB) following Hexagonal Architecture.
type ChatHistoryRepository struct {
	db     db.DB
	logger logger.ContextLogger
}

// New creates a new instance of ChatHistoryRepository with a given database connection and logger.
func New(database db.DB, logger logger.ContextLogger) *ChatHistoryRepository {
	return &ChatHistoryRepository{
		db:     database,
		logger: logger,
	}
}
