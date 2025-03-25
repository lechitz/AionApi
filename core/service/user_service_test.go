package service

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"

	"github.com/lechitz/AionApi/core/domain"
	"github.com/lechitz/AionApi/core/msg"
	mockdb "github.com/lechitz/AionApi/tests/mocks"
	mocksecurity "github.com/lechitz/AionApi/tests/mocks"
	mocktoken "github.com/lechitz/AionApi/tests/mocks"
)

// Tests to CreateUser

func TestCreateUser_WhenValidData_ShouldCreateUserSuccessfully(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockdb.NewMockIUserRepository(ctrl)
	mockTokenService := mocktoken.NewMockITokenService(ctrl)
	mockPasswordService := mocksecurity.NewMockIPasswordService(ctrl)

	logger := zaptest.NewLogger(t).Sugar()
	userService := NewUserService(mockUserRepo, mockTokenService, mockPasswordService, logger)

	ctx := domain.ContextControl{}
	rawPassword := "123456"
	hashedPassword := "$2a$10$abc123"
	inputUser := domain.UserDomain{
		Name:     "Felipe",
		Username: "lechitz",
		Email:    "lechitz@aion.com",
	}

	expectedUser := domain.UserDomain{
		ID:       42,
		Name:     "Felipe",
		Username: "lechitz",
		Email:    "lechitz@aion.com",
		Password: hashedPassword,
	}

	mockUserRepo.
		EXPECT().
		GetUserByUsername(ctx, inputUser.Username).
		Return(domain.UserDomain{}, errors.New("not found"))

	mockUserRepo.
		EXPECT().
		GetUserByEmail(ctx, inputUser.Email).
		Return(domain.UserDomain{}, errors.New("not found"))

	mockPasswordService.
		EXPECT().
		HashPassword(rawPassword).
		Return(hashedPassword, nil)

	mockUserRepo.
		EXPECT().
		CreateUser(ctx, gomock.AssignableToTypeOf(domain.UserDomain{})).
		DoAndReturn(func(_ domain.ContextControl, user domain.UserDomain) (domain.UserDomain, error) {

			user.ID = expectedUser.ID
			return user, nil
		})

	result, err := userService.CreateUser(ctx, inputUser, rawPassword)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, expectedUser.Username, result.Username)
	assert.Equal(t, expectedUser.Email, result.Email)
	assert.Equal(t, expectedUser.Password, hashedPassword)
}

func TestCreateUser_WhenRequiredFieldsMissing_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockdb.NewMockIUserRepository(ctrl)
	mockTokenService := mocktoken.NewMockITokenService(ctrl)
	mockPasswordService := mocksecurity.NewMockIPasswordService(ctrl)

	logger := zaptest.NewLogger(t).Sugar()

	userService := NewUserService(mockUserRepo, mockTokenService, mockPasswordService, logger)

	ctx := domain.ContextControl{}
	emptyUser := domain.UserDomain{}
	password := ""

	user, err := userService.CreateUser(ctx, emptyUser, password)

	assert.Error(t, err)
	assert.Equal(t, msg.NameIsRequired, err.Error())
	assert.Equal(t, uint64(0), user.ID)
}

func TestCreateUser_WhenUsernameAlreadyExists_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockdb.NewMockIUserRepository(ctrl)
	mockTokenService := mocktoken.NewMockITokenService(ctrl)
	mockPasswordService := mocksecurity.NewMockIPasswordService(ctrl)

	logger := zaptest.NewLogger(t).Sugar()

	userService := NewUserService(mockUserRepo, mockTokenService, mockPasswordService, logger)

	ctx := domain.ContextControl{}
	user := domain.UserDomain{
		Name:     "Test",
		Username: "testuser",
		Email:    "test@example.com",
	}
	password := "123456"

	mockUserRepo.
		EXPECT().
		GetUserByUsername(ctx, "testuser").
		Return(domain.UserDomain{ID: 1}, nil)

	createdUser, err := userService.CreateUser(ctx, user, password)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "username 'testuser' is already in use")
	assert.Equal(t, uint64(0), createdUser.ID)
}

