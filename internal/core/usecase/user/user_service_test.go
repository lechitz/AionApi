package user

//
//import (
//	"errors"
//	"github.com/golang/mock/gomock"
//	"github.com/lechitz/AionApi/internal/core/domain"
//	"github.com/lechitz/AionApi/internal/core/usecase/constants"
//	"github.com/lechitz/AionApi/tests/mocks"
//	"github.com/stretchr/testify/assert"
//	"go.uber.org/zap"
//	"go.uber.org/zap/zaptest"
//	"gorm.io/gorm"
//	"testing"
//)
//
//type TestDependencies struct {
//	Controller       *gomock.Controller
//	Logger           *zap.SugaredLogger
//	MockUserRepo     *mocks.MockUserRepository
//	MockTokenService *mocks.MockTokenServiceInterface
//	MockPasswordSvc  *mocks.MockPasswordManager
//	Service          *UserService
//}
//
//func NewUserServiceTestSetup(t *testing.T) *TestDependencies {
//	ctrl := gomock.NewController(t)
//	logger := zaptest.NewLogger(t).Sugar()
//
//	mockUserRepo := mocks.NewMockUserRepository(ctrl)
//	mockTokenSvc := mocks.NewMockTokenServiceInterface(ctrl)
//	mockPasswordSvc := mocks.NewMockPasswordManager(ctrl)
//
//	userService := NewUserService(
//		mockUserRepo,
//		mockTokenSvc,
//		mockPasswordSvc,
//		logger,
//	)
//
//	return &TestDependencies{
//		Controller:       ctrl,
//		Logger:           logger,
//		MockUserRepo:     mockUserRepo,
//		MockTokenService: mockTokenSvc,
//		MockPasswordSvc:  mockPasswordSvc,
//		Service:          userService,
//	}
//}
//
//// CreateUser tests
//
//func TestCreateUser_Success(t *testing.T) {
//	td := NewUserServiceTestSetup(t)
//	defer td.Controller.Finish()
//
//	ctx := domain.ContextControl{}
//	user := domain.UserDomain{
//		Name:     "Test User",
//		Username: "testuser",
//		Email:    "testuser@example.com",
//	}
//	password := "securepassword"
//
//	td.MockUserRepo.EXPECT().
//		GetUserByUsername(ctx, user.Username).
//		Return(domain.UserDomain{}, gorm.ErrRecordNotFound)
//	td.MockUserRepo.EXPECT().
//		GetUserByEmail(ctx, user.Email).
//		Return(domain.UserDomain{}, gorm.ErrRecordNotFound)
//	td.MockPasswordSvc.EXPECT().
//		HashPassword(password).
//		Return("hashedPassword", nil)
//	td.MockUserRepo.EXPECT().
//		CreateUser(ctx, gomock.Any()).
//		DoAndReturn(func(ctx domain.ContextControl, u domain.UserDomain) (domain.UserDomain, error) {
//			u.ID = 1
//			u.Password = "hashedPassword"
//			return u, nil
//		})
//
//	createdUser, err := td.Service.CreateUser(ctx, user, password)
//	assert.NoError(t, err)
//	assert.Equal(t, uint64(1), createdUser.ID)
//	assert.Equal(t, "hashedPassword", createdUser.Password)
//}
//
//func TestCreateUser_UsernameAlreadyExists(t *testing.T) {
//	td := NewUserServiceTestSetup(t)
//	defer td.Controller.Finish()
//
//	ctx := domain.ContextControl{}
//	user := domain.UserDomain{
//		Name:     "Test User",
//		Username: "existinguser",
//		Email:    "newemail@example.com",
//	}
//	password := "securepassword"
//
//	td.MockUserRepo.EXPECT().
//		GetUserByUsername(ctx, user.Username).
//		Return(domain.UserDomain{ID: 1}, nil)
//
//	_, err := td.Service.CreateUser(ctx, user, password)
//	assert.Error(t, err)
//	assert.Equal(t, constants.UsernameIsAlreadyInUse, err.Error())
//}
//
//func TestCreateUser_EmailAlreadyExists(t *testing.T) {
//	td := NewUserServiceTestSetup(t)
//	defer td.Controller.Finish()
//
//	ctx := domain.ContextControl{}
//	user := domain.UserDomain{
//		Name:     "Test User",
//		Username: "newuser",
//		Email:    "existingemail@example.com",
//	}
//	password := "securepassword"
//
//	td.MockUserRepo.EXPECT().
//		GetUserByUsername(ctx, user.Username).
//		Return(domain.UserDomain{}, gorm.ErrRecordNotFound)
//	td.MockUserRepo.EXPECT().
//		GetUserByEmail(ctx, user.Email).
//		Return(domain.UserDomain{ID: 1}, nil)
//
//	_, err := td.Service.CreateUser(ctx, user, password)
//	assert.Error(t, err)
//	assert.Equal(t, constants.EmailIsAlreadyInUse, err.Error())
//}
//
//func TestCreateUser_PasswordHashingFails(t *testing.T) {
//	td := NewUserServiceTestSetup(t)
//	defer td.Controller.Finish()
//
//	ctx := domain.ContextControl{}
//	user := domain.UserDomain{
//		Name:     "Test User",
//		Username: "testuser",
//		Email:    "testuser@example.com",
//	}
//	password := "securepassword"
//
//	td.MockUserRepo.EXPECT().
//		GetUserByUsername(ctx, user.Username).
//		Return(domain.UserDomain{}, gorm.ErrRecordNotFound)
//	td.MockUserRepo.EXPECT().
//		GetUserByEmail(ctx, user.Email).
//		Return(domain.UserDomain{}, gorm.ErrRecordNotFound)
//	td.MockPasswordSvc.EXPECT().
//		HashPassword(password).
//		Return("", errors.New("hashing error"))
//
//	_, err := td.Service.CreateUser(ctx, user, password)
//	assert.EqualError(t, err, constants.ErrorToHashPassword)
//}
//
//func TestCreateUser_CreateUserFails(t *testing.T) {
//	td := NewUserServiceTestSetup(t)
//	defer td.Controller.Finish()
//
//	ctx := domain.ContextControl{}
//	user := domain.UserDomain{
//		Name:     "Test User",
//		Username: "testuser",
//		Email:    "testuser@example.com",
//	}
//	password := "securepassword"
//
//	td.MockUserRepo.EXPECT().
//		GetUserByUsername(ctx, user.Username).
//		Return(domain.UserDomain{}, gorm.ErrRecordNotFound)
//	td.MockUserRepo.EXPECT().
//		GetUserByEmail(ctx, user.Email).
//		Return(domain.UserDomain{}, gorm.ErrRecordNotFound)
//	td.MockPasswordSvc.EXPECT().
//		HashPassword(password).
//		Return("hashedPassword", nil)
//	td.MockUserRepo.EXPECT().
//		CreateUser(ctx, gomock.Any()).
//		Return(domain.UserDomain{}, errors.New(constants.ErrorToCreateUser))
//
//	_, err := td.Service.CreateUser(ctx, user, password)
//	assert.Error(t, err)
//	assert.Equal(t, constants.ErrorToCreateUser, err.Error())
//}
//
//// GetAllUsers tests
//
//// GetUserByID tests
//
//// GetUserByUsername tests
//
//// GetUserByEmail tests
//
//// UpdateUser tests
//
//// UpdateUserPassword tests
//
//// SoftDeleteUser tests
//
//func TestSoftDeleteUser_Success(t *testing.T) {
//	td := NewUserServiceTestSetup(t)
//	defer td.Controller.Finish()
//
//	ctx := domain.ContextControl{}
//	userID := uint64(42)
//
//	td.MockUserRepo.EXPECT().
//		SoftDeleteUser(ctx, userID).
//		Return(nil)
//
//	td.MockTokenService.EXPECT().
//		Delete(ctx, domain.TokenDomain{UserID: userID}).
//		Return(nil)
//
//	err := td.Service.SoftDeleteUser(ctx, userID)
//	assert.NoError(t, err)
//}
//
//func TestSoftDeleteUser_FailureOnUserRepo(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockUserRepo := mocks.NewMockUserRepository(ctrl)
//	mockTokenService := mocks.NewMockTokenServiceInterface(ctrl)
//	mockPasswordService := mocks.NewMockPasswordManager(ctrl)
//	logger := zaptest.NewLogger(t).Sugar()
//
//	service := NewUserService(mockUserRepo, mockTokenService, mockPasswordService, logger)
//
//	ctx := domain.ContextControl{}
//	userID := uint64(99)
//	expectedErr := errors.New(constants.ErrorToSoftDeleteUser)
//
//	mockUserRepo.EXPECT().
//		SoftDeleteUser(ctx, userID).
//		Return(expectedErr)
//
//	err := service.SoftDeleteUser(ctx, userID)
//	assert.EqualError(t, err, expectedErr.Error())
//}
//
//func TestSoftDeleteUser_FailureOnTokenDelete(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockUserRepo := mocks.NewMockUserRepository(ctrl)
//	mockTokenService := mocks.NewMockTokenServiceInterface(ctrl)
//	mockPasswordService := mocks.NewMockPasswordManager(ctrl)
//	logger := zaptest.NewLogger(t).Sugar()
//
//	service := NewUserService(mockUserRepo, mockTokenService, mockPasswordService, logger)
//
//	ctx := domain.ContextControl{}
//	userID := uint64(77)
//	expectedErr := errors.New(constants.ErrorToSoftDeleteUser)
//
//	mockUserRepo.EXPECT().
//		SoftDeleteUser(ctx, userID).
//		Return(nil)
//
//	mockTokenService.EXPECT().
//		Delete(ctx, domain.TokenDomain{UserID: userID}).
//		Return(expectedErr)
//
//	err := service.SoftDeleteUser(ctx, userID)
//	assert.EqualError(t, err, expectedErr.Error())
//}
//
//// TestNormalizeUserData tests
//func TestNormalizeUserData(t *testing.T) {
//	logger := zaptest.NewLogger(t).Sugar()
//	service := &UserService{LoggerSugar: logger}
//
//	raw := &domain.UserDomain{
//		Name:     "  Felipe  ",
//		Username: " lechitz ",
//		Email:    "  LECHITZ@AION.COM ",
//	}
//
//	normalized := service.NormalizeUserData(raw)
//
//	assert.Equal(t, "Felipe", normalized.Name)
//	assert.Equal(t, "lechitz", normalized.Username)
//	assert.Equal(t, "lechitz@aion.com", normalized.Email)
//}
//
//// TestValidateCreateUserRequired tests
//func TestValidateCreateUserRequired(t *testing.T) {
//	service := &UserService{}
//
//	tests := []struct {
//		name     string
//		user     domain.UserDomain
//		password string
//		wantErr  string
//	}{
//		{"missing name", domain.UserDomain{Username: "a", Email: "a@b.com"}, "123", constants.NameIsRequired},
//		{"missing username", domain.UserDomain{Name: "a", Email: "a@b.com"}, "123", constants.UsernameIsRequired},
//		{"missing email", domain.UserDomain{Name: "a", Username: "a"}, "123", constants.EmailIsRequired},
//		{"invalid email", domain.UserDomain{Name: "a", Username: "a", Email: "invalid"}, "123", constants.InvalidEmail},
//		{"missing password", domain.UserDomain{Name: "a", Username: "a", Email: "a@b.com"}, "", constants.PasswordIsRequired},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			err := service.validateCreateUserRequired(tt.user, tt.password)
//			assert.EqualError(t, err, tt.wantErr)
//		})
//	}
//}
