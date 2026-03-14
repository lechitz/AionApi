package usecase_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/category/core/usecase"
	"go.uber.org/mock/gomock"

	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/assert"
)

func TestSoftDeleteCategory_ErrorToSoftDeleteCategory(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	categoryID := uint64(1)
	userID := uint64(3)

	suite.CategoryRepository.EXPECT().
		SoftDelete(gomock.Any(), categoryID, userID).
		Return(errors.New(usecase.FailedToSoftDeleteCategory))

	err := suite.CategoryService.SoftDelete(suite.Ctx, categoryID, userID)
	assert.EqualError(t, err, usecase.FailedToSoftDeleteCategory)
}

func TestSoftDeleteCategory_Success(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	categoryID := uint64(1)
	userID := uint64(3)

	suite.CategoryRepository.EXPECT().
		SoftDelete(gomock.Any(), categoryID, userID).
		Return(nil)

	suite.CategoryCache.EXPECT().
		DeleteCategory(gomock.Any(), categoryID, userID).
		Return(nil)

	suite.CategoryCache.EXPECT().
		DeleteCategoryList(gomock.Any(), userID).
		Return(nil)

	err := suite.CategoryService.SoftDelete(suite.Ctx, categoryID, userID)
	assert.NoError(t, err)
}
