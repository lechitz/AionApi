package controller_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSoftDelete_ServiceError(t *testing.T) {
	expected := errors.New("delete failed")
	svc := &recordServiceStub{
		deleteFn: func(_ context.Context, _, _ uint64) error {
			return expected
		},
	}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	err := h.SoftDelete(t.Context(), 1, 2)
	require.ErrorIs(t, err, expected)
}

func TestSoftDelete_Success(t *testing.T) {
	svc := &recordServiceStub{
		deleteFn: func(_ context.Context, recordID uint64, userID uint64) error {
			require.Equal(t, uint64(1), recordID)
			require.Equal(t, uint64(2), userID)
			return nil
		},
	}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	err := h.SoftDelete(t.Context(), 1, 2)
	require.NoError(t, err)
}

func TestSoftDeleteAll_ServiceError(t *testing.T) {
	expected := errors.New("delete failed")
	svc := &recordServiceStub{
		deleteAllFn: func(_ context.Context, _ uint64) error {
			return expected
		},
	}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	err := h.SoftDeleteAll(t.Context(), 2)
	require.ErrorIs(t, err, expected)
}

func TestSoftDeleteAll_Success(t *testing.T) {
	svc := &recordServiceStub{
		deleteAllFn: func(_ context.Context, userID uint64) error {
			require.Equal(t, uint64(2), userID)
			return nil
		},
	}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	err := h.SoftDeleteAll(t.Context(), 2)
	require.NoError(t, err)
}
