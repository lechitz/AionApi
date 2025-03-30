package user

import (
	"errors"
	"fmt"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/ports/output/db"
	"github.com/lechitz/AionApi/internal/core/ports/output/security"
	"github.com/lechitz/AionApi/internal/core/service/constants"
	"github.com/lechitz/AionApi/internal/core/service/token"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"go.uber.org/zap"
)

type UserService struct {
	UserRepository  db.UserRepository
	TokenService    token.TokenService
	PasswordService security.PasswordManager
	LoggerSugar     *zap.SugaredLogger
}

type UserServiceInterface interface {
	CreateUser(ctx domain.ContextControl, user domain.UserDomain, password string) (domain.UserDomain, error)
	GetAllUsers(ctx domain.ContextControl) ([]domain.UserDomain, error)
	GetUserByID(ctx domain.ContextControl, id uint64) (domain.UserDomain, error)
	GetUserByUsername(ctx domain.ContextControl, username string) (domain.UserDomain, error)
	GetUserByEmail(ctx domain.ContextControl, email string) (domain.UserDomain, error)
	UpdateUser(ctx domain.ContextControl, user domain.UserDomain) (domain.UserDomain, error)
	UpdateUserPassword(ctx domain.ContextControl, user domain.UserDomain, oldPassword, newPassword string) (domain.UserDomain, string, error)
	SoftDeleteUser(ctx domain.ContextControl, userID uint64) error
	NormalizeUserData(user *domain.UserDomain) domain.UserDomain
}

func NewUserService(userRepo db.UserRepository, tokenService token.TokenService, passwordService security.PasswordManager, loggerSugar *zap.SugaredLogger) *UserService {
	return &UserService{
		UserRepository:  userRepo,
		TokenService:    tokenService,
		PasswordService: passwordService,
		LoggerSugar:     loggerSugar,
	}
}

func (s *UserService) CreateUser(ctx domain.ContextControl, user domain.UserDomain, password string) (domain.UserDomain, error) {
	user = s.NormalizeUserData(&user)

	if err := s.validateCreateUserRequired(user, password); err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToValidateCreateUser, constants.Error, err.Error())
		return domain.UserDomain{}, err
	}

	existingByUsername, err := s.UserRepository.GetUserByUsername(ctx, user.Username)
	if err == nil && existingByUsername.ID != 0 {
		return domain.UserDomain{}, fmt.Errorf("username '%s' is already in use", user.Username)
	}

	existingByEmail, err := s.UserRepository.GetUserByEmail(ctx, user.Email)
	if err == nil && existingByEmail.ID != 0 {
		return domain.UserDomain{}, fmt.Errorf("email '%s' is already in use", user.Email)
	}

	hashPassword, err := s.PasswordService.HashPassword(password)
	if err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToHashPassword, constants.Error, err.Error())
		return domain.UserDomain{}, err
	}
	user.Password = hashPassword

	userDB, err := s.UserRepository.CreateUser(ctx, user)
	if err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToCreateUser, constants.Error, err.Error())
		return domain.UserDomain{}, err
	}

	s.LoggerSugar.Infow(constants.SuccessUserCreated, constants.UserID, userDB.ID)
	return userDB, nil
}

func (s *UserService) GetAllUsers(ctx domain.ContextControl) ([]domain.UserDomain, error) {
	users, err := s.UserRepository.GetAllUsers(ctx)
	if err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToGetAllUsers, constants.Error, err.Error())
		return nil, err
	}
	s.LoggerSugar.Infow(constants.SuccessUsersRetrieved, "count", len(users))
	return users, nil
}

func (s *UserService) GetUserByID(ctx domain.ContextControl, id uint64) (domain.UserDomain, error) {
	user, err := s.UserRepository.GetUserByID(ctx, id)
	if err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToGetUserByID, constants.Error, err.Error())
		return domain.UserDomain{}, err
	}
	s.LoggerSugar.Infow(constants.SuccessUserRetrieved, constants.UserID, user.ID)
	return user, nil
}

func (s *UserService) GetUserByUsername(ctx domain.ContextControl, username string) (domain.UserDomain, error) {
	userDB, err := s.UserRepository.GetUserByUsername(ctx, username)
	if err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToGetUserByUserName, constants.Error, err.Error())
		return domain.UserDomain{}, err
	}
	s.LoggerSugar.Infow(constants.SuccessUserRetrieved, constants.UserID, userDB.ID)
	return userDB, nil
}

