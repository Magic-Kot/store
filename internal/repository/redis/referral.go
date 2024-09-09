package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

var (
	errCreateReferral = errors.New("error saving the referral link")
	errNotFound       = errors.New("the short url was not found")
	errGetReferral    = errors.New("url short is already in use")
)

type ReferralRepository struct {
	client redis.Client
}

func NewReferralRepository(client *redis.Client) *ReferralRepository {
	return &ReferralRepository{client: *client}
}

// CreateReferral - saving a new referral link
func (r *ReferralRepository) CreateReferral(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("accessing Redis using the 'CreateReferral' method")
	logger.Debug().Msgf("redis: create session by id: %s", key)

	json := r.client.Set(ctx, key, value, expiration)

	logger.Debug().Msgf("redis create. value: %s, err: %s", json.Val(), json.Err())

	if json.Err() == nil && json.Val() == "" {
		logger.Debug().Msgf("failed to create session. redis: %s", json.Err())
		return "", errCreateReferral
	}

	return json.Val(), nil
}

// GetReferral - getting a referral link by id
func (r *ReferralRepository) GetReferral(ctx context.Context, key string) (string, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("accessing Redis using the 'GetReferral' method")
	logger.Debug().Msgf("redis: get session by key: %s", key)

	referralURL := r.client.Get(ctx, key)

	logger.Debug().Msgf("redis get. value: %s, err: %s", referralURL.Val(), referralURL.Err())

	return referralURL.Result()
}
