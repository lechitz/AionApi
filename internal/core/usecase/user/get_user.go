package user

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain/entity"
	"github.com/lechitz/AionApi/internal/def"

	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
)

// Retriever is an interface for retrieving user data from the system.
// GetUserByID fetches a user by their unique identifier and returns the corresponding domain.UserDomain.
// GetUserByEmail retrieves a user by their email address and returns the associated domain.UserDomain.
// GetUserByUsername fetches a user by their username and returns the related domain.UserDomain.
// GetAllUsers retrieves all users available in the system and returns a slice of domain.UserDomain.
type Retriever interface {
	GetUserByID(ctx context.Context, userID uint64) (entity.UserDomain, error)
	GetUserByEmail(ctx context.Context, email string) (entity.UserDomain, error)
	GetUserByUsername(ctx context.Context, username string) (entity.UserDomain, error)
	GetAllUsers(ctx context.Context) ([]entity.UserDomain, error)
}

// GetUserByID retrieves a user by their unique ID from the database. Returns the user or an error if the operation fails.
func (s *Service) GetUserByID(ctx context.Context, userID uint64) (entity.UserDomain, error) {
	user, err := s.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		s.logger.Errorw(constants.ErrorToGetUserByID, def.Error, err.Error())
		return entity.UserDomain{}, err
	}

	s.logger.Infow(constants.SuccessUserRetrieved, def.CtxUserID, user.ID)

	return user, nil
}

// GetUserByEmail retrieves a user by their email address from the database. Returns the user or an error if the operation fails.
func (s *Service) GetUserByEmail(ctx context.Context, email string) (entity.UserDomain, error) {
	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		s.logger.Errorw(constants.ErrorToGetUserByEmail, def.Error, err.Error())
		return entity.UserDomain{}, err
	}

	s.logger.Infow(constants.SuccessUserRetrieved, def.CtxUserID, user.ID)

	return user, nil
}

// GetUserByUsername retrieves a user by their username from the database. Returns the user or an error if the operation fails.
func (s *Service) GetUserByUsername(
	ctx context.Context,
	username string,
) (entity.UserDomain, error) {
	userDB, err := s.userRepository.GetUserByUsername(ctx, username)
	if err != nil {
		s.logger.Errorw(constants.ErrorToGetUserByUserName, def.Error, err.Error())
		return entity.UserDomain{}, err
	}

	s.logger.Infow(constants.SuccessUserRetrieved, def.CtxUserID, userDB.ID)

	return userDB, nil
}

// GetAllUsers retrieves all users from the system. Returns a slice of UserDomain or an error if the operation fails.
func (s *Service) GetAllUsers(ctx context.Context) ([]entity.UserDomain, error) {
	users, err := s.userRepository.GetAllUsers(ctx)
	if err != nil {
		s.logger.Errorw(constants.ErrorToGetAllUsers, def.Error, err.Error())
		return nil, err
	}

	s.logger.Infow(constants.SuccessUsersRetrieved, def.CtxUsers, len(users))

	return users, nil
}
