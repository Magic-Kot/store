package user

import (
	"context"
	"github.com/rs/zerolog"
	"strconv"
)

type UserRepository interface {
	DeleteUser(ctx context.Context, id int) error
}

type AuthRepository interface {
	DeleteSession(ctx context.Context, key string) error
}

type UserService struct {
	UserRepository UserRepository
	AuthRepository AuthRepository
}

func NewUserService(userRepository UserRepository, authRepository AuthRepository) *UserService {
	return &UserService{
		UserRepository: userRepository,
		AuthRepository: authRepository,
	}
}

// DeleteUser - deleting a user by id
func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("starting the 'DeleteUser' service")

	err := s.UserRepository.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	err = s.AuthRepository.DeleteSession(ctx, strconv.Itoa(id))
	if err != nil {
		return err
	}

	return nil
}
