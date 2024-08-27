package postg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

var errConnectingPostgres = errors.New("error connecting to postgres")

type Client interface {
	// Exec - выполняет запрос, не возвращая никаких строк. Аргументы предназначены для любых параметров-заполнителей в запросе.
	Exec(query string, args ...interface{}) (sql.Result, error)
	// Query - выполняет запрос, который возвращает строки, обычно SELECT. Аргументы предназначены для любых параметров-заполнителей в запросе.
	Query(query string, args ...interface{}) (*sql.Rows, error)
	// QueryRowx - QueryRowContext выполняет запрос, который, как ожидается, вернет не более одной строки. всегда возвращает ненулевое значение.
	//Ошибки откладываются до тех пор, пока не будет вызван метод проверки [Row].
	// Если запрос не выберет ни одной строки, [*Row.Scan] вернет [ErrNoRows]. В противном случае [*Row.Scan] сканирует первую выбранную строку и отбрасывает остальные.
	QueryRowx(query string, args ...interface{}) *sqlx.Row
}

type ConfigDeps struct {
	MaxAttempts int
	Delay       time.Duration
	Username    string
	Password    string
	Host        string
	Port        string
	Database    string
	SSLMode     string
}

// NewClient создает клиента, подключаемый к базе данных по URL: postgres://postgres:12345@localhost:5438/postgres
func NewClient(ctx context.Context, cfg *ConfigDeps) (db *sqlx.DB, err error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("creating a Postgres client")
	logger.Debug().Msgf("config: %+v", cfg)

	fn := func() error {
		db, err = sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.Username, cfg.Database, cfg.Password, cfg.SSLMode))

		if err != nil {
			logger.Debug().Msgf("error connecting to Postgres: %v", err)
			return err
		}
		return nil
	}

	err = Connection(fn, cfg.MaxAttempts, cfg.Delay)

	if err != nil {
		logger.Debug().Msgf("error connecting to Postgres: %v", err)
		return nil, errConnectingPostgres
	}

	logger.Info().Msg("successful connection to Postgres")

	return db, nil
}

func Connection(fn func() error, attempts int, delay time.Duration) (err error) {
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
