package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerHTTP struct {
		ListenAddress   string        `env:"HTTP_LISTEN_ADDRESS,notEmpty"`
		WriteTimeout    time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"90s"`
		ReadTimeout     time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"90s"`
		IdleTimeout     time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"60s"`
		ShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" envDefault:"30s"`
	}

	Postgres struct {
		DSN             string        `env:"PG_DSN,notEmpty" json:"-"`
		MaxIdleConns    int           `env:"PG_MAX_IDLE_CONNS" envDefault:"15"`
		MaxOpenConns    int           `env:"PG_MAX_OPEN_CONNS" envDefault:"15"`
		ConnMaxLifetime time.Duration `env:"PG_CONN_MAX_LIFETIME" envDefault:"5m"`
	}

	Redis struct {
		Username           string `env:"REDIS_USERNAME"`
		Password           string `env:"REDIS_PASSWORD" json:"-"`
		Address            string `env:"REDIS_ADDRESS,notEmpty"`
		DatabaseNumber     int    `env:"REDIS_DATABASE_NUMBER"`
		PoolSize           int    `env:"REDIS_POOL_SIZE" envDefault:"5"`
		MinIdleConnections int    `env:"REDIS_MIN_IDLE_CONNECTIONS" envDefault:"5"`
		MaxIdleConnections int    `env:"REDIS_MAX_IDLE_CONNECTIONS" envDefault:"10"`
	}

	Nats struct {
		URL string `env:"NATS_URL,notEmpty"`
	}

	JWT struct {
		PrivateKey string `env:"JWT_PRIVATE_KEY,notEmpty" json:"-"`
		PublicKey  string `env:"JWT_PUBLIC_KEY,notEmpty"`
		//SigningKey      string        `env:"SIGNING_KEY" yaml:"signingKey" env-default:""`
		AccessTokenTTL  time.Duration `env:"JWT_ACCESS_TOKEN_TTL" envDefault:"5m"`
		RefreshTokenTTL time.Duration `env:"JWT_REFRESH_TOKEN_TTL" envDefault:"50h"`
	}

	Logger struct {
		Level       string `env:"LOGGER_LEVEL" env-default:"info"`
		FieldMaxLen int    `env:"LOG_FIELD_MAX_LEN" envDefault:"2000"`
	}
}

func Load() (Config, error) {
	var config Config

	if err := env.Parse(&config); err != nil {
		return Config{}, fmt.Errorf("env.Parse: %w", err)
	}

	config.JWT.PrivateKey = correctNewlines(config.JWT.PrivateKey)
	config.JWT.PublicKey = correctNewlines(config.JWT.PublicKey)

	return config, nil
}

func correctNewlines(s string) string {
	return strings.NewReplacer(`"`, "", `\n`, "\n").Replace(s)
}
