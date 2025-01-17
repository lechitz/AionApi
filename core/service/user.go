package service

import (
	"errors"
	"strings"

	"github.com/badoux/checkmail"
	"github.com/lechitz/AionApi/core/domain/entities"
	msg "github.com/lechitz/AionApi/core/msg"
	"github.com/lechitz/AionApi/pkg/contextkeys"
	"github.com/lechitz/AionApi/ports/output/cache"
	"github.com/lechitz/AionApi/ports/output/db"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepository  db.IUserRepository
	TokenRepository cache.ITokenRepository
	AuthService     *AuthService
	LoggerSugar     *zap.SugaredLogger
}

func NewUserService(userRepo db.IUserRepository, tokenRepo cache.ITokenRepository, loggerSugar *zap.SugaredLogger, authService *AuthService) *UserService {
	return &UserService{
		UserRepository:  userRepo,
		TokenRepository: tokenRepo,
		LoggerSugar:     loggerSugar,
		AuthService:     authService,
	}
}

func (service *UserService) CreateUser(ctx entities.ContextControl, userDomain entities.UserDomain, passwordReq string) (entities.UserDomain, error) {

	if err := service.validateCreateUserRequired(userDomain, passwordReq); err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToValidateCreateUser, contextkeys.Error, err.Error())
		return entities.UserDomain{}, err
	}

	if err := service.NormalizeUserData(&userDomain); err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToNormalizeUserData, contextkeys.Error, err.Error())
		return entities.UserDomain{}, err
	}

	hashPassword, err := service.HashPassword(passwordReq)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToHashPassword, contextkeys.Error, err.Error())
		return entities.UserDomain{}, err
	}

	userDomain.Password = hashPassword

	userDB, err := service.UserRepository.CreateUser(ctx, userDomain)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToCreateUser, contextkeys.Error, err.Error())
		return entities.UserDomain{}, err
	}

	service.LoggerSugar.Infow(msg.SuccessUserCreated, contextkeys.UserID, userDB.ID)
	return userDB, nil
}

func (service *UserService) GetAllUsers(ctx entities.ContextControl) ([]entities.UserDomain, error) {
	users, err := service.UserRepository.GetAllUsers(ctx)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToGetAllUsers, contextkeys.Error, err.Error())
		return []entities.UserDomain{}, err
	}

	service.LoggerSugar.Infow(msg.SuccessUsersRetrieved, "count", len(users))
	return users, nil
}

func (service *UserService) GetUserByID(ctx entities.ContextControl, ID uint64) (entities.UserDomain, error) {
	user, err := service.UserRepository.GetUserByID(ctx, ID)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToGetUserByID, contextkeys.Error, err.Error())
		return entities.UserDomain{}, err
	}

	service.LoggerSugar.Infow(msg.SuccessUserRetrieved, contextkeys.UserID, user.ID)
	return user, nil
}

func (service *UserService) GetUserByUsername(ctx entities.ContextControl, userDomain entities.UserDomain) (entities.UserDomain, error) {
	user, err := service.UserRepository.GetUserByUsername(ctx, userDomain.Username)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToGetUserByUserName, contextkeys.Error, err.Error())
		return entities.UserDomain{}, err
	}

	service.LoggerSugar.Infow(msg.SuccessUserRetrieved, contextkeys.UserID, user.ID)
	return user, nil
}

func (service *UserService) UpdateUser(ctx entities.ContextControl, userDomain entities.UserDomain) (entities.UserDomain, error) {
	if err := service.NormalizeUserData(&userDomain); err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToFormatUpdateUser, contextkeys.Error, err.Error())
		return entities.UserDomain{}, err
	}

	userDB, err := service.UserRepository.UpdateUser(ctx, userDomain)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToUpdateUser, contextkeys.Error, err.Error())
		return entities.UserDomain{}, err
	}

	userDB, err = service.UserRepository.GetUserByID(ctx, userDB.ID)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToGetUserByID, contextkeys.Error, err.Error())
		return entities.UserDomain{}, err
	}

	service.LoggerSugar.Infow(msg.SuccessUserUpdated, contextkeys.UserID, userDB.ID)
	return userDB, nil
}

func (service *UserService) SoftDeleteUser(ctx entities.ContextControl, userID uint64) error {
	if err := service.UserRepository.SoftDeleteUser(ctx, userID); err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToSoftDeleteUser, contextkeys.Error, err.Error())
		return err
	}

	tokenDomain := entities.TokenDomain{
		UserID: userID,
	}

	if err := service.TokenRepository.DeleteToken(ctx, tokenDomain); err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToDeleteToken, contextkeys.Error, err.Error())
		return err
	}

	service.LoggerSugar.Infow(msg.SuccessUserSoftDeleted, contextkeys.UserID, userID)
	return nil
}

func (service *UserService) UpdateUserPassword(ctx entities.ContextControl, userDomain entities.UserDomain, passwordReq, newPasswordReq string) (entities.UserDomain, string, error) {

	userDB, err := service.UserRepository.GetUserByID(ctx, userDomain.ID)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToGetUserByID, contextkeys.Error, err.Error())
		return entities.UserDomain{}, "", err
	}

	if err := service.AuthService.compareHashAndPassword(userDB.Password, passwordReq); err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToCompareHashAndPassword, contextkeys.Error, err.Error())
		return entities.UserDomain{}, "", err
	}

	tokenDomain := entities.TokenDomain{
		UserID: userDB.ID,
	}
	err = service.TokenRepository.DeleteToken(ctx, tokenDomain)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToDeleteToken, contextkeys.Error, err.Error())
		return entities.UserDomain{}, "", err
	}

	hashedPassword, err := service.HashPassword(newPasswordReq)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToHashPassword, contextkeys.Error, err.Error())
		return entities.UserDomain{}, "", err
	}

	userDB.Password = hashedPassword
	userDB, err = service.UserRepository.UpdatePassword(ctx, userDB)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToUpdatePassword, contextkeys.Error, err.Error())
		return entities.UserDomain{}, "", err
	}

	newToken, err := service.AuthService.TokenService.CreateToken(ctx, tokenDomain)
	// tokenDomain = {UserID: userDB.ID, Token: ""}
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToGenerateToken, contextkeys.Error, err.Error())
		return userDB, "", err
	}

	newTokenDomain := entities.TokenDomain{
		UserID: userDB.ID,
		Token:  newToken,
	}
	if err := service.TokenRepository.SaveToken(ctx, newTokenDomain); err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToSaveToken, contextkeys.Error, err.Error())
		return userDB, "", err
	}

	service.LoggerSugar.Infow(msg.SuccessPasswordUpdated, contextkeys.UserID, userDB.ID)
	return userDB, newToken, nil
}

func (service *UserService) NormalizeUserData(userDomain *entities.UserDomain) error {
	if userDomain.Name != "" {
		userDomain.Name = strings.TrimSpace(userDomain.Name)
	}
	if userDomain.Username != "" {
		userDomain.Username = strings.TrimSpace(userDomain.Username)
	}
	if userDomain.Email != "" {
		userDomain.Email = strings.ToLower(strings.TrimSpace(userDomain.Email))
	}
	return nil
}

func (service *UserService) validateCreateUserRequired(userDomain entities.UserDomain, password string) error {
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