func TestCreateUser_WhenEmailAlreadyExists_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockdb.NewMockIUserRepository(ctrl)
	mockTokenService := mocktoken.NewMockITokenService(ctrl)
	mockPasswordService := mocksecurity.NewMockIPasswordService(ctrl)

	logger := zaptest.NewLogger(t).Sugar()

	userService := NewUserService(mockUserRepo, mockTokenService, mockPasswordService, logger)

	ctx := domain.ContextControl{}
	user := domain.UserDomain{
		Name:     "Test",
		Username: "testuser",
		Email:    "test@example.com",
	}
	password := "123456"

	mockUserRepo.
		EXPECT().
		GetUserByUsername(ctx, "testuser").
		Return(domain.UserDomain{}, errors.New("not found"))

	mockUserRepo.
		EXPECT().
		GetUserByEmail(ctx, "test@example.com").
		Return(domain.UserDomain{ID: 2}, nil)

	createdUser, err := userService.CreateUser(ctx, user, password)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email 'test@example.com' is already in use")
	assert.Equal(t, uint64(0), createdUser.ID)
}

func TestCreateUser_WhenHashPasswordFails_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockdb.NewMockIUserRepository(ctrl)
	mockTokenService := mocktoken.NewMockITokenService(ctrl)
	mockPasswordService := mocksecurity.NewMockIPasswordService(ctrl)

	logger := zaptest.NewLogger(t).Sugar()

	userService := NewUserService(mockUserRepo, mockTokenService, mockPasswordService, logger)

	ctx := domain.ContextControl{}
	user := domain.UserDomain{
		Name:     "Test",
		Username: "testuser",
		Email:    "test@example.com",
	}
	password := "123456"

	mockUserRepo.
		EXPECT().
		GetUserByUsername(ctx, "testuser").
		Return(domain.UserDomain{}, errors.New("not found"))

	mockUserRepo.
		EXPECT().
		GetUserByEmail(ctx, "test@example.com").
		Return(domain.UserDomain{}, errors.New("not found"))

	mockPasswordService.
		EXPECT().
		HashPassword(password).
		Return("", errors.New("hash error"))

	createdUser, err := userService.CreateUser(ctx, user, password)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "hash error")
	assert.Equal(t, uint64(0), createdUser.ID)
}

func TestCreateUser_WhenCreateUserFailsInRepository_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockdb.NewMockIUserRepository(ctrl)
	mockTokenService := mocktoken.NewMockITokenService(ctrl)
	mockPasswordService := mocksecurity.NewMockIPasswordService(ctrl)

	logger := zaptest.NewLogger(t).Sugar()

	userService := NewUserService(mockUserRepo, mockTokenService, mockPasswordService, logger)

	ctx := domain.ContextControl{}
	user := domain.UserDomain{
		Name:     "Test",
		Username: "testuser",
		Email:    "test@example.com",
	}
	password := "123456"
	hashed := "$2a$10$abcdef"

	mockUserRepo.
		EXPECT().
		GetUserByUsername(ctx, "testuser").
		Return(domain.UserDomain{}, errors.New("not found"))

	mockUserRepo.
		EXPECT().
		GetUserByEmail(ctx, "test@example.com").
		Return(domain.UserDomain{}, errors.New("not found"))

	mockPasswordService.
		EXPECT().
		HashPassword(password).
		Return(hashed, nil)

	mockUserRepo.
		EXPECT().
		CreateUser(ctx, gomock.Any()).
		Return(domain.UserDomain{}, errors.New("db error"))

	createdUser, err := userService.CreateUser(ctx, user, password)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
	assert.Equal(t, uint64(0), createdUser.ID)
}

// Tests to GetAllUsers

