package service

import (
	"errors"
	"fmt"
	"github.com/lechitz/AionApi/core/ports/output/db"
	security2 "github.com/lechitz/AionApi/core/ports/output/security"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/lechitz/AionApi/core/domain"
	"github.com/lechitz/AionApi/core/msg"
	"github.com/lechitz/AionApi/pkg/contextkeys"
	"go.uber.org/zap"
)

type UserService struct {
	UserRepository  db.IUserRepository
	TokenService    security2.ITokenService
	PasswordService security2.IPasswordService
	LoggerSugar     *zap.SugaredLogger
}

func NewUserService(userRepo db.IUserRepository, tokenService security2.ITokenService, passwordService security2.IPasswordService, loggerSugar *zap.SugaredLogger) *UserService {
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
		s.LoggerSugar.Errorw(msg.ErrorToValidateCreateUser, contextkeys.Error, err.Error())
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
		s.LoggerSugar.Errorw(msg.ErrorToHashPassword, contextkeys.Error, err.Error())
		return domain.UserDomain{}, err
	}
	user.Password = hashPassword

	userDB, err := s.UserRepository.CreateUser(ctx, user)
	if err != nil {
		s.LoggerSugar.Errorw(msg.ErrorToCreateUser, contextkeys.Error, err.Error())
		return domain.UserDomain{}, err
	}

	s.LoggerSugar.Infow(msg.SuccessUserCreated, contextkeys.UserID, userDB.ID)
	return userDB, nil
}

func (s *UserService) GetAllUsers(ctx domain.ContextControl) ([]domain.UserDomain, error) {
	users, err := s.UserRepository.GetAllUsers(ctx)
	if err != nil {
		s.LoggerSugar.Errorw(msg.ErrorToGetAllUsers, contextkeys.Error, err.Error())
		return nil, err
	}
	s.LoggerSugar.Infow(msg.SuccessUsersRetrieved, "count", len(users))
	return users, nil
}

func (s *UserService) GetUserByID(ctx domain.ContextControl, id uint64) (domain.UserDomain, error) {
	user, err := s.UserRepository.GetUserByID(ctx, id)
	if err != nil {
		s.LoggerSugar.Errorw(msg.ErrorToGetUserByID, contextkeys.Error, err.Error())
		return domain.UserDomain{}, err
	}
	s.LoggerSugar.Infow(msg.SuccessUserRetrieved, contextkeys.UserID, user.ID)
	return user, nil
}

func (s *UserService) GetUserByUsername(ctx domain.ContextControl, username string) (domain.UserDomain, error) {
	userDB, err := s.UserRepository.GetUserByUsername(ctx, username)
	if err != nil {
		s.LoggerSugar.Errorw(msg.ErrorToGetUserByUserName, contextkeys.Error, err.Error())
		return domain.UserDomain{}, err
	}
	s.LoggerSugar.Infow(msg.SuccessUserRetrieved, contextkeys.UserID, userDB.ID)
	return userDB, nil
}

func (s *UserService) GetUserByEmail(ctx domain.ContextControl, email string) (domain.UserDomain, error) {
	user, err := s.UserRepository.GetUserByEmail(ctx, email)
	if err != nil {
		s.LoggerSugar.Errorw(msg.ErrorToGetUserByEmail, contextkeys.Error, err.Error())
		return domain.UserDomain{}, err
	}
	s.LoggerSugar.Infow(msg.SuccessUserRetrieved, contextkeys.UserID, user.ID)
	return user, nil
}

func (s *UserService) UpdateUser(ctx domain.ContextControl, user domain.UserDomain) (domain.UserDomain, error) {
	updateFields := make(map[string]interface{})

	if user.Name != "" {
		updateFields[contextkeys.Name] = user.Name
	}
	if user.Username != "" {
		updateFields[contextkeys.Username] = user.Username
	}
	if user.Email != "" {
		updateFields[contextkeys.Email] = user.Email
	}
	if len(updateFields) == 0 {
		return domain.UserDomain{}, errors.New(msg.ErrorNoFieldsToUpdate)
	}
	updateFields[contextkeys.UpdatedAt] = time.Now().UTC()

	return s.UserRepository.UpdateUser(ctx, user.ID, updateFields)
}

func (s *UserService) UpdateUserPassword(ctx domain.ContextControl, user domain.UserDomain, oldPassword, newPassword string) (domain.UserDomain, string, error) {
	userDB, err := s.UserRepository.GetUserByID(ctx, user.ID)
	if err != nil {
		s.LoggerSugar.Errorw(msg.ErrorToGetUserByID, contextkeys.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	if err := s.PasswordService.ComparePasswords(userDB.Password, oldPassword); err != nil {
		s.LoggerSugar.Errorw(msg.ErrorToCompareHashAndPassword, contextkeys.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	hashedPassword, err := s.PasswordService.HashPassword(newPassword)
	if err != nil {
		s.LoggerSugar.Errorw(msg.ErrorToHashPassword, contextkeys.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	fields := map[string]interface{}{
		contextkeys.Password:  hashedPassword,
		contextkeys.UpdatedAt: time.Now().UTC(),
	}

	updatedUser, err := s.UserRepository.UpdateUser(ctx, user.ID, fields)
	if err != nil {
		s.LoggerSugar.Errorw(msg.ErrorToUpdatePassword, contextkeys.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	tokenDomain := domain.TokenDomain{UserID: user.ID}
	token, err := s.TokenService.Create(ctx, tokenDomain)
	if err != nil {
		s.LoggerSugar.Errorw(msg.ErrorToCreateToken, contextkeys.Error, err.Error())
		return domain.UserDomain{}, "", err
	}
	tokenDomain.Token = token

	if err := s.TokenService.Save(ctx, tokenDomain); err != nil {
		s.LoggerSugar.Errorw(msg.ErrorToSaveToken, contextkeys.Error, err.Error())
		return domain.UserDomain{}, "", errors.New(msg.ErrorToSaveToken)
	}

	s.LoggerSugar.Infow(msg.SuccessPasswordUpdated, contextkeys.UserID, updatedUser.ID)
	return updatedUser, tokenDomain.Token, nil
}

func (s *UserService) SoftDeleteUser(ctx domain.ContextControl, userID uint64) error {
	if err := s.UserRepository.SoftDeleteUser(ctx, userID); err != nil {
		s.LoggerSugar.Errorw(msg.ErrorToSoftDeleteUser, contextkeys.Error, err.Error())
		return err
	}

	tokenDomain := domain.TokenDomain{UserID: userID}
	if err := s.TokenService.Delete(ctx, tokenDomain); err != nil {
		s.LoggerSugar.Errorw(msg.ErrorToDeleteToken, contextkeys.Error, err.Error())
		return err
	}

	s.LoggerSugar.Infow(msg.SuccessUserSoftDeleted, contextkeys.UserID, userID)
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
		return errors.New(msg.NameIsRequired)
	}
	if user.Username == "" {
		return errors.New(msg.UsernameIsRequired)
	}
	if user.Email == "" {
		return errors.New(msg.EmailIsRequired)
	}
	if password == "" {
		return errors.New(msg.PasswordIsRequired)
	}
	if err := checkmail.ValidateFormat(user.Email); err != nil {
		return errors.New(msg.InvalidEmail)
	}
	return nil
}
