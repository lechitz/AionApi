// Package mapper provides conversion functions between domain and database models for Chat context.
package mapper

import (
	"encoding/json"

	"github.com/lechitz/aion-api/internal/chat/adapter/secondary/db/model"
	"github.com/lechitz/aion-api/internal/chat/core/domain"
)

// ChatHistoryToDB converts a domain ChatHistory to a database ChatHistoryDB model.
func ChatHistoryToDB(ch domain.ChatHistory) model.ChatHistoryDB {
	db := model.ChatHistoryDB{
		ChatID:          ch.ChatID,
		UserID:          ch.UserID,
		Message:         ch.Message,
		Response:        ch.Response,
		TokensUsed:      ch.TokensUsed,
		CreatedAt:       ch.CreatedAt,
		UpdatedAt:       ch.UpdatedAt,
		SessionID:       ch.SessionID,
		ExecutionTimeMs: ch.ExecutionTimeMs,
		ToolCount:       ch.ToolCount,
		ErrorCount:      ch.ErrorCount,
	}

	// Convert FunctionCalls map to JSONB
	if len(ch.FunctionCalls) > 0 {
		if jsonBytes, err := json.Marshal(ch.FunctionCalls); err == nil {
			db.FunctionCalls = jsonBytes
		}
	}

	// Handle soft-delete
	if ch.DeletedAt != nil {
		db.DeletedAt.Time = *ch.DeletedAt
		db.DeletedAt.Valid = true
	}

	return db
}

// ChatHistoryFromDB converts a database ChatHistoryDB model to a domain ChatHistory.
func ChatHistoryFromDB(db model.ChatHistoryDB) domain.ChatHistory {
	ch := domain.ChatHistory{
		ChatID:          db.ChatID,
		UserID:          db.UserID,
		Message:         db.Message,
		Response:        db.Response,
		TokensUsed:      db.TokensUsed,
		CreatedAt:       db.CreatedAt,
		UpdatedAt:       db.UpdatedAt,
		SessionID:       db.SessionID,
		ExecutionTimeMs: db.ExecutionTimeMs,
		ToolCount:       db.ToolCount,
		ErrorCount:      db.ErrorCount,
	}

	// Convert JSONB to FunctionCalls map
	if len(db.FunctionCalls) > 0 {
		var functionCalls map[string]string
		if err := json.Unmarshal(db.FunctionCalls, &functionCalls); err == nil {
			ch.FunctionCalls = functionCalls
		}
	}

	// Handle soft-delete
	if db.DeletedAt.Valid {
		ch.DeletedAt = &db.DeletedAt.Time
	}

	return ch
}

// ChatHistoriesFromDB converts a slice of database ChatHistoryDB models to domain ChatHistory slice.
func ChatHistoriesFromDB(dbList []model.ChatHistoryDB) []domain.ChatHistory {
	histories := make([]domain.ChatHistory, len(dbList))
	for i, db := range dbList {
		histories[i] = ChatHistoryFromDB(db)
	}
	return histories
}