func TestGetAllUsers_WhenRepositoryReturnsUsers_ShouldReturnUserList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockdb.NewMockIUserRepository(ctrl)
	mockTokenService := mocktoken.NewMockITokenService(ctrl)
	mockPasswordService := mocksecurity.NewMockIPasswordService(ctrl)
	logger := zaptest.NewLogger(t).Sugar()

	userService := NewUserService(mockUserRepo, mockTokenService, mockPasswordService, logger)

	ctx := domain.ContextControl{}
	expectedUsers := []domain.UserDomain{
		{ID: 1, Name: "Felipe", Username: "lechitz"},
		{ID: 2, Name: "Bravo", Username: "bravolechitz"},
	}

	mockUserRepo.
		EXPECT().
		GetAllUsers(ctx).
		Return(expectedUsers, nil)

	result, err := userService.GetAllUsers(ctx)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, expectedUsers[0].Username, result[0].Username)
	assert.Equal(t, expectedUsers[1].ID, result[1].ID)
}

func TestGetAllUsers_WhenRepositoryFails_ShouldReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockdb.NewMockIUserRepository(ctrl)
	mockTokenService := mocktoken.NewMockITokenService(ctrl)
	mockPasswordService := mocksecurity.NewMockIPasswordService(ctrl)
	logger := zaptest.NewLogger(t).Sugar()

	userService := NewUserService(mockUserRepo, mockTokenService, mockPasswordService, logger)

	ctx := domain.ContextControl{}
	expectedErr := errors.New("db error")

	mockUserRepo.
		EXPECT().
		GetAllUsers(ctx).
		Return(nil, expectedErr)

	users, err := userService.GetAllUsers(ctx)

	assert.Error(t, err)
	assert.Nil(t, users)
	assert.Equal(t, expectedErr, err)
}

// Tests to GetUserByID

func TestGetUserByID_WhenUserExists_ShouldReturnUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockdb.NewMockIUserRepository(ctrl)
	mockTokenService := mocktoken.NewMockITokenService(ctrl)
	mockPasswordService := mocksecurity.NewMockIPasswordService(ctrl)
	logger := zaptest.NewLogger(t).Sugar()

	userService := NewUserService(mockUserRepo, mockTokenService, mockPasswordService, logger)

	ctx := domain.ContextControl{}
	userID := uint64(42)

	expectedUser := domain.UserDomain{
		ID:       userID,
		Name:     "Felipe",
		Username: "lechitz",
		Email:    "lechitz@aion.com",
	}

	mockUserRepo.
		EXPECT().
		GetUserByID(ctx, userID).
		Return(expectedUser, nil)

	result, err := userService.GetUserByID(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, expectedUser.Username, result.Username)
}

func TestGetUserByID_WhenRepoFails_ShouldReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockdb.NewMockIUserRepository(ctrl)
	mockTokenService := mocktoken.NewMockITokenService(ctrl)
	mockPasswordService := mocksecurity.NewMockIPasswordService(ctrl)
	logger := zaptest.NewLogger(t).Sugar()

	userService := NewUserService(mockUserRepo, mockTokenService, mockPasswordService, logger)

	ctx := domain.ContextControl{}
	userID := uint64(42)
	expectedErr := errors.New("record not found")

	mockUserRepo.
		EXPECT().
		GetUserByID(ctx, userID).
		Return(domain.UserDomain{}, expectedErr)

	result, err := userService.GetUserByID(ctx, userID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, uint64(0), result.ID)
}

// Tests to GetUserByUsername

func TestGetUserByUsername_WhenUserExists_ShouldReturnUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockdb.NewMockIUserRepository(ctrl)
	mockTokenService := mocktoken.NewMockITokenService(ctrl)
	mockPasswordService := mocksecurity.NewMockIPasswordService(ctrl)
	logger := zaptest.NewLogger(t).Sugar()

	userService := NewUserService(mockUserRepo, mockTokenService, mockPasswordService, logger)

	ctx := domain.ContextControl{}
	username := "lechitz"
	expected := domain.UserDomain{ID: 42, Username: username, Name: "Felipe"}

	mockUserRepo.
		EXPECT().
		GetUserByUsername(ctx, username).
		Return(expected, nil)

	user, err := userService.GetUserByUsername(ctx, username)

	assert.NoError(t, err)
	assert.Equal(t, expected.ID, user.ID)
	assert.Equal(t, expected.Username, user.Username)
}

