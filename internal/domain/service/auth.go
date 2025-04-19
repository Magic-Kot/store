package service

import (
	"context"
	"fmt"
	"git.appkode.ru/pub/go/failure"
	"time"

	"github.com/rs/xid"

	"github.com/Magic-Kot/store/internal/domain/entity"
	"github.com/Magic-Kot/store/internal/domain/value"
	"github.com/Magic-Kot/store/pkg/errcodes"
	"github.com/Magic-Kot/store/pkg/utils/hash"
)

type refreshTokenParser interface {
	Parse(token value.RefreshToken) (value.RefreshTokenClaims, error)
}

type accessTokenGenerator interface {
	Generate(value.PersonID) (value.AccessToken, error)
}

type refreshTokenGenerator interface {
	Generate(value.PersonID, value.RefreshTokenID) (value.RefreshToken, error)
}

type refreshTokenSource interface {
	Create(context.Context, value.PersonID, value.RefreshTokenID, time.Duration) error
	Find(context.Context, value.PersonID, value.RefreshTokenID) error
	Delete(context.Context, value.PersonID, value.RefreshTokenID) error
}

type postgresRepository interface {
	CreateUser(context.Context, entity.CreateUser) error
	UserByLogin(context.Context, value.Login) (value.UserAuth, error)
	PersonIDLatest(context.Context) (value.PersonID, error)
}

type Auth struct {
	refreshTokenParser    refreshTokenParser
	accessTokenGenerator  accessTokenGenerator
	refreshTokenGenerator refreshTokenGenerator
	refreshTokens         refreshTokenSource
	refreshTokenTTL       time.Duration
	PostgresRepository    postgresRepository
}

func NewAuth(
	refreshTokenParser refreshTokenParser,
	accessTokenGenerator accessTokenGenerator,
	refreshTokenGenerator refreshTokenGenerator,
	refreshTokens refreshTokenSource,
	refreshTokenTTL time.Duration,
	postgresRepository postgresRepository,
) Auth {
	return Auth{
		refreshTokenParser:    refreshTokenParser,
		accessTokenGenerator:  accessTokenGenerator,
		refreshTokenGenerator: refreshTokenGenerator,
		refreshTokens:         refreshTokens,
		refreshTokenTTL:       refreshTokenTTL,
		PostgresRepository:    postgresRepository,
	}
}

func (a Auth) Registration(ctx context.Context, login value.Login, password value.Password) error {
	personID, err := a.PostgresRepository.PersonIDLatest(ctx)
	if err != nil && !failure.IsNotFoundError(err) {
		return fmt.Errorf("PostgresRepository.PersonIDLatest: %w", err)
	}

	personID = 10000

	user := entity.CreateUser{
		ID:           xid.New().String(),
		PersonID:     personID + 1,
		Login:        login,
		PasswordHash: hash.GenerateHash(password),
		CreatedAt:    time.Now(),
	}

	if err := a.PostgresRepository.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("PostgresRepository.CreateUser: %w", err)
	}

	return nil
}

func (a Auth) Authenticate(ctx context.Context, login value.Login, password value.Password) (value.TokenPair, error) {
	passwordHash := hash.GenerateHash(password)

	user, err := a.PostgresRepository.UserByLogin(ctx, login)
	if err != nil {
		return value.TokenPair{}, fmt.Errorf("PostgresRepository.UserByLogin: %w", err)
	}

	if passwordHash != user.Password.String() {
		return value.TokenPair{}, failure.NewInvalidArgumentError(
			"invalid login or password",
			failure.WithCode(errcodes.DirectoriesBusy),
			failure.WithDescription("invalid login or password"),
		)
	}

	refreshTokenID := value.NewRefreshTokenID(xid.New())

	pair, err := a.tokenPair(user.PersonID, refreshTokenID)
	if err != nil {
		return value.TokenPair{}, fmt.Errorf("tokenPair: %w", err)
	}

	if err = a.refreshTokens.Create(ctx, user.PersonID, refreshTokenID, a.refreshTokenTTL); err != nil {
		return value.TokenPair{}, fmt.Errorf("refreshTokens.Create: %w", err)
	}

	return value.TokenPair{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
	}, nil
}

func (a Auth) Refresh(ctx context.Context, refreshToken value.RefreshToken) (value.TokenPair, error) {
	refreshClaims, err := a.refreshTokenParser.Parse(refreshToken)
	if err != nil {
		return value.TokenPair{}, fmt.Errorf("refreshTokenParser.Parse: %w", err)
	}

	refreshTokenID, err := value.ParseRefreshTokenID(refreshClaims.ID)
	if err != nil {
		return value.TokenPair{}, fmt.Errorf("value.ParseRefreshTokenID: %w", err)
	}

	if err = a.refreshTokens.Find(ctx, refreshClaims.PersonID, refreshTokenID); err != nil {
		return value.TokenPair{}, fmt.Errorf("refreshTokens.Find: %w", err)
	}

	NewRefreshTokenID := value.NewRefreshTokenID(xid.New())

	pair, err := a.tokenPair(refreshClaims.PersonID, NewRefreshTokenID)
	if err != nil {
		return value.TokenPair{}, fmt.Errorf("tokenPair: %w", err)
	}

	if err = a.refreshTokens.Create(ctx, refreshClaims.PersonID, NewRefreshTokenID, a.refreshTokenTTL); err != nil {
		return value.TokenPair{}, fmt.Errorf("refreshTokens.Create: %w", err)
	}

	if err = a.refreshTokens.Delete(ctx, refreshClaims.PersonID, refreshTokenID); err != nil {
		return value.TokenPair{}, fmt.Errorf("refreshTokens.Delete: %w", err)
	}

	return value.TokenPair{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
	}, nil
}

func (a Auth) Logout(ctx context.Context, refreshToken value.RefreshToken) error {
	refreshClaims, err := a.refreshTokenParser.Parse(refreshToken)
	if err != nil {
		return fmt.Errorf("refreshTokenParser.Parse: %w", err)
	}

	refreshTokenID, err := value.ParseRefreshTokenID(refreshClaims.ID)
	if err != nil {
		return fmt.Errorf("value.ParseRefreshTokenID: %w", err)
	}

	if err = a.refreshTokens.Delete(ctx, refreshClaims.PersonID, refreshTokenID); err != nil {
		return fmt.Errorf("refreshTokens.Delete: %w", err)
	}

	return nil
}

func (a Auth) tokenPair(
	personID value.PersonID,
	refreshTokenID value.RefreshTokenID,
) (value.TokenPair, error) {
	access, err := a.accessTokenGenerator.Generate(personID)
	if err != nil {
		return value.TokenPair{}, fmt.Errorf("accessTokenGenerator.Generate: %w", err)
	}

	refresh, err := a.refreshTokenGenerator.Generate(personID, refreshTokenID)
	if err != nil {
		return value.TokenPair{}, fmt.Errorf("refreshTokenGenerator.Generate: %w", err)
	}

	return value.TokenPair{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}
