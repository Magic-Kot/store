package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Magic-Kot/store/internal/models"
	"github.com/Magic-Kot/store/pkg/utils/hash"
	"github.com/Magic-Kot/store/pkg/utils/jwtToken"
)

type UserRepository interface {
	GetUser(ctx context.Context, user *models.User) (*models.User, error)
	CreateUser(ctx context.Context, login string, passwordHash string) (int, error)
	SignIn(ctx context.Context, user *models.UserAuthorization) (*models.UserAuthorization, error)
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
	user, err := s.UserRepository.GetUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// CreateUser - регистрация нового пользователя
func (s *UserService) CreateUser(ctx context.Context, login string, password string) (int, error) {
	passwordHash := hash.GenerateHash(password)

	id, err := s.UserRepository.CreateUser(ctx, login, passwordHash)
	if err != nil {
		return id, err
	}

	return id, nil
}

// SignIn - индетификация, аутентификация пользователя, получение токена
func (s *UserService) SignIn(ctx context.Context, user *models.UserAuthorization) (string, error) {
	passwordHash := hash.GenerateHash(user.Password)

	user, err := s.UserRepository.SignIn(ctx, user)
	if err != nil {
		return "", err
	}

	if passwordHash != user.Password {
		err = errors.New("invalid password")
		return "", err
	}

	token, err := jwtToken.GenerateToken(user.ID)

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
