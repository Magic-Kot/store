package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Magic-Kot/store/pkg/client/postg"

	"github.com/rs/zerolog"
)

var (
	errUserNotFound = errors.New("postgres: user not found")
)

type UserRepository struct {
	client postg.Client
}

func NewUserRepository(client postg.Client) *UserRepository {
	return &UserRepository{
		client: client,
	}
}

// GetUser - получение сущности пользователя по ID
func (r *UserRepository) GetUser(ctx context.Context, id int) (string, string, error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("starting the GET request handler")

	q := `
		SELECT username, email
		FROM users
		WHERE id = $1
	`

	logger.Debug().Msgf("postgres: get user by id: %d\n", id)
	var username, email string
	//что лучше использовать map[string]string, struct User?

	// как обычно обращаются к базе?
	//tx, err := r.client.Begin(ctx)
	//row := tx.QueryRow(q, id)

	err := r.client.QueryRowx(q, id).Scan(&username, &email)

	logger.Debug().Msgf("postgres returned: login: %s, email: %s, err: %s\n", username, email, err)

	if errors.Is(err, sql.ErrNoRows) {
		return "", "", errUserNotFound
	} else if err != nil {
		return "", "", err
	}

	return username, email, nil
}

// GetAllUser - получение login всех пользователей

// CreateUser - создание нового пользователя
func (r *UserRepository) CreateUser(ctx context.Context, username string, passwordHash string) (int, error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("starting the POST request handler")

	q := `
		INSERT INTO users 
		    (username, password) 
		VALUES 
		       ($1, $2) 
		RETURNING id
	`

	var id int

	if err := r.client.QueryRowx(q, username, passwordHash).Scan(&id); err != nil {
		return 0, errors.New(fmt.Sprint("failed to create user. postgres: ", err))
	}

	return id, nil
}

// AuthorizationUser - авторизация пользователя
func (r *UserRepository) AuthorizationUser(ctx context.Context, login string) (string, error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("starting the POST request handler")

	return "", nil
}

// UpdateUser - обновление данных пользователя по ID
func (r *UserRepository) UpdateUser(ctx context.Context, value string, arg []interface{}) error {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("starting the PUT request handler")
	logger.Debug().Msgf("postgres: update user by id: %d, value: %s, arg: %v\n", arg[0], value, arg[1:])

	q := fmt.Sprintf(`UPDATE users SET %s WHERE id = $1`, value)

	commandTag, err := r.client.Exec(q, arg...)
	//тут должна быть обработка ошибок:
	// pq: syntax error at or near \"WHERE\"

	if err != nil {
		fmt.Println(fmt.Sprintf("%T", err))
		return errors.New(fmt.Sprint("failed to update login user. ", err))
	}

	if str, _ := commandTag.RowsAffected(); str != 1 {
		return errUserNotFound
	}

	return nil
}

// DeleteUser - удаление пользователя из таблицы users по login
func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("starting the DELETE request handler")

	q := `
		DELETE FROM users
		WHERE id = $1
	`

	commandTag, err := r.client.Exec(q, id)
	// тут должна быть обработка ошибок:
	//

	if err != nil {
		return errors.New(fmt.Sprint("failed to delete user: ", err))
	}

	if str, _ := commandTag.RowsAffected(); str != 1 {
		return errUserNotFound
	}

	return nil
}
