package main

import (
	"context"
	"fmt"
	"github.com/Magic-Kot/store/internal/config"
	"github.com/Magic-Kot/store/internal/controllers"
	"github.com/Magic-Kot/store/internal/delivery/httpecho"
	"github.com/Magic-Kot/store/internal/repository/postgres"
	"github.com/Magic-Kot/store/internal/services/user"
	"github.com/Magic-Kot/store/pkg/client/postg"
	"github.com/Magic-Kot/store/pkg/httpserver"
	"github.com/Magic-Kot/store/pkg/logging"
	"github.com/Magic-Kot/store/pkg/utils/jwt_token"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func main() {
	// read config
	//var cfg httpserver.ServerDeps
	var cfg config.Config

	err := cleanenv.ReadConfig("internal/config/config.yml", &cfg)
	//err := cleanenv.ReadConfig("config.yml", &cfg) // for docker
	//err := cleanenv.ReadConfig("internal/config/config.env", &cfg)
	//err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("error initializing config")
	}

	// create logger
	logCfg := logging.LoggerDeps{
		LogLevel: cfg.LoggerDeps.LogLevel,
	}

	logger, err := logging.NewLogger(&logCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init logger")
	}

	logger.Info().Msg("init logger")

	ctx := context.Background()
	ctx = logger.WithContext(ctx)

	logger.Debug().Msgf("config: %+v", cfg)

	// create server
	serv := httpserver.ConfigDeps{
		Host:    cfg.ServerDeps.Host,
		Port:    cfg.ServerDeps.Port,
		Timeout: cfg.ServerDeps.Timeout,
		//Logger:  logger,
	}

	server := httpserver.NewServer(&serv)

	// create client
	repo := postg.ConfigDeps{
		MaxAttempts: cfg.PostgresDeps.MaxAttempts,
		Delay:       cfg.PostgresDeps.Delay,
		Username:    cfg.PostgresDeps.Username,
		Password:    cfg.PostgresDeps.Password,
		Host:        cfg.PostgresDeps.Host,
		Port:        cfg.PostgresDeps.Port,
		Database:    cfg.PostgresDeps.Database,
		SSLMode:     cfg.PostgresDeps.SSLMode,
	}

	pool, err := postg.NewClient(ctx, &repo)
	if err != nil {
		logger.Fatal().Err(err).Msg(fmt.Sprint("NewClient: ", err))
	}

	// create repository
	db := postgres.NewUserRepository(pool)

	// create tokenJWT
	tokenCfg := jwt_token.TokenJWTDeps{
		SigningKey:      cfg.AuthDeps.SigningKey,
		AccessTokenTTL:  cfg.AuthDeps.AccessTokenTTL,
		RefreshTokenTTL: cfg.AuthDeps.RefreshTokenTTL,
	}

	tokenJWT, err := jwt_token.NewTokenJWT(&tokenCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init tokenJWT")
	}

	// create service
	service := user.NewUserService(db, tokenJWT)

	// create validator
	validate := validator.New()

	// create controller
	contr := controllers.NewApiController(service, logger, validate, tokenJWT)

	// set routes
	httpecho.SetUserRoutes(server.Server(), contr)

	// start server
	logger.Info().Msg("starting server")

	if err := server.Start(); err != nil {
		logger.Fatal().Msg(fmt.Sprint("serverStart:", err))
	}
}
