package user

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Magic-Kot/store/internal/models"
	"github.com/Magic-Kot/store/pkg/utils/hash"
	"github.com/Magic-Kot/store/pkg/utils/jwt_token"

	"github.com/rs/zerolog"
)

var (
	errAutorizationUser = errors.New("invalid refresh token")
)

type UserRepository interface {
	GetUser(ctx context.Context, user *models.User) (*models.User, error)
	CreateUser(ctx context.Context, login string, passwordHash string) (int, error)
	SignIn(ctx context.Context, user *models.UserAuthorization) (*models.UserAuthorization, error)
	CreateSession(ctx context.Context, value string, arg []interface{}) (int, error)
	GetSession(ctx context.Context, table string, column string, value string, arg []interface{}) (string, error)
	UpdateUser(ctx context.Context, table string, column string, value string, arg []interface{}) error
	DeleteUser(ctx context.Context, id int) error
}

type UserService struct {
	UserRepository UserRepository
	token          *jwt_token.Manager
}

func NewUserService(userRepository UserRepository, token *jwt_token.Manager) *UserService {
	return &UserService{
		UserRepository: userRepository,
		token:          token,
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

// SignIn - аутентификация пользователя, получение токена
func (s *UserService) SignIn(ctx context.Context, user *models.UserAuthorization) (models.Tokens, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("starting the handler 'SignIn'")

	var res models.Tokens

	passwordHash := hash.GenerateHash(user.Password)

	user, err := s.UserRepository.SignIn(ctx, user)
	if err != nil {
		return res, err
	}

	if passwordHash != user.Password {
		logger.Debug().Msg("invalid password")

		err = errors.New("invalid password")
		return res, err
	}

	res.AccessToken, err = s.token.NewJWT(strconv.Itoa(user.ID))
	if err != nil {
		return res, err
	}

	res.RefreshToken = s.token.NewRefreshToken(strconv.Itoa(user.ID))

	RefreshTokenHash, err := hash.GenerateHashBcrypt(res.RefreshToken)
	if err != nil {
		return res, err
	}

	session := models.Session{
		RefreshToken: RefreshTokenHash,
		ExpiresAt:    time.Now().Add(s.token.RefreshTokenTTL()),
	}

	// checking for an open session
	// GetSession - метод возвращает ошибку в случае отсутствия сессии
	arg := make([]interface{}, 0)
	arg = append(arg, user.ID)

	_, err = s.UserRepository.GetSession(ctx, "sessions", "userId", "1", arg)

	if err != nil {
		arg = append(arg, user.GUID, session.RefreshToken, session.ExpiresAt)

		_, err := s.UserRepository.CreateSession(ctx, "userId, guid, refreshToken, expiresAt", arg)
		if err != nil {
			logger.Debug().Msgf("failed to create session: %s", err)

			return res, err
		}

		return res, nil
	}

	// update session
	arg = append(arg, session.RefreshToken, session.ExpiresAt)

	err = s.UserRepository.UpdateUser(ctx, "sessions", "userId", "refreshToken=$2, expiresAt=$3", arg)
	if err != nil {
		logger.Debug().Msgf("failed to update session: %s", err)

		return res, err
	}

	return res, nil
}

// RefreshToken - получение новых refresh и access токенов
func (s *UserService) RefreshToken(ctx context.Context, refresh string) (models.Tokens, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("starting the handler 'RefreshToken'")

	var res models.Tokens

	arg := make([]interface{}, 0)

	passwordDecode, err := s.token.ParseRefreshToken(refresh)
	if err != nil {
		logger.Debug().Msgf("invalid refresh token: %s", err)
		// redirecting to the authorization page
		return res, errAutorizationUser
	}

	userId := strings.Fields(passwordDecode)

	arg = append(arg, userId[1])

	passwordHash, err := s.UserRepository.GetSession(ctx, "sessions", "userId", "refreshToken", arg)
	if err != nil {
		logger.Debug().Msgf("invalid refresh token: %s", err)
		// redirecting to the authorization page
		return res, errAutorizationUser
	}

	err = hash.CompareHashBcrypt(refresh, passwordHash)
	if err != nil {
		logger.Debug().Msgf("invalid refresh token: %s", err)
		// redirecting to the authorization page
		return res, errAutorizationUser
	}

	// TODO: код ниже вынести в отдельную функцию

	res.AccessToken, err = s.token.NewJWT(userId[1])
	if err != nil {
		return res, err
	}

	res.RefreshToken = s.token.NewRefreshToken(userId[1])

	RefreshTokenHash, err := hash.GenerateHashBcrypt(res.RefreshToken)
	if err != nil {
		return res, err
	}

	session := models.Session{
		RefreshToken: RefreshTokenHash,
		ExpiresAt:    time.Now().Add(s.token.RefreshTokenTTL()),
	}

	// update session
	arg = append(arg, session.RefreshToken, session.ExpiresAt)

	err = s.UserRepository.UpdateUser(ctx, "sessions", "userId", "refreshToken=$2, expiresAt=$3", arg)
	if err != nil {
		return res, err
	}

	return res, nil
}

// UpdateUser - обновление данных пользователя по id
func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("starting the service 'UpdateUser'")

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
			fmt.Println(types.Field(i).Name, values.Field(i))

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

// DeleteUser - удаление пользователя по id
func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	err := s.UserRepository.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