func TestGetUserByUsername_WhenRepositoryFails_ShouldReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockdb.NewMockIUserRepository(ctrl)
	mockTokenService := mocktoken.NewMockITokenService(ctrl)
	mockPasswordService := mocksecurity.NewMockIPasswordService(ctrl)
	logger := zaptest.NewLogger(t).Sugar()

	userService := NewUserService(mockUserRepo, mockTokenService, mockPasswordService, logger)

	ctx := domain.ContextControl{}
	username := "lechitz"
	expectedErr := errors.New("db failure")

	mockUserRepo.
		EXPECT().
		GetUserByUsername(ctx, username).
		Return(domain.UserDomain{}, expectedErr)

	_, err := userService.GetUserByUsername(ctx, username)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

// Tests to GetUserByEmail

func TestGetUserByEmail_WhenUserExists_ShouldReturnUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockdb.NewMockIUserRepository(ctrl)
	mockTokenService := mocktoken.NewMockITokenService(ctrl)
	mockPasswordService := mocksecurity.NewMockIPasswordService(ctrl)
	logger := zaptest.NewLogger(t).Sugar()

	userService := NewUserService(mockUserRepo, mockTokenService, mockPasswordService, logger)

	ctx := domain.ContextControl{}
	email := "lechitz@aion.com"
	expectedUser := domain.UserDomain{
		ID:       99,
		Name:     "Felipe",
		Username: "lechitz",
		Email:    email,
	}

	mockUserRepo.
		EXPECT().
		GetUserByEmail(ctx, email).
		Return(expectedUser, nil)

	result, err := userService.GetUserByEmail(ctx, email)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, result)
}

func TestGetUserByEmail_WhenUserNotFound_ShouldReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockdb.NewMockIUserRepository(ctrl)
	mockTokenService := mocktoken.NewMockITokenService(ctrl)
	mockPasswordService := mocksecurity.NewMockIPasswordService(ctrl)
	logger := zaptest.NewLogger(t).Sugar()

	userService := NewUserService(mockUserRepo, mockTokenService, mockPasswordService, logger)

	ctx := domain.ContextControl{}
	email := "notfound@aion.com"

	mockUserRepo.
		EXPECT().
		GetUserByEmail(ctx, email).
		Return(domain.UserDomain{}, errors.New("not found"))

	result, err := userService.GetUserByEmail(ctx, email)

	assert.Error(t, err)
	assert.Equal(t, uint64(0), result.ID)
}

// Tests to UpdateUser

func TestUpdateUser_WhenValidFields_ShouldUpdateUserSuccessfully(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockdb.NewMockIUserRepository(ctrl)
	mockTokenService := mocktoken.NewMockITokenService(ctrl)
	mockPasswordService := mocksecurity.NewMockIPasswordService(ctrl)
	logger := zaptest.NewLogger(t).Sugar()

	userService := NewUserService(mockUserRepo, mockTokenService, mockPasswordService, logger)

	ctx := domain.ContextControl{}
	inputUser := domain.UserDomain{
		ID:       1,
		Name:     "New Name",
		Username: "newusername",
		Email:    "new@email.com",
	}

	expectedUser := inputUser

	mockUserRepo.
		EXPECT().
		UpdateUser(ctx, inputUser.ID, gomock.Any()).
		Return(expectedUser, nil)

	result, err := userService.UpdateUser(ctx, inputUser)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Name, result.Name)
	assert.Equal(t, expectedUser.Username, result.Username)
	assert.Equal(t, expectedUser.Email, result.Email)
}

