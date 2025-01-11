package service

import (
	"errors"
	"github.com/badoux/checkmail"
	"github.com/lechitz/AionApi/internal/core/constants"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/pkg/utils"
	"github.com/lechitz/AionApi/ports/output"
	"go.uber.org/zap"
	"strings"
)

type UserService struct {
	UserDomainDataBaseRepository output.IUserDomainDataBaseRepository
	UserDomainCacheRepository    output.IAuthRepository
	LoggerSugar                  *zap.SugaredLogger
}

func (service *UserService) CreateUser(contextControl domain.ContextControl, user domain.UserDomain) (domain.UserDomain, error) {

	if err := service.validateCreateUser(user); err != nil {
		service.LoggerSugar.Errorw(constants.ErrorToValidateCreateUser, "error", err.Error())
		return domain.UserDomain{}, err
	}

	if err := service.formatCreateUser(&user); err != nil {
		service.LoggerSugar.Errorw(constants.ErrorToFormatCreateUser, "error", err.Error())
		return domain.UserDomain{}, err
	}

	user, err := service.UserDomainDataBaseRepository.CreateUser(contextControl, user)
	if err != nil {
		service.LoggerSugar.Errorw(constants.ErrorToCreateUser, "error", err.Error())
		return domain.UserDomain{}, err
	}

	service.LoggerSugar.Infow(constants.SuccessUserCreated, "userID", user.ID)
	return user, nil
}

func (service *UserService) GetAllUsers(contextControl domain.ContextControl) ([]domain.UserDomain, error) {

	users, err := service.UserDomainDataBaseRepository.GetAllUsers(contextControl)
	if err != nil {
		service.LoggerSugar.Errorw(constants.ErrorToGetAllUsers, "error", err.Error())
		return []domain.UserDomain{}, err
	}

	service.LoggerSugar.Infow(constants.SuccessUsersRetrieved, "users", users)
	return users, nil
}

func (service *UserService) GetUserByID(contextControl domain.ContextControl, ID uint64) (domain.UserDomain, error) {

	user, err := service.UserDomainDataBaseRepository.GetUserByID(contextControl, ID)
	if err != nil {
		service.LoggerSugar.Errorw(constants.ErrorToGetUserByID, "error", err.Error())
		return domain.UserDomain{}, err
	}

	service.LoggerSugar.Infow(constants.SuccessUserRetrieved, "userID", user.ID)
	return user, nil
}

func (service *UserService) GetUserByUsername(contextControl domain.ContextControl, userDomain domain.UserDomain) (domain.UserDomain, error) {

	user, err := service.UserDomainDataBaseRepository.GetUserByUsername(contextControl, userDomain)
	if err != nil {
		service.LoggerSugar.Errorw(constants.ErrorToGetUserByUserName, "error", err.Error())
		return domain.UserDomain{}, err
	}

	service.LoggerSugar.Infow(constants.SuccessUserRetrieved, "userID", user.ID)
	return user, nil
}

func (service *UserService) UpdateUser(contextControl domain.ContextControl, user domain.UserDomain) (domain.UserDomain, error) {

	if err := service.formatUpdateUser(user); err != nil {
		service.LoggerSugar.Errorw(constants.ErrorToFormatUpdateUser, "error", err.Error())
		return domain.UserDomain{}, err
	}

	user, err := service.UserDomainDataBaseRepository.UpdateUser(contextControl, user)
	if err != nil {
		service.LoggerSugar.Errorw(constants.ErrorToUpdateUser, "error", err.Error())
		return domain.UserDomain{}, err
	}

	user, err = service.UserDomainDataBaseRepository.GetUserByID(contextControl, user.ID)
	if err != nil {
		service.LoggerSugar.Errorw(constants.ErrorToGetUserByID, "error", err.Error())
		return domain.UserDomain{}, err
	}

	service.LoggerSugar.Infow(constants.SuccessUserUpdated, "userID", user.ID)
	return user, nil
}

func (service *UserService) UpdatePassword(contextControl domain.ContextControl, userDomain domain.UserDomain, password, newPassword string) (domain.UserDomain, string, error) {

	userDB, err := service.UserDomainDataBaseRepository.GetUserByID(contextControl, userDomain.ID)
	if err != nil {
		service.LoggerSugar.Errorw(constants.ErrorToGetUserByID, "error", err.Error())
		return domain.UserDomain{}, "", err
	}

	if err := utils.VerifyPassword(userDB.Password, password); err != nil {
		service.LoggerSugar.Errorw(constants.ErrorToVerifyPassword, "error", err.Error())
		return domain.UserDomain{}, "", err
	}

	err = service.UserDomainCacheRepository.DeleteTokenByUserID(contextControl.Context, userDB.ID)
	if err != nil {
		service.LoggerSugar.Errorw(constants.ErrorToDeleteToken, "error", err.Error())
		return domain.UserDomain{}, "", err
	}

	hashedPassword, err := utils.Hash(newPassword)
	if err != nil {
		service.LoggerSugar.Errorw(constants.ErrorToHashPassword, "error", err.Error())
		return domain.UserDomain{}, "", err
	}

	userDB.Password = string(hashedPassword)

	userDB, err = service.UserDomainDataBaseRepository.UpdatePassword(contextControl, userDB)
	if err != nil {
		service.LoggerSugar.Errorw(constants.ErrorToUpdatePassword, "error", err.Error())
		return domain.UserDomain{}, "", err
	}

	newToken, err := service.UserDomainCacheRepository.CreateToken(contextControl.Context, userDB.ID)

	service.LoggerSugar.Infow(constants.SuccessPasswordUpdated, "userID", userDB.ID)
	return userDB, newToken, nil
}

func (service *UserService) SoftDeleteUser(contextControl domain.ContextControl, ID uint64) error {

	err := service.UserDomainDataBaseRepository.SoftDeleteUser(contextControl, ID)
	if err != nil {
		service.LoggerSugar.Errorw(constants.ErrorToSoftDeleteUser, "error", err.Error())
		return err
	}

	service.LoggerSugar.Infow(constants.SuccessUserSoftDeleted, "userID", ID)
	return nil
}

func (service *UserService) validateCreateUser(user domain.UserDomain) error {
	if user.Name == "" {
		return errors.New(constants.NameIsRequired)
	}
	if user.Username == "" {
		return errors.New(constants.UsernameIsRequired)
	}
	if user.Email == "" {
		return errors.New(constants.EmailIsRequired)
	}

	if user.Password == "" {
		return errors.New(constants.PasswordIsRequired)
	}

	if err := checkmail.ValidateFormat(user.Email); err != nil {
		return errors.New(constants.InvalidEmail)
	}

	return nil
}

func (service *UserService) formatCreateUser(user *domain.UserDomain) error {
	user.Name = strings.TrimSpace(user.Name)
	user.Username = strings.TrimSpace(user.Username)
	user.Email = strings.TrimSpace(user.Email)

	hashedPassword, err := utils.Hash(user.Password)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	return nil
}

func (service *UserService) formatUpdateUser(user domain.UserDomain) error {
	if user.Name != "" {
		user.Name = strings.TrimSpace(user.Name)
	}

	if user.Username != "" {
		user.Username = strings.TrimSpace(user.Username)
	}

	if user.Email != "" {
		user.Email = strings.TrimSpace(user.Email)
	}

	return nil
}
