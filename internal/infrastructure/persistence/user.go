package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"git.appkode.ru/pub/go/failure"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Magic-Kot/store/internal/domain/entity"
	"github.com/Magic-Kot/store/internal/domain/value"
)

type DBUser struct {
	db *sqlx.DB
}

func NewDBUser(db *sqlx.DB) DBUser {
	return DBUser{db: db}
}

func (u DBUser) UserInfo(ctx context.Context, personID value.PersonID) (entity.UserData, error) {
	const query = `
		SELECT
		    username,
		    name,
		    surname,
		    age,
		    email
		FROM
		    users
		WHERE
		    person_id = $1
	`

	var userInfo entity.UserData

	if err := u.db.GetContext(ctx, userInfo, query, personID); err != nil {
		return userInfo, fmt.Errorf("db.GetContext: %w", err)
	}

	return userInfo, nil
}

func (u DBUser) UpdateUser(ctx context.Context, personID value.PersonID, age int, username, name, surname string) error {
	const query = `
		UPDATE
		    users
		SET 
		    username = COALESCE($2, username),
		    name = COALESCE($3, name),
		    surname = COALESCE($4, surname),
		    age = COALESCE($5, age),
		    updated_at = $6
		WHERE
		    person_id = $1 and
		    deleted_at is null 
	`

	res, err := u.db.ExecContext(ctx, query, personID, username, name, surname, age, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("db.ExecContext: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return failure.NewNotFoundError("not found")
	}

	return nil
}

func (u DBUser) DeletedAtByPersonID(ctx context.Context, personID value.PersonID) (*time.Time, error) {
	const query = `
		SELECT
			deleted_at
		FROM
			users
		WHERE
			person_id = $1
		`

	var deletedAt sql.NullTime

	if err := u.db.QueryRowxContext(ctx, query, personID).Scan(&deletedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, failure.NewNotFoundError("not found")
		}

		return nil, fmt.Errorf("db.QueryRowxContext: %w", err)
	}

	if deletedAt.Valid {
		return &deletedAt.Time, nil
	}

	return nil, nil
}

func (u DBUser) Remove(ctx context.Context, personID value.PersonID) error {
	const sqlQuery = `
		UPDATE users
		SET
		    deleted_at = $2
		WHERE 
		    person_id = $1
	`

	res, err := u.db.ExecContext(ctx, sqlQuery, personID, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("db.ExecContext(): %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return failure.NewNotFoundError("not found")
	}

	return nil
}
