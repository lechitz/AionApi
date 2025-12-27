// Package setup provides test suite builders and common helpers for unit tests.
package setup

import (
	"context"
	"testing"

	"github.com/lechitz/AionApi/internal/record/core/usecase"
	"github.com/lechitz/AionApi/tests/mocks"
	"go.uber.org/mock/gomock"
)

// RecordServiceTestSuite groups mocked dependencies and the SUT (RecordService)
// to simplify Record-related unit tests.
type RecordServiceTestSuite struct {
	Ctrl             *gomock.Controller
	Logger           *mocks.MockContextLogger
	RecordRepository *mocks.MockRecordRepository
	RecordCache      *mocks.MockRecordCache
	TagRepository    *mocks.MockTagRepository
	RecordService    *usecase.Service
	Ctx              context.Context
}

// RecordServiceTest initializes and returns a RecordServiceTestSuite with the
// correct mocked output ports (RecordRepository, RecordCache, TagRepository, ContextLogger). Use this helper to
// bootstrap each test and ensure proper teardown via Ctrl.Finish().
func RecordServiceTest(t *testing.T) *RecordServiceTestSuite {
	ctrl := gomock.NewController(t)

	recordRepository := mocks.NewMockRecordRepository(ctrl)
	recordCache := mocks.NewMockRecordCache(ctrl)
	tagRepository := mocks.NewMockTagRepository(ctrl)
	logger := mocks.NewMockContextLogger(ctrl)

	// Set default, non-intrusive expectations for the logger.
	ExpectLoggerDefaultBehavior(logger)

	svc := usecase.NewService(recordRepository, recordCache, tagRepository, logger)

	return &RecordServiceTestSuite{
		Ctrl:             ctrl,
		Logger:           logger,
		RecordRepository: recordRepository,
		RecordCache:      recordCache,
		TagRepository:    tagRepository,
		RecordService:    svc,
		Ctx:              t.Context(),
	}
}
