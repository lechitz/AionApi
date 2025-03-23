package service

import (
	"errors"
	"github.com/lechitz/AionApi/core/domain"
	"github.com/lechitz/AionApi/ports/output/db"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/lechitz/AionApi/core/msg"
	"github.com/lechitz/AionApi/pkg/contextkeys"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepository db.IUserRepository
	TokenService   *TokenService
	AuthService    *AuthService
	LoggerSugar    *zap.SugaredLogger
}

func NewUserService(userRepo db.IUserRepository, tokenService *TokenService, authService *AuthService, loggerSugar *zap.SugaredLogger) *UserService {
	return &UserService{
		UserRepository: userRepo,
		TokenService:   tokenService,
		AuthService:    authService,
		LoggerSugar:    loggerSugar,
	}
}

func (service *UserService) CreateUser(ctx domain.ContextControl, userDomain domain.UserDomain, passwordReq string) (domain.UserDomain, error) {

	normalizedUser := service.NormalizeUserData(&userDomain)

	if err := service.validateCreateUserRequired(normalizedUser, passwordReq); err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToValidateCreateUser, contextkeys.Error, err.Error())
		return domain.UserDomain{}, err
	}

	hashPassword, err := service.HashPassword(passwordReq)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToHashPassword, contextkeys.Error, err.Error())
		return domain.UserDomain{}, err
	}

	normalizedUser.Password = hashPassword

	userDB, err := service.UserRepository.CreateUser(ctx, normalizedUser)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToCreateUser, contextkeys.Error, err.Error())
		return domain.UserDomain{}, err
	}

	service.LoggerSugar.Infow(msg.SuccessUserCreated, contextkeys.UserID, userDB.ID)
	return userDB, nil
}

func (service *UserService) GetAllUsers(ctx domain.ContextControl) ([]domain.UserDomain, error) {
	users, err := service.UserRepository.GetAllUsers(ctx)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToGetAllUsers, contextkeys.Error, err.Error())
		return []domain.UserDomain{}, err
	}

	service.LoggerSugar.Infow(msg.SuccessUsersRetrieved, "count", len(users))
	return users, nil
}

func (service *UserService) GetUserByID(ctx domain.ContextControl, ID uint64) (domain.UserDomain, error) {
	user, err := service.UserRepository.GetUserByID(ctx, ID)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToGetUserByID, contextkeys.Error, err.Error())
		return domain.UserDomain{}, err
	}

	service.LoggerSugar.Infow(msg.SuccessUserRetrieved, contextkeys.UserID, user.ID)
	return user, nil
}

func (service *UserService) GetUserByUsername(ctx domain.ContextControl, userDomain domain.UserDomain) (domain.UserDomain, error) {
	user, err := service.UserRepository.GetUserByUsername(ctx, userDomain.Username)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToGetUserByUserName, contextkeys.Error, err.Error())
		return domain.UserDomain{}, err
	}

	service.LoggerSugar.Infow(msg.SuccessUserRetrieved, contextkeys.UserID, user.ID)
	return user, nil
}

func (service *UserService) UpdateUser(ctx domain.ContextControl, user domain.UserDomain) (domain.UserDomain, error) {
	updateFields := map[string]interface{}{}

	if user.Name != "" {
		updateFields["name"] = user.Name
	}
	if user.Username != "" {
		updateFields["username"] = user.Username
	}
	if user.Email != "" {
		updateFields["email"] = user.Email
	}

	if len(updateFields) == 0 {
		return domain.UserDomain{}, errors.New("no fields to update")
	}

	updateFields["updated_at"] = time.Now()

	updatedUser, err := service.UserRepository.UpdateUserFields(ctx, user.ID, updateFields)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToUpdateUser, contextkeys.Error, err.Error())
		return domain.UserDomain{}, err
	}

	return updatedUser, nil
}

func (service *UserService) SoftDeleteUser(ctx domain.ContextControl, userID uint64) error {
	fields := map[string]interface{}{
		"deleted_at": time.Now(),
	}

	if _, err := service.UserRepository.UpdateUserFields(ctx, userID, fields); err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToSoftDeleteUser, contextkeys.Error, err.Error())
		return err
	}

	tokenDomain := domain.TokenDomain{
		UserID: userID,
	}

	if err := service.TokenService.DeleteToken(ctx, tokenDomain); err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToDeleteToken, contextkeys.Error, err.Error())
		return err
	}

	service.LoggerSugar.Infow(msg.SuccessUserSoftDeleted, contextkeys.UserID, userID)
	return nil
}

func (service *UserService) UpdateUserPassword(ctx domain.ContextControl, userDomain domain.UserDomain, passwordReq, newPasswordReq string) (domain.UserDomain, string, error) {

	userDB, err := service.UserRepository.GetUserByID(ctx, userDomain.ID)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToGetUserByID, contextkeys.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	if err = service.AuthService.compareHashAndPassword(userDB.Password, passwordReq); err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToCompareHashAndPassword, contextkeys.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	hashedPassword, err := service.HashPassword(newPasswordReq)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToHashPassword, contextkeys.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	updateFields := map[string]interface{}{
		contextkeys.Password:  hashedPassword,
		contextkeys.UpdatedAt: time.Now(),
	}

	updatedUser, err := service.UserRepository.UpdateUserFields(ctx, userDomain.ID, updateFields)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToUpdatePassword, contextkeys.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	tokenDomain := domain.TokenDomain{
		UserID: userDomain.ID,
	}
	token, err := service.TokenService.CreateToken(ctx, tokenDomain)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToCreateToken, contextkeys.Error, err.Error())
		return domain.UserDomain{}, "", err
	}
	tokenDomain.Token = token

	if err := service.TokenService.SaveToken(ctx, tokenDomain); err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToSaveToken, contextkeys.Error, err.Error())
		return domain.UserDomain{}, "", errors.New("error to save token")
	}

	service.LoggerSugar.Infow(msg.SuccessPasswordUpdated, contextkeys.UserID, updatedUser.ID)
	return updatedUser, tokenDomain.Token, nil
}

func (service *UserService) NormalizeUserData(userDomain *domain.UserDomain) domain.UserDomain {
	if userDomain.Name != "" {
		userDomain.Name = strings.TrimSpace(userDomain.Name)
	}
	if userDomain.Username != "" {
		userDomain.Username = strings.TrimSpace(userDomain.Username)
	}
	if userDomain.Email != "" {
		userDomain.Email = strings.ToLower(strings.TrimSpace(userDomain.Email))
	}

	return *userDomain
}

func (service *UserService) validateCreateUserRequired(userDomain domain.UserDomain, password string) error {
	if userDomain.Name == "" {
		return errors.New(msg.NameIsRequired)
	}
	if userDomain.Username == "" {
		return errors.New(msg.UsernameIsRequired)
	}
	if userDomain.Email == "" {
		return errors.New(msg.EmailIsRequired)
	}
	if password == "" {
		return errors.New(msg.PasswordIsRequired)
	}
	if err := checkmail.ValidateFormat(userDomain.Email); err != nil {
		return errors.New(msg.InvalidEmail)
	}
	return nil
}

func (service *UserService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
