package constants

const (
	// Erros de negócios relacionados a usuários

	ErrorToCreateUser         = "error to create user"
	ErrorToGetUser            = "error to get user"
	ErrorToGetUsers           = "error to get users"
	ErrorToUpdateUser         = "error to update user"
	ErrorToDeleteUser         = "error to delete user"
	ErrorToExtractUserID      = "error to extract user ID"
	ErrorUserPermissionDenied = "user permission denied"
	ErrorToFormatUpdateUser   = "error to format update user"
	ErrorToFormatCreateUser   = "error to format create user"
	ErrorToValidateCreateUser = "error to validate user"
	ErrorToGetUserByID        = "error to get user by ID"
	ErrorToGetUserByUserName  = "error to get user by username"

	// Erros de validação de dados de usuários

	ErrorMissingFields = "missing fields"
	NameIsRequired     = "name is required"
	UsernameIsRequired = "username is required"
	EmailIsRequired    = "email is required"
	PasswordIsRequired = "password is required"
	UserIDIsRequired   = "user ID is required"
	InvalidEmail       = "invalid email"

	// Mensagens de Sucesso

	SuccessToCreateUser = "user created successfully"
	SuccessToGetUser    = "user get successfully"
	SuccessToGetUsers   = "users get successfully"
	SuccessToUpdateUser = "user updated successfully"
	SuccessToDeleteUser = "user deleted successfully"
)
