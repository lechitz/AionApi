package category_test

import (
	"errors"
	"github.com/lechitz/AionApi/internal/core/domain"
	"testing"

	"github.com/lechitz/AionApi/internal/core/usecase/category/constants"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/lechitz/AionApi/tests/testdata"
	"github.com/stretchr/testify/assert"
)

func TestSoftDeleteCategory_ErrorToGetCategoryByID(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	suite.CategoryRepository.EXPECT().
		GetCategoryByID(suite.Ctx, testdata.PerfectCategory).
		Return(domain.Category{}, errors.New(constants.FailedToGetCategoryByID))

	err := suite.CategoryService.SoftDeleteCategory(suite.Ctx, testdata.PerfectCategory)
	assert.EqualError(t, err, constants.FailedToGetCategoryByID)
}

func TestSoftDeleteCategory_ErrorToSoftDeleteCategory(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	suite.CategoryRepository.EXPECT().
		GetCategoryByID(suite.Ctx, testdata.PerfectCategory).
		Return(domain.Category{ID: testdata.PerfectCategory.ID}, nil)

	suite.CategoryRepository.EXPECT().
		SoftDeleteCategory(suite.Ctx, domain.Category{ID: testdata.PerfectCategory.ID}).
		Return(errors.New(constants.FailedToSoftDeleteCategory))

	err := suite.CategoryService.SoftDeleteCategory(suite.Ctx, testdata.PerfectCategory)
	assert.EqualError(t, err, constants.FailedToSoftDeleteCategory)
}

func TestSoftDeleteCategory_Success(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	suite.CategoryRepository.EXPECT().
		GetCategoryByID(suite.Ctx, domain.Category{ID: testdata.PerfectCategory.ID}).
		Return(domain.Category{ID: testdata.PerfectCategory.ID}, nil)

	suite.CategoryRepository.EXPECT().
		SoftDeleteCategory(suite.Ctx, domain.Category{ID: testdata.PerfectCategory.ID}).
		Return(nil)

	err := suite.CategoryService.SoftDeleteCategory(
		suite.Ctx,
		domain.Category{ID: testdata.PerfectCategory.ID},
	)

	assert.NoError(t, err)
}
