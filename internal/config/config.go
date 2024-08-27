package config

import "time"

type Config struct {
	ServerDeps   ServerDeps   `env:"SERVER" yaml:"server"`
	PostgresDeps PostgresDeps `env:"REPOSITORY" yaml:"repository"`
	LoggerDeps   LoggerDeps   `env:"LOGGER" yaml:"logger"`
	AuthDeps     AuthDeps     `env:"AUTH" yaml:"auth"`
}

type ServerDeps struct {
	Host    string        `env:"HOST"  yaml:"host" env-default:"localhost"`
	Port    string        `env:"PORT" yaml:"port" env-default:":8000"`
	Timeout time.Duration `env:"TIMEOUT" yaml:"timeout" env-default:"5s"`
}

type PostgresDeps struct {
	MaxAttempts int           `env:"MAX_ATTEMPTS" yaml:"maxAttempts" env-default:"3"`
	Delay       time.Duration `env:"DELAY" yaml:"delay" env-default:"10s"`
	Username    string        `env:"USERNAMEPOSTGRES"  yaml:"username" env-default:"postgres"`
	Password    string        `env:"PASSWORD" yaml:"password" env-default:"postgres"`
	Host        string        `env:"HOST" yaml:"host" env-default:"127.0.0.1"`
	Port        string        `env:"PORT" yaml:"port" env-default:"5432"`
	Database    string        `env:"DATABASE" yaml:"database" env-default:"postgres"`
	SSLMode     string        `env:"MODELESS" yaml:"sslMode" env-default:"disable"`
}

type LoggerDeps struct {
	LogLevel string `env:"LOG_LEVEL" yaml:"logLevel" env-default:"info"`
}

type AuthDeps struct {
	SigningKey      string        `env:"SIGNING_KEY" yaml:"signingKey" env-default:""`
	AccessTokenTTL  time.Duration `env:"ACCESS_TOKEN_TTL" yaml:"accessTokenTTL" env-default:"1h"`
	RefreshTokenTTL time.Duration `env:"REFRESH_TOKEN_TTL" yaml:"refreshTokenTTL" env-default:"4h"`
}
