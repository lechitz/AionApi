package controller

import (
	"context"

	"github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/internal/record/core/ports/input"
)

// RecordController is the contract used by GraphQL resolvers.
type RecordController interface {
	Create(ctx context.Context, in model.CreateRecordInput, userID uint64) (*model.Record, error)
	GetByID(ctx context.Context, recordID, userID uint64) (*model.Record, error)
	ListByUser(ctx context.Context, userID uint64, limit int, afterEventTime *string, afterID *int64) ([]*model.Record, error)
	ListByTag(ctx context.Context, tagID, userID uint64, limit int) ([]*model.Record, error)
	ListByDay(ctx context.Context, userID uint64, date string) ([]*model.Record, error)
	ListAllUntil(ctx context.Context, userID uint64, until string, limit int) ([]*model.Record, error)
	ListAllBetween(ctx context.Context, userID uint64, startDate, endDate string, limit int) ([]*model.Record, error)
	Update(ctx context.Context, in model.UpdateRecordInput, userID uint64) (*model.Record, error)
	SoftDelete(ctx context.Context, recordID, userID uint64) error
	SoftDeleteAll(ctx context.Context, userID uint64) error
}

// controller is the controller for the record service.
type controller struct {
	RecordService input.RecordService
	Logger        logger.ContextLogger
}

// NewController wires dependencies and returns a Controller.
func NewController(svc input.RecordService, logger logger.ContextLogger) RecordController {
	return &controller{
		RecordService: svc,
		Logger:        logger,
	}
}
