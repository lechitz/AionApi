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

func TestListAllBetween_UserIDMissing(t *testing.T) {
	svc := &recordServiceStub{}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.ListAllBetween(t.Context(), 0, "2024-01-01T00:00:00Z", "2024-01-02T00:00:00Z", 10)
	require.ErrorIs(t, err, controller.ErrUserIDNotFound)
	assert.Nil(t, out)
}

func TestListAllBetween_InvalidStartDate(t *testing.T) {
	svc := &recordServiceStub{}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.ListAllBetween(t.Context(), 1, "bad", "2024-01-02T00:00:00Z", 10)
	require.Error(t, err)
	assert.Equal(t, "invalid start date, expected RFC3339", err.Error())
	assert.Nil(t, out)
}

func TestListAllBetween_InvalidEndDate(t *testing.T) {
	svc := &recordServiceStub{}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.ListAllBetween(t.Context(), 1, "2024-01-01T00:00:00Z", "bad", 10)
	require.Error(t, err)
	assert.Equal(t, "invalid end date, expected RFC3339", err.Error())
	assert.Nil(t, out)
}

func TestListAllBetween_ServiceError(t *testing.T) {
	expected := errors.New("list failed")
	svc := &recordServiceStub{
		listAllBetweenFn: func(_ context.Context, _ uint64, _ time.Time, _ time.Time, _ int) ([]domain.Record, error) {
			return nil, expected
		},
	}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.ListAllBetween(t.Context(), 1, "2024-01-01T00:00:00Z", "2024-01-02T00:00:00Z", 10)
	require.ErrorIs(t, err, expected)
	assert.Nil(t, out)
}

func TestListAllBetween_Success(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	svc := &recordServiceStub{
		listAllBetweenFn: func(_ context.Context, userID uint64, gotStart time.Time, gotEnd time.Time, limit int) ([]domain.Record, error) {
			require.Equal(t, uint64(9), userID)
			require.Equal(t, start, gotStart)
			require.Equal(t, end, gotEnd)
			require.Equal(t, 10, limit)
			return []domain.Record{{ID: 1, UserID: userID, TagID: 2, EventTime: start, CreatedAt: start, UpdatedAt: start}}, nil
		},
	}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.ListAllBetween(t.Context(), 9, start.Format(time.RFC3339), end.Format(time.RFC3339), 10)
	require.NoError(t, err)
	require.Len(t, out, 1)
	assert.Equal(t, "1", out[0].ID)
}
