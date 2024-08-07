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

// GetUser - получение login по user_name
func (r *UserRepository) GetUser(ctx context.Context, user_name string) (string, error) {
	fmt.Println("обращение к Postgres GetUser")

	//r.log.Info("Пуск работы обработчика GET запроса")

	q := `
		SELECT login
		FROM users
		WHERE user_name = $1
	`

	var resURL string

	//r.log.Info("Alias: ", alias)

	err := r.client.QueryRow(ctx, q, user_name).Scan(&resURL)
	if err != nil {
		return "", fmt.Errorf("postgres: %w", err)
	}

	//r.log.Info("url has been found")

	return resURL, nil
}

// CreateUser - создание нового пользователя
func (r *UserRepository) CreateUser(ctx context.Context, name string, login string) (string, error) {
	fmt.Println("обращение к Postgres CreateUser")

	q := `
		INSERT INTO users 
		    (user_name, login) 
		VALUES 
		       ($1, $2) 
		RETURNING id
	`

	var id int

	if err := r.client.QueryRow(ctx, q, name, login).Scan(&id); err != nil {
		fmt.Printf("failed to create user: %w", err)

		return fmt.Sprintf("failed to create user"), fmt.Errorf("postgres: %w", err)
	}

	//r.log.Info("a new url has been saved")

	return fmt.Sprintf("successfully created user, id: %d", id), nil
}

// UpdateUser - обновление пользователя в таблице users
func (r *UserRepository) UpdateUser(ctx context.Context, name string, login string) (string, error) {
	q := `
	UPDATE users
	SET login = $1
	WHERE user_name = $2
`

	commandTag, err := r.client.Exec(ctx, q, login, name)
	if err != nil {
		return fmt.Sprintf("failed to update login user"), fmt.Errorf("postgres: %w", err)
	}

	if commandTag.RowsAffected() != 1 {
		return fmt.Sprintf("no row found to update"), errors.New("postgres:no row found to update")
	}

	return fmt.Sprintf("successfully updated login user"), nil
}

// DeleteUser - удаление пользователя из таблицы users по user_name
func (r *UserRepository) DeleteUser(ctx context.Context, name string) (string, error) {
	q := `
		DELETE FROM users
		WHERE user_name = $1
	`

	commandTag, err := r.client.Exec(ctx, q, name)
	if err != nil {
		return fmt.Sprintf("failed to delete user"), fmt.Errorf("postgres: %w", err)
	}

	if commandTag.RowsAffected() != 1 {
		return fmt.Sprintf("no row found to delete"), errors.New("postgres:no row found to delete")
	}

	//r.log.Info("url has been deleted")

	return fmt.Sprintf("successfully deleted user"), nil
}
