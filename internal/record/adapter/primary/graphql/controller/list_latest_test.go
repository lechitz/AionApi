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

func TestListLatest_UserIDMissing(t *testing.T) {
	svc := &recordServiceStub{}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.ListLatest(t.Context(), 0, 10)
	require.ErrorIs(t, err, controller.ErrUserIDNotFound)
	assert.Nil(t, out)
}

func TestListLatest_DefaultLimitApplied(t *testing.T) {
	svc := &recordServiceStub{
		listLatestFn: func(_ context.Context, _ uint64, limit int) ([]domain.Record, error) {
			require.Equal(t, 10, limit)
			return nil, nil
		},
	}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	_, err := h.ListLatest(t.Context(), 1, 0)
	require.NoError(t, err)
}

func TestListLatest_ServiceError(t *testing.T) {
	expected := errors.New("list failed")
	svc := &recordServiceStub{
		listLatestFn: func(_ context.Context, _ uint64, _ int) ([]domain.Record, error) {
			return nil, expected
		},
	}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.ListLatest(t.Context(), 1, 10)
	require.ErrorIs(t, err, expected)
	assert.Nil(t, out)
}

func TestListLatest_Success(t *testing.T) {
	when := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	svc := &recordServiceStub{
		listLatestFn: func(_ context.Context, userID uint64, limit int) ([]domain.Record, error) {
			require.Equal(t, uint64(9), userID)
			require.Equal(t, 5, limit)
			return []domain.Record{{ID: 1, UserID: userID, TagID: 2, EventTime: when, CreatedAt: when, UpdatedAt: when}}, nil
		},
	}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.ListLatest(t.Context(), 9, 5)
	require.NoError(t, err)
	require.Len(t, out, 1)
	assert.Equal(t, "1", out[0].ID)
}
