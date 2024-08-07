package postg

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

// NewClient создает клиента, подключаемый к базе данных по URL: postgres://postgres:12345@localhost:5438/postgres
func NewClient(ctx context.Context, maxAttempts int, username, password, host, port, database string) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", username, password, host, port, database)
	err = DoWithTries(func() error {
		// With Timeout возвращает значение с указанием крайнего срока.
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		// Connect создает новый пул и немедленно устанавливает одно соединение. ctx можно использовать для отмены этого первоначального соединения.
		//Информацию о формате connString смотрите в ParseConfig.
		pool, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			return fmt.Errorf("connect function error: %w", err)
		}

		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		log.Fatal("error do with tries postg: %w", err)
	}

	return pool, nil
}

func DoWithTries(fn func() error, attemtps int, delay time.Duration) (err error) {
	for attemtps > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attemtps--

			continue
		}

		return nil
	}

	return
}