func TestUpdateUser_WhenNoFieldsProvided_ShouldReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockdb.NewMockIUserRepository(ctrl)
	mockTokenService := mocktoken.NewMockITokenService(ctrl)
	mockPasswordService := mocksecurity.NewMockIPasswordService(ctrl)
	logger := zaptest.NewLogger(t).Sugar()

	userService := NewUserService(mockUserRepo, mockTokenService, mockPasswordService, logger)

	ctx := domain.ContextControl{}
	inputUser := domain.UserDomain{ID: 1}

	_, err := userService.UpdateUser(ctx, inputUser)

	assert.Error(t, err)
	assert.Equal(t, msg.ErrorNoFieldsToUpdate, err.Error())
}

func TestUpdateUser_WhenRepositoryFails_ShouldReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockdb.NewMockIUserRepository(ctrl)
	mockTokenService := mocktoken.NewMockITokenService(ctrl)
	mockPasswordService := mocksecurity.NewMockIPasswordService(ctrl)
	logger := zaptest.NewLogger(t).Sugar()

	userService := NewUserService(mockUserRepo, mockTokenService, mockPasswordService, logger)

	ctx := domain.ContextControl{}
	inputUser := domain.UserDomain{
		ID:    1,
		Email: "novo@email.com",
	}

	expectedErr := errors.New("db error")

	mockUserRepo.
		EXPECT().
		UpdateUser(ctx, inputUser.ID, gomock.Any()).
		Return(domain.UserDomain{}, expectedErr)

	_, err := userService.UpdateUser(ctx, inputUser)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

// Tests to UpdateUserPassword

