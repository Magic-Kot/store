package config

import "time"

type Config struct {
	ServerDeps       ServerDeps       `yaml:"server"`
	RepositoryConfig RepositoryConfig `yaml:"repository"`
}

type ServerDeps struct {
	Host    string        `yaml:"host" env:"HOST" env-default:"localhost"`
	Port    string        `yaml:"port" env:"PORT" env-default:":8000"`
	Timeout time.Duration `yaml:"timeout" env:"TIMEOUT" env-default:"5s"`
}

type RepositoryConfig struct {
	MaxAttempts int    `yaml:"maxAttempts" env-default:"2"`
	Username    string `yaml:"username" env-default:"postgres"`
	Password    string `yaml:"password" env-default:"postgres"`
	Host        string `yaml:"host" env-default:"127.0.0.1"`
	Port        string `yaml:"port" env-default:"5432"`
	Database    string `yaml:"database" env-default:"postgres"`
}
