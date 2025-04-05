package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	redis "github.com/redis/go-redis/v9"

	"github.com/Magic-Kot/store/internal/domain/value"
)

type RefreshTokens struct {
	redisClient *redis.Client
}

func NewRefreshTokens(
	redisClient *redis.Client,
) RefreshTokens {
	return RefreshTokens{
		redisClient: redisClient,
	}
}

func (r RefreshTokens) Create(
	ctx context.Context,
	personID value.PersonID,
	refreshTokenID value.RefreshTokenID,
	expiration time.Duration,
) error {
	if err := r.redisClient.Set(
		ctx,
		r.redisKey(personID, refreshTokenID),
		"",
		expiration,
	).Err(); err != nil {
		return fmt.Errorf("redisClient.Set: %w", err)
	}

	return nil
}

func (r RefreshTokens) Find(
	ctx context.Context,
	personID value.PersonID,
	refreshTokenID value.RefreshTokenID,
) error {
	if err := r.redisClient.Get(ctx, r.redisKey(personID, refreshTokenID)).Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("refresh token not found")
		}

		return fmt.Errorf("redisClient.Get: %w", err)
	}

	return nil
}

func (r RefreshTokens) Delete(
	ctx context.Context,
	personID value.PersonID,
	refreshTokenID value.RefreshTokenID,
) error {
	if err := r.redisClient.Del(ctx, r.redisKey(personID, refreshTokenID)).Err(); err != nil {
		return fmt.Errorf("redisClient.Del: %w", err)
	}

	return nil
}

func (r RefreshTokens) redisKey(userID value.PersonID, refreshTokenID value.RefreshTokenID) string {
	return fmt.Sprintf("token:refresh:%s:%s", userID, refreshTokenID)
}