func TestUpdateUserPassword(t *testing.T) {

	t.Run("should update password successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mockdb.NewMockIUserRepository(ctrl)
		mockTokenService := mocktoken.NewMockITokenService(ctrl)
		mockPasswordService := mocksecurity.NewMockIPasswordService(ctrl)
		logger := zaptest.NewLogger(t).Sugar()

		service := NewUserService(mockUserRepo, mockTokenService, mockPasswordService, logger)

		ctx := domain.ContextControl{}
		oldPassword := "old-password"
		newPassword := "new-password"
		hashedOld := "hashed-old"
		hashedNew := "hashed-new"

		user := domain.UserDomain{ID: 1}

		userFromDB := domain.UserDomain{
			ID:       1,
			Password: hashedOld,
		}

		mockUserRepo.EXPECT().
			GetUserByID(ctx, user.ID).
			Return(userFromDB, nil)

		mockPasswordService.EXPECT().
			ComparePasswords(hashedOld, oldPassword).
			Return(nil)

		mockPasswordService.EXPECT().
			HashPassword(newPassword).
			Return(hashedNew, nil)

		mockUserRepo.EXPECT().
			UpdateUser(ctx, user.ID, gomock.Any()).
			Return(userFromDB, nil)

		mockTokenService.EXPECT().
			Create(ctx, gomock.Any()).
			Return("new-token", nil)

		mockTokenService.EXPECT().
			Save(ctx, gomock.Any()).
			Return(nil)

		updatedUser, token, err := service.UpdateUserPassword(ctx, user, oldPassword, newPassword)

		assert.NoError(t, err)
		assert.Equal(t, user.ID, updatedUser.ID)
		assert.Equal(t, "new-token", token)
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mockdb.NewMockIUserRepository(ctrl)
		mockTokenService := mocktoken.NewMockITokenService(ctrl)
		mockPasswordService := mocksecurity.NewMockIPasswordService(ctrl)
		logger := zaptest.NewLogger(t).Sugar()

		service := NewUserService(mockUserRepo, mockTokenService, mockPasswordService, logger)

		ctx := domain.ContextControl{}
		user := domain.UserDomain{ID: 99}

		mockUserRepo.EXPECT().
			GetUserByID(ctx, user.ID).
			Return(domain.UserDomain{}, errors.New("not found"))

		_, _, err := service.UpdateUserPassword(ctx, user, "old", "new")

		assert.Error(t, err)
		assert.Equal(t, "not found", err.Error())
	})

	t.Run("should return error when password does not match", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mockdb.NewMockIUserRepository(ctrl)
		mockTokenService := mocktoken.NewMockITokenService(ctrl)
		mockPasswordService := mocksecurity.NewMockIPasswordService(ctrl)
		logger := zaptest.NewLogger(t).Sugar()

		service := NewUserService(mockUserRepo, mockTokenService, mockPasswordService, logger)

		ctx := domain.ContextControl{}
		user := domain.UserDomain{ID: 1}
		hashed := "hashed-password"

		mockUserRepo.EXPECT().
			GetUserByID(ctx, user.ID).
			Return(domain.UserDomain{ID: 1, Password: hashed}, nil)

		mockPasswordService.EXPECT().
			ComparePasswords(hashed, "wrong-password").
			Return(errors.New("invalid password"))

		_, _, err := service.UpdateUserPassword(ctx, user, "wrong-password", "new")

		assert.Error(t, err)
		assert.Equal(t, "invalid password", err.Error())
	})

	t.Run("should return error when hashing new password fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mockdb.NewMockIUserRepository(ctrl)
		mockTokenService := mocktoken.NewMockITokenService(ctrl)
		mockPasswordService := mocksecurity.NewMockIPasswordService(ctrl)
		logger := zaptest.NewLogger(t).Sugar()

		service := NewUserService(mockUserRepo, mockTokenService, mockPasswordService, logger)

		ctx := domain.ContextControl{}
		user := domain.UserDomain{ID: 1}
		hashed := "hashed-password"

		mockUserRepo.EXPECT().
			GetUserByID(ctx, user.ID).
			Return(domain.UserDomain{ID: 1, Password: hashed}, nil)

		mockPasswordService.EXPECT().
			ComparePasswords(hashed, "old").
			Return(nil)

		mockPasswordService.EXPECT().
			HashPassword("new").
			Return("", errors.New("hash error"))

		_, _, err := service.UpdateUserPassword(ctx, user, "old", "new")

		assert.Error(t, err)
		assert.Equal(t, "hash error", err.Error())
	})

	t.Run("should return error when saving new token fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mockdb.NewMockIUserRepository(ctrl)
		mockTokenService := mocktoken.NewMockITokenService(ctrl)
		mockPasswordService := mocksecurity.NewMockIPasswordService(ctrl)
		logger := zaptest.NewLogger(t).Sugar()

		service := NewUserService(mockUserRepo, mockTokenService, mockPasswordService, logger)

		ctx := domain.ContextControl{}
		user := domain.UserDomain{ID: 1}
		hashed := "hashed-password"
		hashedNew := "new-hash"

		mockUserRepo.EXPECT().
			GetUserByID(ctx, user.ID).
			Return(domain.UserDomain{ID: 1, Password: hashed}, nil)

		mockPasswordService.EXPECT().
			ComparePasswords(hashed, "old").
			Return(nil)

		mockPasswordService.EXPECT().
			HashPassword("new").
			Return(hashedNew, nil)

		mockUserRepo.EXPECT().
			UpdateUser(ctx, user.ID, gomock.Any()).
			Return(domain.UserDomain{ID: 1}, nil)

		mockTokenService.EXPECT().
			Create(ctx, domain.TokenDomain{UserID: user.ID}).
			Return("token", nil)

		mockTokenService.EXPECT().
			Save(ctx, domain.TokenDomain{
				UserID: user.ID,
				Token:  "token",
			}).
			Return(errors.New("redis down"))

		_, _, err := service.UpdateUserPassword(ctx, user, "old", "new")

		assert.Error(t, err)
		assert.Equal(t, "error saving token", err.Error())
	})
}

// Tests to SoftDeleteUser

