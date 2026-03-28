// Package setup provides test suite builders and common helpers for unit tests.
package setup

import (
	"context"
	"testing"

	"github.com/lechitz/aion-api/internal/tag/core/usecase"
	"github.com/lechitz/aion-api/tests/mocks"
	"go.uber.org/mock/gomock"
)

// TagServiceTestSuite groups mocked dependencies and the SUT (TagService)
// to simplify Tag-related unit tests.
type TagServiceTestSuite struct {
	Ctrl          *gomock.Controller
	Logger        *mocks.MockContextLogger
	TagRepository *mocks.MockTagRepository
	TagCache      *mocks.MockTagCache
	TagService    *usecase.Service
	Ctx           context.Context
}

// TagServiceTest initializes and returns a TagServiceTestSuite with the
// correct mocked output ports (TagRepository, TagCache, ContextLogger). Use this helper to
// bootstrap each test and ensure proper teardown via Ctrl.Finish().
func TagServiceTest(t *testing.T) *TagServiceTestSuite {
	ctrl := gomock.NewController(t)

	tagRepository := mocks.NewMockTagRepository(ctrl)
	tagCache := mocks.NewMockTagCache(ctrl)
	logger := mocks.NewMockContextLogger(ctrl)

	// Set default, non-intrusive expectations for the logger.
	ExpectLoggerDefaultBehavior(logger)

	svc := usecase.NewService(tagRepository, tagCache, logger)

	return &TagServiceTestSuite{
		Ctrl:          ctrl,
		Logger:        logger,
		TagRepository: tagRepository,
		TagCache:      tagCache,
		TagService:    svc,
		Ctx:           t.Context(),
	}
}
