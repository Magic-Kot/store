package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Magic-Kot/store/internal/models"
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
func (r *UserRepository) GetUser(ctx context.Context, user *models.User) (*models.User, error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("starting the GET request handler")

	q := `
		SELECT username, email
		FROM users
		WHERE id = $1
	`

	logger.Debug().Msgf("postgres: get user by id: %d\n", user.ID)
	var username, email string

	// как обычно обращаются к базе?
	//tx, err := r.client.Begin(ctx)
	//row := tx.QueryRow(q, id)

	err := r.client.QueryRowx(q, user.ID).Scan(&user.Username, &user.Email)

	logger.Debug().Msgf("postgres returned: login: %s, email: %s, err: %s\n", username, email, err)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errUserNotFound
	} else if err != nil {
		return nil, err
	}

	return user, nil
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

// SignIn - аутентификация пользователя, получение
func (r *UserRepository) SignIn(ctx context.Context, user *models.UserAuthorization) (*models.UserAuthorization, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("accessing Postgres using the 'SignIn' method")
	logger.Debug().Msgf("postgres SignIn: by login: %s", user.Username)

	q := `
		SELECT id, password
		FROM users
		WHERE username = $1
	`

	err := r.client.QueryRowx(q, user.Username).Scan(&user.ID, &user.Password)

	if errors.Is(err, sql.ErrNoRows) {
		logger.Debug().Msgf("user not found. postgres: %s", err)
		return nil, errUserNotFound
	} else if err != nil {
		logger.Debug().Msgf("failed to get user. postgres: %s", err)
		return nil, err
	}

	return user, nil
}

// CreateSession - создание сессии пользователя
func (r *UserRepository) CreateSession(ctx context.Context, value string, arg []interface{}) (int, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("accessing Postgres using the 'CreateSession' method")
	logger.Debug().Msgf("postgres: create session by id: %d, value: %s, arg: %v", arg[0], value, arg[1:])

	q := fmt.Sprintf(`INSERT INTO sessions (%s) VALUES ($1, $2, $3, $4) RETURNING id`, value)

	var id int

	if err := r.client.QueryRowx(q, arg...).Scan(&id); err != nil {
		logger.Debug().Msgf("failed to create session. postgres: %s", err)

		return 0, errors.New(fmt.Sprint("failed to create session. postgres: ", err))
	}

	return id, nil
}

// GetSession - получение сессии по userId пользователя
func (r *UserRepository) GetSession(ctx context.Context, table string, column string, value string, arg []interface{}) (string, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("accessing Postgres using the 'GetSession' method")
	logger.Debug().Msgf("postgres: get session by table: %s, column: %s, value: %s, arg: %v", table, column, value, arg)

	q := fmt.Sprintf(`SELECT %s FROM %s WHERE %s = $1`, value, table, column)

	var check string

	err := r.client.QueryRowx(q, arg...).Scan(&check)

	if errors.Is(err, sql.ErrNoRows) {
		logger.Debug().Msgf("session not found. postgres: %s", err)
		return "", errUserNotFound
	} else if err != nil {
		logger.Debug().Msgf("failed to get session. postgres: %s", err)
		return "", err
	}

	return check, nil
}

// UpdateUser - обновление данных в указанной таблице по указанному столбцу
func (r *UserRepository) UpdateUser(ctx context.Context, table string, column string, value string, arg []interface{}) error {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("accessing Postgres using the 'UpdateUser' method")
	logger.Debug().Msgf("postgres: update table by table: %s, column: %s, value: %s, arg: %v", table, column, value, arg)

	q := fmt.Sprintf(`UPDATE %s SET %s WHERE %s = $1`, table, value, column)

	commandTag, err := r.client.Exec(q, arg...)

	if err != nil {
		logger.Debug().Msgf("failed table updates. postgres: %s", err)
		return errors.New(fmt.Sprint("failed table updates: ", err))
	}

	if str, _ := commandTag.RowsAffected(); str != 1 {
		logger.Debug().Msgf("user not found. postgres: %s", err)
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
