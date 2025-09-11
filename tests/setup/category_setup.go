// Package setup provides test suite builders and common helpers for unit tests.
package setup

import (
	"context"
	"testing"

	"github.com/lechitz/AionApi/internal/category/core/usecase"
	"github.com/lechitz/AionApi/tests/mocks"
	"go.uber.org/mock/gomock"
)

// CategoryServiceTestSuite groups mocked dependencies and the SUT (CategoryService)
// to simplify Category-related unit tests.
type CategoryServiceTestSuite struct {
	Ctrl               *gomock.Controller
	Logger             *mocks.MockContextLogger
	CategoryRepository *mocks.MockCategoryRepository
	CategoryService    *usecase.Service
	Ctx                context.Context
}

// CategoryServiceTest initializes and returns a CategoryServiceTestSuite with the
// correct mocked output ports (CategoryStore, ContextLogger). Use this helper to
// bootstrap each test and ensure proper teardown via Ctrl.Finish().
func CategoryServiceTest(t *testing.T) *CategoryServiceTestSuite {
	ctrl := gomock.NewController(t)

	categoryRepository := mocks.NewMockCategoryRepository(ctrl)
	logger := mocks.NewMockContextLogger(ctrl)

	// Set default, non-intrusive expectations for the logger.
	ExpectLoggerDefaultBehavior(logger)

	svc := usecase.NewService(categoryRepository, logger)

	return &CategoryServiceTestSuite{
		Ctrl:               ctrl,
		Logger:             logger,
		CategoryRepository: categoryRepository,
		CategoryService:    svc,
		Ctx:                t.Context(),
	}
}
