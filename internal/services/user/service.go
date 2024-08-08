package user

import (
	"context"
)

type UserRepository interface {
	GetUser(ctx context.Context, login string) (int, error)
	CreateUser(ctx context.Context, login string, password string) (string, error)
	UpdateUser(ctx context.Context, id int, login string, password string) (string, error)
	DeleteUser(ctx context.Context, id int) (string, error)
}

type UserService struct {
	UserRepository UserRepository
}

func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{
		UserRepository: userRepository,
	}
}

func (s *UserService) GetUser(ctx context.Context, login string) (int, error) {
	id, err := s.UserRepository.GetUser(ctx, login)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *UserService) CreateUser(ctx context.Context, login string, password string) (string, error) {
	rez, err := s.UserRepository.CreateUser(ctx, login, password)
	if err != nil {
		return rez, err
	}

	return rez, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int, login string, password string) (string, error) {
	rez, err := s.UserRepository.UpdateUser(ctx, id, login, password)
	if err != nil {
		return rez, err
	}

	return rez, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int) (string, error) {
	rez, err := s.UserRepository.DeleteUser(ctx, id)
	if err != nil {
		return rez, err
	}

	return rez, nil
}
