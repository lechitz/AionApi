package usecase_test

import (
	"errors"
	"testing"
	"time"

	"github.com/lechitz/aion-api/internal/record/core/domain"
	"github.com/lechitz/aion-api/internal/record/core/usecase"
	"github.com/lechitz/aion-api/tests/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_GetByID_Success_CacheHit(t *testing.T) {
	suite := setup.RecordServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	recordID := uint64(100)

	cachedRecord := domain.Record{
		ID:        recordID,
		UserID:    userID,
		TagID:     10,
		EventTime: time.Now().UTC(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	// Mock: Cache hit
	suite.RecordCache.EXPECT().
		GetRecord(gomock.Any(), recordID, userID).
		Return(cachedRecord, nil)

	// Execute
	result, err := suite.RecordService.GetByID(suite.Ctx, recordID, userID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, recordID, result.ID)
	assert.Equal(t, userID, result.UserID)
}

func TestService_GetByID_Success_CacheMiss(t *testing.T) {
	suite := setup.RecordServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	recordID := uint64(100)

	dbRecord := domain.Record{
		ID:        recordID,
		UserID:    userID,
		TagID:     10,
		EventTime: time.Now().UTC(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	// Mock: Cache miss
	suite.RecordCache.EXPECT().
		GetRecord(gomock.Any(), recordID, userID).
		Return(domain.Record{}, errors.New("cache miss"))

	// Mock: Repository returns record
	suite.RecordRepository.EXPECT().
		GetByID(gomock.Any(), recordID, userID).
		Return(dbRecord, nil)

	// Mock: Save to cache (best-effort)
	suite.RecordCache.EXPECT().
		SaveRecord(gomock.Any(), dbRecord, gomock.Any()).
		Return(nil)

	// Execute
	result, err := suite.RecordService.GetByID(suite.Ctx, recordID, userID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, recordID, result.ID)
	assert.Equal(t, userID, result.UserID)
}

func TestService_GetByID_Errors(t *testing.T) {
	tests := []struct {
		name      string
		recordID  uint64
		userID    uint64
		setupMock func(*setup.RecordServiceTestSuite, uint64, uint64)
		wantErr   error
	}{
		{
			name:     "error - record not found",
			recordID: 999,
			userID:   1,
			setupMock: func(s *setup.RecordServiceTestSuite, recordID, userID uint64) {
				s.RecordCache.EXPECT().
					GetRecord(gomock.Any(), recordID, userID).
					Return(domain.Record{}, errors.New("cache miss"))

				s.RecordRepository.EXPECT().
					GetByID(gomock.Any(), recordID, userID).
					Return(domain.Record{}, errors.New("not found"))
			},
			wantErr: usecase.ErrGetRecord,
		},
		{
			name:     "error - record belongs to different user",
			recordID: 100,
			userID:   999,
			setupMock: func(s *setup.RecordServiceTestSuite, recordID, userID uint64) {
				s.RecordCache.EXPECT().
					GetRecord(gomock.Any(), recordID, userID).
					Return(domain.Record{}, errors.New("cache miss"))

				s.RecordRepository.EXPECT().
					GetByID(gomock.Any(), recordID, userID).
					Return(domain.Record{}, errors.New("unauthorized"))
			},
			wantErr: usecase.ErrGetRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := setup.RecordServiceTest(t)
			defer suite.Ctrl.Finish()

			tt.setupMock(suite, tt.recordID, tt.userID)

			// Execute
			result, err := suite.RecordService.GetByID(suite.Ctx, tt.recordID, tt.userID)

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr.Error())
			assert.Equal(t, domain.Record{}, result)
		})
	}
}
