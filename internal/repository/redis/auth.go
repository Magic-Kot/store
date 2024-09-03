package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"time"
)

type AuthRepository struct {
	client redis.Client
}

func NewAuthRepository(client *redis.Client) *AuthRepository {
	return &AuthRepository{client: *client}
}

// CreateSession - creating a user session or updating it
func (a *AuthRepository) CreateSession(ctx context.Context, key string, value interface{}) (string, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("accessing Redis using the 'CreateSession' method")
	logger.Debug().Msgf("redis: create session by key: %s, value: %s", key, value)

	json := a.client.JSONSet(ctx, key, ".", value)

	a.client.Expire(ctx, key, 4*time.Hour) // Hard code

	if json.Err() == nil && json.Val() == "" {
		logger.Debug().Msgf("failed to create session. redis: %s", json.Err())
		return "", errors.New("failed to create session")
	}

	return json.String(), nil
}

// GetSession - получение сессии по userId пользователя
func (a *AuthRepository) GetSession(ctx context.Context, key string) (string, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("accessing Redis using the 'GetSession' method")
	logger.Debug().Msgf("redis: get session by key: %s", key)

	session := a.client.JSONGet(ctx, key)

	fmt.Printf("метод GetSession получил сессию: %s\n", session.String())

	if session.Err() == nil && session.Val() == "" {
		logger.Debug().Msgf("failed to get session: %s", session.Err())
		return "", errors.New("failed to get session")
	}

	return session.Result()
}

// DeleteSession - удаление сессии по userId пользователя
func (a *AuthRepository) DeleteSession(ctx context.Context, key string) error {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("accessing Redis using the 'DeleteSession' method")
	logger.Debug().Msgf("redis: delete session by key: %s", key)

	return nil
}
