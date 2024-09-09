package referral

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/Magic-Kot/store/internal/models"
	"github.com/Magic-Kot/store/pkg/utils/short_url"

	"github.com/rs/zerolog"
)

var (
	errNotFound = errors.New("the short url was not found")
)

type ReferralRepositoryPostgres interface {
	CreateShortUrl(ctx context.Context, userId int, arg interface{}) error
}

type ReferralRepositoryRedis interface {
	CreateReferral(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error)
	GetReferral(ctx context.Context, key string) (string, error)
}

type ReferralService struct {
	ReferralRepositoryPostgres ReferralRepositoryPostgres
	ReferralRepositoryRedis    ReferralRepositoryRedis
}

func NewReferralService(ReferralRepositoryPostgres ReferralRepositoryPostgres, ReferralRepositoryRedis ReferralRepositoryRedis) *ReferralService {
	return &ReferralService{
		ReferralRepositoryPostgres: ReferralRepositoryPostgres,
		ReferralRepositoryRedis:    ReferralRepositoryRedis,
	}
}

// CreateReferral - generating a referral link
func (rs *ReferralService) CreateReferral(ctx context.Context, body *models.Request) (models.Response, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("starting the 'CreateReferral' service")

	resp := models.Response{}
	var id string

	if body.CustomShort == "" {
		id = short_url.Base62Encode(rand.Uint64())
	} else {
		id = body.CustomShort
	}

	// checking redis for collisions
	url, _ := rs.ReferralRepositoryRedis.GetReferral(ctx, id)
	if url != "" {
		logger.Debug().Msg("the short url was not found")
		return resp, errNotFound
	}

	if body.Expiry == 0 {
		body.Expiry = 24 * time.Hour // Hard code
	}

	_, err := rs.ReferralRepositoryRedis.CreateReferral(ctx, id, body.URL, body.Expiry)
	if err != nil {
		logger.Debug().Msg("URL short is already in use")
		return resp, err
	}

	// saving a short referral link in Postgres
	err = rs.ReferralRepositoryPostgres.CreateShortUrl(ctx, body.UserId, id)

	resp = models.Response{
		URL:         body.URL,
		CustomShort: "",
		Expiry:      body.Expiry,
	}

	//resp.CustomShort = os.Getenv("DOMAIN") + "/" + id
	resp.CustomShort = "https://siriusfuture.ru/baf/" + id

	return resp, nil
}

func (rs *ReferralService) GetReferral(ctx context.Context, shortUrl string) (string, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("starting the 'GetReferral' service")

	url, _ := rs.ReferralRepositoryRedis.GetReferral(ctx, shortUrl)
	if url == "" {
		logger.Debug().Msg("the short url was not found")
		return "", errNotFound
	}

	return url, nil
}
