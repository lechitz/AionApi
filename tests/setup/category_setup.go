// Package setup provides test suite builders and common helpers for unit tests.
package setup

import (
	"context"
	"testing"

	"github.com/lechitz/AionApi/internal/core/usecase/category"
	"github.com/lechitz/AionApi/tests/mocks"
	"go.uber.org/mock/gomock"
)

// CategoryServiceTestSuite groups mocked dependencies and the SUT (CategoryService)
// to simplify Category-related unit tests.
type CategoryServiceTestSuite struct {
	Ctrl               *gomock.Controller
	Logger             *mocks.ContextLogger
	CategoryRepository *mocks.CategoryRepository
	CategoryService    *category.Service
	Ctx                context.Context
}

// CategoryServiceTest initializes and returns a CategoryServiceTestSuite with the
// correct mocked output ports (CategoryStore, ContextLogger). Use this helper to
// bootstrap each test and ensure proper teardown via Ctrl.Finish().
func CategoryServiceTest(t *testing.T) *CategoryServiceTestSuite {
	ctrl := gomock.NewController(t)

	categoryRepository := mocks.NewCategoryRepository(ctrl)
	logger := mocks.NewContextLogger(ctrl)

	// Set default, non-intrusive expectations for the logger.
	ExpectLoggerDefaultBehavior(logger)

	svc := category.NewService(categoryRepository, logger)

	return &CategoryServiceTestSuite{
		Ctrl:               ctrl,
		Logger:             logger,
		CategoryRepository: categoryRepository,
		CategoryService:    svc,
		Ctx:                t.Context(),
	}
}
