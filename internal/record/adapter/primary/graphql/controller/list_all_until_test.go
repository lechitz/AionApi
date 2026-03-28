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

func TestListAllUntil_UserIDMissing(t *testing.T) {
	svc := &recordServiceStub{}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.ListAllUntil(t.Context(), 0, time.Now().UTC().Format(time.RFC3339), 10)
	require.ErrorIs(t, err, controller.ErrUserIDNotFound)
	assert.Nil(t, out)
}

func TestListAllUntil_InvalidTimestamp(t *testing.T) {
	svc := &recordServiceStub{}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.ListAllUntil(t.Context(), 1, "bad", 10)
	require.Error(t, err)
	assert.Equal(t, "invalid until timestamp, expected RFC3339", err.Error())
	assert.Nil(t, out)
}

func TestListAllUntil_ServiceError(t *testing.T) {
	expected := errors.New("list failed")
	svc := &recordServiceStub{
		listAllUntilFn: func(_ context.Context, _ uint64, _ time.Time, _ int) ([]domain.Record, error) {
			return nil, expected
		},
	}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.ListAllUntil(t.Context(), 1, time.Now().UTC().Format(time.RFC3339), 10)
	require.ErrorIs(t, err, expected)
	assert.Nil(t, out)
}

func TestListAllUntil_Success(t *testing.T) {
	until := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	svc := &recordServiceStub{
		listAllUntilFn: func(_ context.Context, userID uint64, gotUntil time.Time, limit int) ([]domain.Record, error) {
			require.Equal(t, uint64(9), userID)
			require.Equal(t, until, gotUntil)
			require.Equal(t, 10, limit)
			return []domain.Record{{ID: 1, UserID: userID, TagID: 2, EventTime: until, CreatedAt: until, UpdatedAt: until}}, nil
		},
	}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.ListAllUntil(t.Context(), 9, until.Format(time.RFC3339), 10)
	require.NoError(t, err)
	require.Len(t, out, 1)
	assert.Equal(t, "1", out[0].ID)
}
