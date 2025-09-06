package usecase_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/core/category/domain"
	"github.com/lechitz/AionApi/internal/core/category/usecase"
	"go.uber.org/mock/gomock"

	"github.com/lechitz/AionApi/tests/setup"
	"github.com/lechitz/AionApi/tests/testdata"
	"github.com/stretchr/testify/assert"
)

func TestSoftDeleteCategory_ErrorToSoftDeleteCategory(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	suite.CategoryRepository.EXPECT().
		SoftDelete(gomock.Any(), testdata.PerfectCategory).
		Return(errors.New(usecase.FailedToSoftDeleteCategory))

	err := suite.CategoryService.SoftDelete(suite.Ctx, testdata.PerfectCategory)
	assert.EqualError(t, err, usecase.FailedToSoftDeleteCategory)
}

func TestSoftDeleteCategory_Success(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.Category{ID: testdata.PerfectCategory.ID}

	suite.CategoryRepository.EXPECT().
		SoftDelete(gomock.Any(), input).
		Return(nil)

	err := suite.CategoryService.SoftDelete(suite.Ctx, input)
	assert.NoError(t, err)
}
