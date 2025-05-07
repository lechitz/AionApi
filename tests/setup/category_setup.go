package setup

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/lechitz/AionApi/internal/core/usecase/category"
	mockCategory "github.com/lechitz/AionApi/tests/mocks/category"
	mockLogger "github.com/lechitz/AionApi/tests/mocks/logger"
	"testing"
)

type CategoryServiceTestSuite struct {
	Ctrl               *gomock.Controller
	Logger             *mockLogger.MockLogger
	CategoryRepository *mockCategory.MockCategoryStore
	CategoryService    *category.CategoryService
	Ctx                context.Context
}

func SetupCategoryServiceTest(t *testing.T) *CategoryServiceTestSuite {
	ctrl := gomock.NewController(t)

	mockCategoryRepository := mockCategory.NewMockCategoryStore(ctrl)
	mockLog := mockLogger.NewMockLogger(ctrl)

	ExpectLoggerDefaultBehavior(mockLog)

	categoryService := category.NewCategoryService(mockCategoryRepository, mockLog)

	return &CategoryServiceTestSuite{
		Ctrl:               ctrl,
		Logger:             mockLog,
		CategoryService:    categoryService,
		CategoryRepository: mockCategoryRepository,
		Ctx:                context.Background(),
	}
}
