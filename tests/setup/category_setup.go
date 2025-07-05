// Package setup contains test setup utilities for testing Service.
package setup

import (
	"context"
	"testing"

	"github.com/lechitz/AionApi/internal/core/usecase/category"
	"github.com/lechitz/AionApi/tests/mocks"
	"go.uber.org/mock/gomock"
)

// CategoryServiceTestSuite is a test suite structure for testing methods in the CategoryService, holding mock dependencies and context.
type CategoryServiceTestSuite struct {
	Ctrl               *gomock.Controller
	Logger             *mocks.MockLogger
	CategoryRepository *mocks.MockCategoryStore
	CategoryService    *category.Service
	Ctx                context.Context
}

// CategoryServiceTest initializes and returns a CategoryServiceTestSuite with mock dependencies for testing Service logic.
func CategoryServiceTest(t *testing.T) *CategoryServiceTestSuite {
	ctrl := gomock.NewController(t)

	mockCategoryRepository := mocks.NewMockCategoryStore(ctrl)
	mockLog := mocks.NewMockLogger(ctrl)

	ExpectLoggerDefaultBehavior(mockLog)

	categoryService := category.NewCategoryService(mockCategoryRepository, mockLog)

	return &CategoryServiceTestSuite{
		Ctrl:               ctrl,
		Logger:             mockLog,
		CategoryService:    categoryService,
		CategoryRepository: mockCategoryRepository,
		Ctx:                t.Context(),
	}
}
