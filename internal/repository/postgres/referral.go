package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Magic-Kot/store/pkg/client/postg"

	"github.com/rs/zerolog"
)

var (
	errTransaction = errors.New("transaction error")
)

type ReferralRepository struct {
	client postg.Client
}

func NewReferralRepository(client postg.Client) *ReferralRepository {
	return &ReferralRepository{
		client: client,
	}
}

// CreateShortUrl - saving a short referral link to the "referral" table
func (r *ReferralRepository) CreateShortUrl(ctx context.Context, userId int, arg interface{}) error {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("accessing Postgres using the 'CreateShortUrl' method")
	logger.Debug().Msgf("postgres. userId: %d, short_url: %v", userId, arg)

	tx, err := r.client.Begin()
	if err != nil {
		logger.Debug().Msgf("transaction creation error. err: %s", err)
		return errTransaction
	}

	var id int
	createUrlQuery := fmt.Sprint("INSERT INTO referral (short_url) VALUES ($1) RETURNING id")
	row := tx.QueryRow(createUrlQuery, arg)
	if err = row.Scan(&id); err != nil {
		logger.Debug().Msgf("error writing to the 'referral' table. err: %s", err)

		tx.Rollback()
		return errTransaction
	}

	createUsersReferralQuery := fmt.Sprint("INSERT INTO users_referral (user_id, referral_id) VALUES ($1, $2)")
	_, err = tx.Exec(createUsersReferralQuery, userId, id)
	if err != nil {
		logger.Debug().Msgf("error writing to the 'users_referral' table. err: %s", err)

		tx.Rollback()
		return errTransaction
	}

	return tx.Commit()
}
