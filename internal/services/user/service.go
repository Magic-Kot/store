package user

import (
	"context"
)

type UserRepository interface {
	GetUser(ctx context.Context, user_name string) (string, error)
	CreateUser(ctx context.Context, name string, login string) (string, error)
	UpdateUser(ctx context.Context, name string, login string) (string, error)
	DeleteUser(ctx context.Context, name string) (string, error)
}

type UserService struct {
	UserRepository UserRepository
}

func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{
		UserRepository: userRepository,
	}
}

func (s *UserService) GetUser(ctx context.Context, name string) (string, error) {
	name, err := s.UserRepository.GetUser(ctx, name)
	if err != nil {
		return "", err
	}

	return name, nil
}

func (s *UserService) CreateUser(ctx context.Context, name string, login string) (string, error) {
	rez, err := s.UserRepository.CreateUser(ctx, name, login)
	if err != nil {
		return rez, err
	}

	return rez, nil
}

func (s *UserService) UpdateUser(ctx context.Context, name string, login string) (string, error) {
	rez, err := s.UserRepository.UpdateUser(ctx, name, login)
	if err != nil {
		return rez, err
	}

	return rez, nil
}

func (s *UserService) DeleteUser(ctx context.Context, name string) (string, error) {
	rez, err := s.UserRepository.DeleteUser(ctx, name)
	if err != nil {
		return rez, err
	}

	return rez, nil
}
