package postgres

import (
	"context"
	"errors"
	"fmt"

	"online-store/pkg/client/postg"
)

type UserRepository struct {
	client postg.Client
}

func NewUserRepository(client postg.Client) *UserRepository {
	return &UserRepository{
		client: client,
	}
}

// GetUser - получение ID по login
func (r *UserRepository) GetUser(ctx context.Context, login string) (int, error) {
	fmt.Println("обращение к Postgres GetUser")

	//r.log.Info("Пуск работы обработчика GET запроса")

	q := `
		SELECT id
		FROM users
		WHERE login = $1
	`

	var id int

	//r.log.Info("Alias: ", alias)

	err := r.client.QueryRow(ctx, q, login).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("postgres: %w", err)
	}

	//r.log.Info("url has been found")

	return id, nil
}

// CreateUser - создание нового пользователя
func (r *UserRepository) CreateUser(ctx context.Context, login string, password string) (string, error) {
	fmt.Println("обращение к Postgres CreateUser")

	q := `
		INSERT INTO users 
		    (login, password) 
		VALUES 
		       ($1, $2) 
		RETURNING id
	`

	var id int

	if err := r.client.QueryRow(ctx, q, login, password).Scan(&id); err != nil {
		fmt.Printf("failed to create user: %w", err)

		return fmt.Sprintf("failed to create user"), fmt.Errorf("postgres: %w", err)
	}

	//r.log.Info("a new url has been saved")

	return fmt.Sprintf("successfully created user, id: %d", id), nil
}

// UpdateUser - обновление login, password пользователя в таблице users
func (r *UserRepository) UpdateUser(ctx context.Context, id int, login string, password string) (string, error) {
	q := `
	UPDATE users
	SET login = $2, password = $3
	WHERE id = $1
`

	commandTag, err := r.client.Exec(ctx, q, id, login, password)
	if err != nil {
		return fmt.Sprintf("failed to update login, password user"), fmt.Errorf("postgres: %w", err)
	}

	if commandTag.RowsAffected() != 1 {
		return fmt.Sprintf("no row found to update"), errors.New("postgres: the user does not exist")
	}

	return fmt.Sprintf("successfully updated login, password user"), nil
}

// DeleteUser - удаление пользователя из таблицы users по login
func (r *UserRepository) DeleteUser(ctx context.Context, id int) (string, error) {
	q := `
		DELETE FROM users
		WHERE id = $1
	`

	commandTag, err := r.client.Exec(ctx, q, id)
	if err != nil {
		return fmt.Sprintf("failed to delete user"), fmt.Errorf("postgres: %w", err)
	}

	if commandTag.RowsAffected() != 1 {
		return fmt.Sprintf("no row found to delete"), errors.New("postgres: the user does not exist")
	}

	//r.log.Info("url has been deleted")

	return fmt.Sprintf("successfully deleted user"), nil
}
