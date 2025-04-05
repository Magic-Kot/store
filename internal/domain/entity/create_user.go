package entity

import (
	"time"

	"github.com/Magic-Kot/store/internal/domain/value"
)

type CreateUser struct {
	ID           string         `db:"id"`
	PersonID     value.PersonID `db:"person_id"`
	Login        value.Login    `db:"login"`
	PasswordHash string         `db:"password"`
	CreatedAt    time.Time      `db:"created_at"`
}
