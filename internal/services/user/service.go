package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/Magic-Kot/store/internal/models"
	"github.com/Magic-Kot/store/pkg/utils/hash"
)

type UserRepository interface {
	GetUser(ctx context.Context, id int) (string, string, error)
	CreateUser(ctx context.Context, login string, passwordHash string) (int, error)
	AuthorizationUser(ctx context.Context, login string) (string, error)
	UpdateUser(ctx context.Context, value string, arg []interface{}) error
	DeleteUser(ctx context.Context, id int) error
}

type UserService struct {
	UserRepository UserRepository
}

func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{
		UserRepository: userRepository,
	}
}

// GetUser - получение сущности пользователя по ID
func (s *UserService) GetUser(ctx context.Context, user *models.User) (*models.User, error) {
	login, email, err := s.UserRepository.GetUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	user.Username = login
	user.Email = email

	return user, nil
}

// CreateUser - регистрация нового пользователя
func (s *UserService) CreateUser(ctx context.Context, login string, password string) (int, error) {
	passwordHash := hash.GeneratePasswordHash(password)

	id, err := s.UserRepository.CreateUser(ctx, login, passwordHash)
	if err != nil {
		return id, err
	}

	return id, nil
}

// AuthorizationUser - авторизация пользователя
func (s *UserService) AuthorizationUser(ctx context.Context, login string, password string) (int, error) {
	passwordHash := hash.GeneratePasswordHash(password)

	passwordHashDB, err := s.UserRepository.AuthorizationUser(ctx, login)
	if err != nil {
		return 0, err
	}

	if passwordHash != passwordHashDB {
		err = fmt.Errorf("invalid password")
		return 0, err
	}

	token := 0

	return token, nil
}

// UpdateUser - обновление данных пользователя по ID
func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	value := make([]string, 0)
	arg := make([]interface{}, 0)
	argId := 2

	arg = append(arg, user.ID)

	if user.Username != "" {
		value = append(value, fmt.Sprintf("username=$%d", argId)) //username=$2
		arg = append(arg, user.Username)
		argId++
	}

	if user.Email != "" {
		value = append(value, fmt.Sprintf("email=$%d", argId)) //username=$2 email=$3
		arg = append(arg, user.Email)
		argId++
	}

	valueQuery := strings.Join(value, ", ")

	err := s.UserRepository.UpdateUser(ctx, valueQuery, arg)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser - удаление пользователя по ID
func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	err := s.UserRepository.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
