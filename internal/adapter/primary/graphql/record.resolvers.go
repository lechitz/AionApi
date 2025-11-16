package graphql

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
)

// CreateRecord is the resolver for the createRecord field.
func (m *mutationResolver) CreateRecord(ctx context.Context, input model.CreateRecordInput) (*model.Record, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return m.RecordController().Create(ctx, input, uid)
}

// RecordByID is the resolver for the recordByID field.
func (q *queryResolver) RecordByID(ctx context.Context, recordID string) (*model.Record, error) {
	id, err := strconv.ParseUint(recordID, 10, 64)
	if err != nil {
		return nil, err
	}

	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return q.RecordController().GetByID(ctx, id, uid)
}

// RecordsByTag is the resolver for the recordsByTag field.
func (q *queryResolver) RecordsByTag(ctx context.Context, tagID string, limit *int32) ([]*model.Record, error) {
	tid, err := strconv.ParseUint(tagID, 10, 64)
	if err != nil {
		return nil, err
	}

	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)

	lim := 50 // default
	if limit != nil && *limit > 0 {
		lim = int(*limit)
	}

	return q.RecordController().ListByTag(ctx, tid, uid, lim)
}

// RecordsByDay is the resolver for the recordsByDay field.
func (q *queryResolver) RecordsByDay(ctx context.Context, date string) ([]*model.Record, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return q.RecordController().ListByDay(ctx, uid, date)
}

// RecordsUntil is the resolver for the recordsUntil field.
func (q *queryResolver) RecordsUntil(ctx context.Context, until string, limit *int32) ([]*model.Record, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)

	lim := 50 // default
	if limit != nil && *limit > 0 {
		lim = int(*limit)
	}

	return q.RecordController().ListAllUntil(ctx, uid, until, lim)
}

// RecordsBetween is the resolver for the recordsBetween field.
func (q *queryResolver) RecordsBetween(ctx context.Context, startDate string, endDate string, limit *int32) ([]*model.Record, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)

	lim := 50 // default
	if limit != nil && *limit > 0 {
		lim = int(*limit)
	}

	return q.RecordController().ListAllBetween(ctx, uid, startDate, endDate, lim)
}

// Records is the resolver for the records field (list by user with optional cursors).
func (q *queryResolver) Records(ctx context.Context, limit *int32, afterEventTime *string, afterID *string) ([]*model.Record, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	lim := 50
	if limit != nil && *limit > 0 {
		lim = int(*limit)
	}
	var afterIDInt *int64
	if afterID != nil && *afterID != "" {
		if v, err := strconv.ParseInt(*afterID, 10, 64); err == nil {
			afterIDInt = &v
		}
	}

	return q.RecordController().ListByUser(ctx, uid, lim, afterEventTime, afterIDInt)
}

// UpdateRecord is the resolver for the updateRecord field.
func (m *mutationResolver) UpdateRecord(ctx context.Context, input model.UpdateRecordInput) (*model.Record, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return m.RecordController().Update(ctx, input, uid)
}

// SoftDeleteRecord is the resolver for the softDeleteRecord field.
func (m *mutationResolver) SoftDeleteRecord(ctx context.Context, input model.DeleteRecordInput) (bool, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	id, err := strconv.ParseUint(input.ID, 10, 64)
	if err != nil {
		return false, err
	}
	if err := m.RecordController().SoftDelete(ctx, id, uid); err != nil {
		return false, err
	}
	return true, nil
}

// SoftDeleteAllRecords is the resolver for the softDeleteAllRecords field.
func (m *mutationResolver) SoftDeleteAllRecords(ctx context.Context) (bool, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	if err := m.RecordController().SoftDeleteAll(ctx, uid); err != nil {
		return false, err
	}
	return true, nil
}

// Additional unimplemented resolvers can be added below as needed.
