package user

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/Magic-Kot/store/internal/models"
	"github.com/rs/zerolog"
)

type UserRepository interface {
	GetUser(ctx context.Context, user *models.User) (*models.User, error)
	UpdateUser(ctx context.Context, table string, column string, value string, arg []interface{}) error
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

// GetUser - getting a user by id
func (s *UserService) GetUser(ctx context.Context, user *models.User) (*models.User, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("starting the 'GetUser' service")

	user, err := s.UserRepository.GetUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser - updating user data by ID
func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("starting the 'UpdateUser' service")

	value := make([]string, 0)
	arg := make([]interface{}, 0)
	argId := 2

	arg = append(arg, user.ID)

	values := reflect.ValueOf(*user)
	types := values.Type()

	if user.Age != 0 {
		value = append(value, fmt.Sprintf("age=$%d", argId)) //age=$2
		arg = append(arg, user.Age)
		argId++
	}

	for i := 2; i < values.NumField(); i++ {
		if values.Field(i).String() != "" {
			value = append(value, fmt.Sprintf("%s=$%d", types.Field(i).Name, argId))
			arg = append(arg, values.Field(i).String())

			argId++
		}
	}

	valueQuery := strings.Join(value, ", ")

	err := s.UserRepository.UpdateUser(ctx, "users", "id", valueQuery, arg)
	if err != nil {
		return err
	}

	return nil
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
