package user

import (
	"context"
	"encoding/json"
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
	UpdateUser(ctx context.Context, table string, column string, value string, arg []interface{}) error
	DeleteUser(ctx context.Context, id int) error
}

type AuthRepository interface {
	CreateSession(ctx context.Context, key string, value interface{}) (string, error)
	GetSession(ctx context.Context, key string) (string, error)
	DeleteSession(ctx context.Context, key string) error
}

type UserService struct {
	UserRepository UserRepository
	AuthRepository AuthRepository
	token          *jwt_token.Manager
}

func NewUserService(userRepository UserRepository, authRepository AuthRepository, token *jwt_token.Manager) *UserService {
	return &UserService{
		UserRepository: userRepository,
		AuthRepository: authRepository,
		token:          token,
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

// SignUp - registering a new user
func (s *UserService) SignUp(ctx context.Context, login string, password string) (int, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("starting the 'SignUp' service")

	passwordHash := hash.GenerateHash(password)

	id, err := s.UserRepository.CreateUser(ctx, login, passwordHash)
	if err != nil {
		return id, err
	}

	return id, nil
}

// SignIn - user authentication
func (s *UserService) SignIn(ctx context.Context, user *models.UserAuthorization) (models.Tokens, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("starting the 'SignIn' service")

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

	_, err = s.AuthRepository.CreateSession(ctx, strconv.Itoa(user.ID), session)
	if err != nil {
		logger.Debug().Msgf("failed to create session: %s", err)
		return res, errors.New("failed to create session")
	}

	return res, nil
}

// RefreshToken - getting new refresh and access tokens
func (s *UserService) RefreshToken(ctx context.Context, refresh string) (models.Tokens, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("starting the 'RefreshToken' service")

	var (
		res   models.Tokens
		token models.Session
	)

	passwordDecode, err := s.token.ParseRefreshToken(refresh)
	if err != nil {
		logger.Debug().Msgf("invalid refresh token: %s", err)
		// redirecting to the authorization page
		return res, errAutorizationUser
	}

	userId := strings.Fields(passwordDecode)

	session, err := s.AuthRepository.GetSession(ctx, userId[1])

	err = json.Unmarshal([]byte(session), &token)
	if err != nil {
		logger.Debug().Msgf("invalid refresh token: %s", err)
		// redirecting to the authorization page
		return res, errAutorizationUser
	}

	err = hash.CompareHashBcrypt(refresh, token.RefreshToken)
	if err != nil {
		logger.Debug().Msgf("invalid refresh token: %s", err)
		// redirecting to the authorization page
		return res, errAutorizationUser
	}

	// generation of new access and refresh tokens
	res.AccessToken, err = s.token.NewJWT(userId[1])
	if err != nil {
		return res, err
	}

	res.RefreshToken = s.token.NewRefreshToken(userId[1])

	RefreshTokenHash, err := hash.GenerateHashBcrypt(res.RefreshToken)
	if err != nil {
		return res, err
	}

	newSession := models.Session{
		RefreshToken: RefreshTokenHash,
		ExpiresAt:    time.Now().Add(s.token.RefreshTokenTTL()),
	}

	// update session
	_, err = s.AuthRepository.CreateSession(ctx, userId[1], newSession)
	if err != nil {
		logger.Debug().Msgf("failed to update session: %s", err)
		return res, errors.New("failed to update session")
	}

	return res, nil
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
