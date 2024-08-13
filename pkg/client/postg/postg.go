package postg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

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
	Username    string
	Password    string
	Host        string
	Port        string
	Database    string
	SSLMode     string
}

func NewClient(ctx context.Context, cfg *ConfigDeps) (*sqlx.DB, error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("creating a Postgres client")
	logger.Debug().Msgf("config: %+v", cfg)

	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Database, cfg.Password, cfg.SSLMode))
	if err != nil {
		logger.Error().Msg(fmt.Sprint("errOpen:error connecting to Postgres:", err))
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logger.Error().Msg(fmt.Sprint("errPing: error connecting to Postgres:", err))
		return nil, err
	}

	return db, nil

	//delay := 5 * time.Second
	//attempts := cfg.MaxAttempts

	//for attempts > 0 {
	//	if db, err = Connection(cfg, attempts, delay); err != nil {
	//		time.Sleep(delay)
	//		attempts--
	//
	//		continue
	//	}
	//}
	//
	//if err != nil {
	//	logger.Info().Msg(fmt.Sprint("error connecting to Postgres:", err))
	//	return nil, err
	//}
	//
	//err = db.Ping()
	//if err != nil {
	//	logger.Info().Msg(fmt.Sprint("error connecting to Postgres:", err))
	//	return nil, err
	//}
	//
	//logger.Info().Msg("successful connection to Postgres")
	//
	//return db, err
}

//func Connection(cfg *ConfigDeps, attempts int, delay time.Duration) (db *sqlx.DB, err error) {
//	db, err = sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
//		cfg.Host, cfg.Port, cfg.Username, cfg.Database, cfg.Password, cfg.SSLMode))
//
//	return db, err
//}