func TestSoftDeleteUser_WhenSuccess_ShouldDeleteUserAndToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockdb.NewMockIUserRepository(ctrl)
	mockTokenService := mocktoken.NewMockITokenService(ctrl)
	mockPasswordService := mocksecurity.NewMockIPasswordService(ctrl)
	logger := zaptest.NewLogger(t).Sugar()

	service := NewUserService(mockUserRepo, mockTokenService, mockPasswordService, logger)

	ctx := domain.ContextControl{}
	userID := uint64(42)

	mockUserRepo.
		EXPECT().
		SoftDeleteUser(ctx, userID).
		Return(nil)

	mockTokenService.
		EXPECT().
		Delete(ctx, domain.TokenDomain{UserID: userID}).
		Return(nil)

	err := service.SoftDeleteUser(ctx, userID)

	assert.NoError(t, err)
}

func TestSoftDeleteUser_WhenUserRepoFails_ShouldReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockdb.NewMockIUserRepository(ctrl)
	mockTokenService := mocktoken.NewMockITokenService(ctrl)
	mockPasswordService := mocksecurity.NewMockIPasswordService(ctrl)
	logger := zaptest.NewLogger(t).Sugar()

	service := NewUserService(mockUserRepo, mockTokenService, mockPasswordService, logger)

	ctx := domain.ContextControl{}
	userID := uint64(42)
	expectedErr := errors.New("repo failure")

	mockUserRepo.
		EXPECT().
		SoftDeleteUser(ctx, userID).
		Return(expectedErr)

	err := service.SoftDeleteUser(ctx, userID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestSoftDeleteUser_WhenTokenServiceFails_ShouldReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockdb.NewMockIUserRepository(ctrl)
	mockTokenService := mocktoken.NewMockITokenService(ctrl)
	mockPasswordService := mocksecurity.NewMockIPasswordService(ctrl)
	logger := zaptest.NewLogger(t).Sugar()

	service := NewUserService(mockUserRepo, mockTokenService, mockPasswordService, logger)

	ctx := domain.ContextControl{}
	userID := uint64(42)
	expectedErr := errors.New("token deletion error")

	mockUserRepo.
		EXPECT().
		SoftDeleteUser(ctx, userID).
		Return(nil)

	mockTokenService.
		EXPECT().
		Delete(ctx, domain.TokenDomain{UserID: userID}).
		Return(expectedErr)

	err := service.SoftDeleteUser(ctx, userID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

// Tests NormalizeUserData

func TestNormalizeUserData_ShouldTrimSpacesAndLowercaseEmail(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	svc := &UserService{LoggerSugar: logger}

	raw := &domain.UserDomain{
		Name:     "  Felipe  ",
		Username: " lechitz ",
		Email:    "  LECHITZ@AION.COM  ",
	}

	normalized := svc.NormalizeUserData(raw)

	assert.Equal(t, "Felipe", normalized.Name)
	assert.Equal(t, "lechitz", normalized.Username)
	assert.Equal(t, "lechitz@aion.com", normalized.Email)
}

// Tests validateCreateUserRequired

func TestValidateCreateUserRequired_WhenAllFieldsValid_ShouldReturnNil(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	svc := &UserService{LoggerSugar: logger}

	user := domain.UserDomain{
		Name:     "Felipe",
		Username: "lechitz",
		Email:    "lechitz@aion.com",
	}
	err := svc.validateCreateUserRequired(user, "123456")
	assert.NoError(t, err)
}

func TestValidateCreateUserRequired_WhenMissingName_ShouldReturnError(t *testing.T) {
	svc := &UserService{}
	user := domain.UserDomain{
		Username: "lechitz",
		Email:    "lechitz@aion.com",
	}
	err := svc.validateCreateUserRequired(user, "123456")
	assert.EqualError(t, err, msg.NameIsRequired)
}

func TestValidateCreateUserRequired_WhenInvalidEmail_ShouldReturnError(t *testing.T) {
	svc := &UserService{}
	user := domain.UserDomain{
		Name:     "Felipe",
		Username: "lechitz",
		Email:    "invalid-email",
	}
	err := svc.validateCreateUserRequired(user, "123456")
	assert.EqualError(t, err, msg.InvalidEmail)
}
