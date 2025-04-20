package persistence

import (
	"context"
	"fmt"

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
		    id = $1
	`

	var userInfo entity.UserData

	if err := u.db.GetContext(ctx, userInfo, query, personID); err != nil {
		return userInfo, fmt.Errorf("db.GetContext: %w", err)
	}

	return userInfo, nil
}
