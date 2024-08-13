package config

import "time"

type Config struct {
	ServerDeps   ServerDeps   `yaml:"server"`
	PostgresDeps PostgresDeps `yaml:"repository"`
	LoggerDeps   LoggerDeps   `yaml:"logger"`
}

type ServerDeps struct {
	Host    string        `yaml:"host" env:"HOST" env-default:"localhost"`
	Port    string        `yaml:"port" env:"PORT" env-default:":8000"`
	Timeout time.Duration `yaml:"timeout" env:"TIMEOUT" env-default:"5s"`
}

type PostgresDeps struct {
	MaxAttempts int    `yaml:"maxAttempts" env:"MAX_ATTEMPTS" env-default:"2"`
	Username    string `yaml:"username" env:"USERNAMEPOSTGRES" env-default:"postgres"`
	Password    string `yaml:"password" env:"PASSWORD" env-default:"postgres"`
	Host        string `yaml:"host" env:"HOST" env-default:"127.0.0.1"`
	Port        string `yaml:"port" env:"PORT" env-default:"5432"`
	Database    string `yaml:"database" env:"DATABASE" env-default:"postgres"`
	SSLMode     string `yaml:"sslMode" env:"MODELESS" env-default:"disable"`
}

type LoggerDeps struct {
	LogLevel string `yaml:"logLevel" env:"LOG_LEVEL" env-default:"info"`
}
