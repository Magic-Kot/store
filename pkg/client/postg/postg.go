package postg

import (
	"context"
	"fmt"
	"log"
	"online-store/internal/config"
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

//type Config struct {
//	MaxAttempts int    `yaml:"maxAttempts" env:"MAX_ATTEMPTS" env-default:"2"`
//	Username    string `yaml:"username" env:"USERNAME" env-default:"postgres"`
//	Password    string `yaml:"password" env:"PASSWORD" env-default:"postgres"`
//	Host        string `yaml:"host" env:"HOST" env-default:"127.0.0.1"`
//	Port        string `yaml:"port" env:"PORT" env-default:"5432"`
//	Database    string `yaml:"database" env:"DATABASE" env-default:"postgres"`
//}

// NewClient создает клиента, подключаемый к базе данных по URL: postgres://postgres:12345@localhost:5438/postgres
// func NewClient(ctx context.Context, config *Config) (pool *pgxpool.Pool, err error) {
func NewClient(ctx context.Context, config *config.RepositoryConfig) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", config.Username, config.Password, config.Host, config.Port, config.Database)
	err = DoWithTries(func() error {
		// With Timeout возвращает значение с указанием крайнего срока.
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		// Connect создает новый пул и немедленно устанавливает одно соединение. ctx можно использовать для отмены этого первоначального соединения.
		pool, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			return fmt.Errorf("connect function error: %w", err)
		}

		return nil
	}, config.MaxAttempts, 5*time.Second)

	if err != nil {
		log.Fatal("error do with tries postg: %w", err)
	}

	return pool, nil
}

func DoWithTries(fn func() error, attempts int, delay time.Duration) (err error) {
	for attempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attempts--

			continue
		}

		return nil
	}

	return
}
