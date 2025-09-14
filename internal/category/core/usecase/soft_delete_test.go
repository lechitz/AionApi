package usecase_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/category/core/usecase"
	"go.uber.org/mock/gomock"

	"github.com/lechitz/AionApi/tests/setup"
	"github.com/lechitz/AionApi/tests/testdata"
	"github.com/stretchr/testify/assert"
)

func TestSoftDeleteCategory_ErrorToSoftDeleteCategory(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	suite.CategoryRepository.EXPECT().
		SoftDelete(gomock.Any(), testdata.PerfectCategory.ID, testdata.PerfectCategory.UserID).
		Return(errors.New(usecase.FailedToSoftDeleteCategory))

	err := suite.CategoryService.SoftDelete(suite.Ctx, testdata.PerfectCategory.ID, testdata.PerfectCategory.UserID)
	assert.EqualError(t, err, usecase.FailedToSoftDeleteCategory)
}

func TestSoftDeleteCategory_Success(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory

	suite.CategoryRepository.EXPECT().
		SoftDelete(gomock.Any(), category.ID, category.UserID).
		Return(nil)

	err := suite.CategoryService.SoftDelete(suite.Ctx, category.ID, category.UserID)
	assert.NoError(t, err)
}
