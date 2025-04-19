package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"git.appkode.ru/pub/go/failure"
	"github.com/jmoiron/sqlx"

	"github.com/Magic-Kot/store/internal/domain/entity"
	"github.com/Magic-Kot/store/internal/domain/value"
	"github.com/Magic-Kot/store/pkg/errcodes"
)

type AuthPostgresRepository struct {
	db *sqlx.DB
}

func NewAuthPostgresRepository(db *sqlx.DB) *AuthPostgresRepository {
	return &AuthPostgresRepository{
		db: db,
	}
}

func (r *AuthPostgresRepository) CreateUser(ctx context.Context, user entity.CreateUser) error {
	query := `
		INSERT INTO users (
			id,
			person_id,
			login,
			password,
			created_at
		) VALUES (
			:id,
			:person_id,
			:login,
			:password,
			:created_at
		)
	`

	if _, err := r.db.NamedExecContext(ctx, query, user); err != nil {
		return fmt.Errorf("r.db.ExecContext: %w", err)
	}

	return nil
}

func (r *AuthPostgresRepository) UserByLogin(ctx context.Context, login value.Login) (value.UserAuth, error) {
	query := `
		SELECT
		    person_id,
		    password
		FROM
		    users
		WHERE
		    login = $1
	`

	var user value.UserAuth

	err := r.db.GetContext(ctx, &user, query, login)

	if errors.Is(err, sql.ErrNoRows) {
		return user, failure.NewNotFoundError(
			fmt.Sprintf("user with login: %s not found", login),
			failure.WithCode(errcodes.UserNotFound),
			failure.WithDescription(errcodes.UserNotFoundMessage.String()))
	} else if err != nil {
		return user, fmt.Errorf("r.db.GetContext: %w", err)
	}

	return user, nil
}

func (r *AuthPostgresRepository) PersonIDLatest(ctx context.Context) (value.PersonID, error) {
	query := `
		SELECT
		    person_id
		FROM
		    users
		ORDER BY
		    created_at DESC
		LIMIT 1
	`

	var personID value.PersonID

	if err := r.db.GetContext(ctx, &personID, query); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return personID, failure.NewNotFoundError(
				fmt.Sprintf("personID not found"),
				failure.WithCode(errcodes.NotFound),
				failure.WithDescription(errcodes.NotFound.String()))
		}

		return personID, fmt.Errorf("db.GetContext: %w", err)
	}

	return personID, nil
}
