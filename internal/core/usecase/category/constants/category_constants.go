// Package constants contains constants related to category operations.
package constants

// ErrToValidateCategory indicates a validation error in a category operation.
const ErrToValidateCategory = "category validation error"

// CategoryAlreadyExists is returned when trying to create a category that already exists.
const CategoryAlreadyExists = "category already exists"

// FailedToCreateCategory indicates failure to create a category.
const FailedToCreateCategory = "failed to create category"

// FailedToUpdateCategory indicates failure to update a category.
const FailedToUpdateCategory = "failed to update category"

// SuccessfullyCreatedCategory is used to format a success message when a category is created.
const SuccessfullyCreatedCategory = "successfully created category %s"

// SuccessfullyUpdatedCategory is used to format a success message when a category is updated.
const SuccessfullyUpdatedCategory = "successfully updated category %s"

// SuccessfullySoftDeletedCategory is used to format a success message when a category is soft deleted.
const SuccessfullySoftDeletedCategory = "successfully soft deleted category %s"

// FailedToGetCategoryByID indicates failure to retrieve a category by its ID.
const FailedToGetCategoryByID = "failed to get category by id"

// FailedToGetCategoryByName indicates failure to retrieve a category by its name.
const FailedToGetCategoryByName = "failed to get category by name"

// FailedToGetAllCategories indicates failure to retrieve all categories.
const FailedToGetAllCategories = "failed to get all categories"

// FailedToSoftDeleteCategory indicates failure to soft delete a category.
const FailedToSoftDeleteCategory = "failed to soft delete category"

// CategoryIDIsRequired indicates the category ID is required.
const CategoryIDIsRequired = "category id is required"

// CategoryNameIsRequired indicates the category name is required.
const CategoryNameIsRequired = "category name is required"

// CategoryDescriptionIsTooLong indicates the category description is too long.
const CategoryDescriptionIsTooLong = "category description cannot exceed 200 characters"

// CategoryColorIsTooLong indicates the category color is too long.
const CategoryColorIsTooLong = "category color cannot exceed 7 characters"

// CategoryIconIsTooLong indicates the category icon is too long.
const CategoryIconIsTooLong = "category icon cannot exceed 50 characters"

// TODO: separar o que foi de pra common

// CategoryName is the key for category name in context or responses.
const CategoryName = "name"

// CategoryDescription is the key for category description in context or responses.
const CategoryDescription = "description"

// CategoryColor is the key for category color in context or responses.
const CategoryColor = "color_hex"

// CategoryIcon is the key for category icon in context or responses.
const CategoryIcon = "icon"
