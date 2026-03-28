package controller_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lechitz/aion-api/internal/record/adapter/primary/graphql/controller"
	"github.com/lechitz/aion-api/internal/record/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListByUser_UserIDMissing(t *testing.T) {
	svc := &recordServiceStub{}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.ListByUser(t.Context(), 0, 10, nil, nil)
	require.ErrorIs(t, err, controller.ErrUserIDNotFound)
	assert.Nil(t, out)
}

func TestListByUser_ServiceError(t *testing.T) {
	expected := errors.New("list failed")
	svc := &recordServiceStub{
		listByUserFn: func(_ context.Context, _ uint64, _ int, _ *string, _ *int64) ([]domain.Record, error) {
			return nil, expected
		},
	}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.ListByUser(t.Context(), 1, 10, nil, nil)
	require.ErrorIs(t, err, expected)
	assert.Nil(t, out)
}

func TestListByUser_Success(t *testing.T) {
	when := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	svc := &recordServiceStub{
		listByUserFn: func(_ context.Context, userID uint64, limit int, _ *string, _ *int64) ([]domain.Record, error) {
			require.Equal(t, uint64(9), userID)
			require.Equal(t, 10, limit)
			return []domain.Record{{ID: 1, UserID: userID, TagID: 2, EventTime: when, CreatedAt: when, UpdatedAt: when}}, nil
		},
	}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.ListByUser(t.Context(), 9, 10, nil, nil)
	require.NoError(t, err)
	require.Len(t, out, 1)
	assert.Equal(t, "1", out[0].ID)
}
