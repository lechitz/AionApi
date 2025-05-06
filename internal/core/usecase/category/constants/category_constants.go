package constants

const (
	ErrToValidateCategory        = "category validation error"
	CategoryAlreadyExists        = "category already exists"
	FailedToCreateCategory       = "failed to create category"
	FailedToUpdateCategory       = "failed to update category"
	SuccessfullyCreatedCategory  = "successfully created category %s"
	SuccessfullyUpdatedCategory  = "successfully updated category %s"
	FailedToGetCategoryByID      = "failed to get category by id"
	FailedToGetCategoryByName    = "failed to get category by name"
	FailedToGetAllCategories     = "failed to get all categories"
	CategoryIDIsRequired         = "category id is required"
	CategoryNameIsRequired       = "category name is required"
	CategoryDescriptionIsTooLong = "category description cannot exceed 200 characters"
	CategoryColorIsTooLong       = "category color cannot exceed 7 characters"
	CategoryIconIsTooLong        = "category icon cannot exceed 50 characters"
)

const (
	Error               = "error"
	CategoryID          = "category_id"
	CategoryName        = "name"
	CategoryDescription = "description"
	CategoryColor       = "color_hex"
	CategoryIcon        = "icon"
	Category            = "category"
)