func (s *UserService) GetUserByEmail(ctx domain.ContextControl, email string) (domain.UserDomain, error) {
	user, err := s.UserRepository.GetUserByEmail(ctx, email)
	if err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToGetUserByEmail, constants.Error, err.Error())
		return domain.UserDomain{}, err
	}
	s.LoggerSugar.Infow(constants.SuccessUserRetrieved, constants.UserID, user.ID)
	return user, nil
}

func (s *UserService) UpdateUser(ctx domain.ContextControl, user domain.UserDomain) (domain.UserDomain, error) {
	updateFields := make(map[string]interface{})

	if user.Name != "" {
		updateFields[constants.Name] = user.Name
	}
	if user.Username != "" {
		updateFields[constants.Username] = user.Username
	}
	if user.Email != "" {
		updateFields[constants.Email] = user.Email
	}
	if len(updateFields) == 0 {
		return domain.UserDomain{}, errors.New(constants.ErrorNoFieldsToUpdate)
	}
	updateFields[constants.UpdatedAt] = time.Now().UTC()

	return s.UserRepository.UpdateUser(ctx, user.ID, updateFields)
}

func (s *UserService) UpdateUserPassword(ctx domain.ContextControl, user domain.UserDomain, oldPassword, newPassword string) (domain.UserDomain, string, error) {
	userDB, err := s.UserRepository.GetUserByID(ctx, user.ID)
	if err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToGetUserByID, constants.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	if err := s.PasswordService.ComparePasswords(userDB.Password, oldPassword); err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToCompareHashAndPassword, constants.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	hashedPassword, err := s.PasswordService.HashPassword(newPassword)
	if err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToHashPassword, constants.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	fields := map[string]interface{}{
		constants.Password:  hashedPassword,
		constants.UpdatedAt: time.Now().UTC(),
	}

	updatedUser, err := s.UserRepository.UpdateUser(ctx, user.ID, fields)
	if err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToUpdatePassword, constants.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	tokenDomain := domain.TokenDomain{UserID: user.ID}
	token, err := s.TokenService.Create(ctx, tokenDomain)
	if err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToCreateToken, constants.Error, err.Error())
		return domain.UserDomain{}, "", err
	}
	tokenDomain.Token = token

	if err := s.TokenService.Save(ctx, tokenDomain); err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToSaveToken, constants.Error, err.Error())
		return domain.UserDomain{}, "", errors.New(constants.ErrorToSaveToken)
	}

	s.LoggerSugar.Infow(constants.SuccessPasswordUpdated, constants.UserID, updatedUser.ID)
	return updatedUser, tokenDomain.Token, nil
}

func (s *UserService) SoftDeleteUser(ctx domain.ContextControl, userID uint64) error {
	if err := s.UserRepository.SoftDeleteUser(ctx, userID); err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToSoftDeleteUser, constants.Error, err.Error())
		return err
	}

	tokenDomain := domain.TokenDomain{UserID: userID}
	if err := s.TokenService.Delete(ctx, tokenDomain); err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToDeleteToken, constants.Error, err.Error())
		return err
	}

	s.LoggerSugar.Infow(constants.SuccessUserSoftDeleted, constants.UserID, userID)
	return nil
}

func (s *UserService) NormalizeUserData(user *domain.UserDomain) domain.UserDomain {
	if user.Name != "" {
		user.Name = strings.TrimSpace(user.Name)
	}
	if user.Username != "" {
		user.Username = strings.TrimSpace(user.Username)
	}
	if user.Email != "" {
		user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	}
	return *user
}

func (s *UserService) validateCreateUserRequired(user domain.UserDomain, password string) error {
	if user.Name == "" {
		return errors.New(constants.NameIsRequired)
	}
	if user.Username == "" {
		return errors.New(constants.UsernameIsRequired)
	}
	if user.Email == "" {
		return errors.New(constants.EmailIsRequired)
	}
	if password == "" {
		return errors.New(constants.PasswordIsRequired)
	}
	if err := checkmail.ValidateFormat(user.Email); err != nil {
		return errors.New(constants.InvalidEmail)
	}
	return nil
}
